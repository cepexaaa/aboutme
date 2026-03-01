from sklearn.linear_model import LogisticRegression
import numpy as np
from sklearn.base import BaseEstimator

# Вычисляет хи-квадрат статистику для каждого признака
# Бинаризует признаки по медиане (выше/ниже медианы)
# Строит таблицу сопряженности 2xN_classes
# Вычисляет хи-квадрат = sum((набл - ожид)^2 / ожид)
class FilterFeatureSelector:
    def __init__(self, k: int = 1000):
        self.k = k
        self.selected_indices_ = None
        self.scores_ = None

    def chi_square(self, X, y):
        n_samples, n_features = X.shape
        n_classes = len(np.unique(y))

        if hasattr(X, 'toarray'):
            X_dense = X.toarray()
        else:
            X_dense = X

        scores = np.zeros(n_features)

        for feature_idx in range(n_features):
            feature_values = X_dense[:, feature_idx]
            threshold = np.median(feature_values)
            feature_binary = (feature_values > threshold).astype(int)
            contingency_table = np.zeros((2, n_classes))

            for class_idx in range(n_classes):
                mask = (y == class_idx)
                contingency_table[0, class_idx] = np.sum((feature_binary[mask] == 0))
                contingency_table[1, class_idx] = np.sum((feature_binary[mask] == 1))

            row_sums = contingency_table.sum(axis=1, keepdims=True)
            col_sums = contingency_table.sum(axis=0, keepdims=True)
            total = contingency_table.sum()

            # Вероятность иметь значение признака i: P(признак=i) = row_sum_i / total
            # Вероятность принадлежать классу j: P(класс = j) = col_sum_j / total
            # P(признак = i И класс = j) = P(признак = i) × P(класс = j)
            expected = row_sums @ col_sums / total

            expected = np.where(expected == 0, 1e-10, expected)

            chi2 = np.sum((contingency_table - expected) ** 2 / expected)
            scores[feature_idx] = chi2

        return scores

    def fit(self, X, y):
        self.scores_ = self.chi_square(X, y)

        if self.k > len(self.scores_):
            self.k = len(self.scores_)

        self.selected_indices_ = np.argsort(self.scores_)[-self.k:][::-1]

        return self

    def transform(self, X):
        return X[:, self.selected_indices_]

    def fit_transform(self, X, y):
        self.fit(X, y)
        return self.transform(X)

# Начинает со всех признаков (support_ = все True)
# Пока не останется n_features_to_select:
#   Обучает модель на текущем подмножестве
#   Получает важность признаков (coef_ или feature_importances_)
#   Удаляет step наименее важных признаков
# Запоминает ranking_ (порядок удаления)
class WrapperFeatureSelector:
    def __init__(self, estimator: BaseEstimator, n_features_to_select: int = 100,
                 step: int = 50, scoring: str = 'f1'):
        self.estimator = estimator
        self.n_features_to_select = n_features_to_select
        self.step = step
        self.scoring = scoring
        self.support_ = None
        self.ranking_ = None
        self.feature_importances_ = None

    def _get_feature_importance(self, estimator):
        if hasattr(estimator, 'coef_'):
            if estimator.coef_.ndim == 1:
                return np.abs(estimator.coef_)
            else:
                return np.mean(np.abs(estimator.coef_), axis=0)
        elif hasattr(estimator, 'feature_importances_'):
            return estimator.feature_importances_
        else:
            raise ValueError("The model must have coef_ or feature_importances_")

    def fit(self, X, y):
        n_samples, n_features = X.shape

        self.support_ = np.ones(n_features, dtype=bool)
        self.ranking_ = np.ones(n_features, dtype=int)

        current_ranking = 1

        while np.sum(self.support_) > self.n_features_to_select:
            features = np.where(self.support_)[0]

            if len(features) == 0:
                break

            X_subset = X[:, features]
            self.estimator.fit(X_subset, y)
            importance = self._get_feature_importance(self.estimator)
            sorted_indices = np.argsort(importance)
            n_to_eliminate = min(self.step, len(features) - self.n_features_to_select)

            if n_to_eliminate > 0:
                to_eliminate_subset = sorted_indices[:n_to_eliminate]
                to_eliminate_original = features[to_eliminate_subset]
                self.support_[to_eliminate_original] = False
                self.ranking_[to_eliminate_original] = current_ranking
                current_ranking += 1
        self.ranking_[self.support_] = current_ranking
        X_final = X[:, self.support_]
        self.estimator.fit(X_final, y)
        final_importance = self._get_feature_importance(self.estimator)
        self.feature_importances_ = np.zeros(n_features)
        self.feature_importances_[self.support_] = final_importance
        return self

    def transform(self, X):
        if self.support_ is None:
            raise ValueError("need to run fit()")

        return X[:, self.support_]

    def get_support(self, indices=False):
        if indices:
            return np.where(self.support_)[0]
        return self.support_

    def fit_transform(self, X, y):
        self.fit(X, y)
        return self.transform(X)

# Использует L1-регуляризованную логистическую регрессию
# L1-регуляризация (лассо) зануляет неважные веса
# Отбирает признаки с ненулевыми коэффициентами
# Если все занулились - выбирает топ N//10 признаков
class EmbeddedFeatureSelector:
    def __init__(self, C: float = 1.0, max_iter: int = 1000,
                 threshold: float = 1e-4, solver: str = 'saga'):
        self.C = C
        self.max_iter = max_iter
        self.threshold = threshold
        self.solver = solver
        self.coef_ = None
        self.support_ = None
        self.estimator_ = None

    def fit(self, X, y):
        self.estimator_ = LogisticRegression(
            penalty='l1',
            C=self.C,
            max_iter=self.max_iter,
            solver=self.solver,
            random_state=42
        )

        self.estimator_.fit(X, y)

        if self.estimator_.coef_.ndim == 1:
            self.coef_ = self.estimator_.coef_
        else:
            self.coef_ = np.max(np.abs(self.estimator_.coef_), axis=0)

        self.support_ = np.abs(self.coef_) > self.threshold

        if not np.any(self.support_):
            n_features = len(self.coef_)
            n_to_select = max(1, n_features // 10)
            top_indices = np.argsort(np.abs(self.coef_))[-n_to_select:]
            self.support_ = np.zeros(n_features, dtype=bool)
            self.support_[top_indices] = True

        return self

    def transform(self, X):
        if self.support_ is None:
            raise ValueError("need to run fit()")

        return X[:, self.support_]

    def fit_transform(self, X, y):
        self.fit(X, y)
        return self.transform(X)

    def get_feature_importance(self):
        if self.coef_ is None:
            raise ValueError("need to run fit()")

        return np.abs(self.coef_)

    def get_nonzero_coef_count(self):
        if self.support_ is None:
            raise ValueError("need to run fit()")

        return np.sum(self.support_)

