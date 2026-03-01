import numpy as np
from sklearn.base import BaseEstimator
from sklearn.utils.validation import check_X_y, check_array, check_is_fitted
from sklearn.utils.multiclass import unique_labels
from sklearn.metrics import accuracy_score

# Наследование от этих классов даёт нам геттеры и сеттеры полей класса
# и другие классы наследуют этот класс.
# Так что это для поддержания совместимости
class RidgeLinearClassifier(BaseEstimator):
    def __init__(self, alpha=1.0, fit_intercept=True):
        self.alpha = alpha
        self.fit_intercept = fit_intercept

    # add column of 1
    def _prepare_features(self, X):
        if self.fit_intercept:
            X = np.column_stack([np.ones(X.shape[0]), X])
        return X

    # Learning
    def fit(self, X, y):
        X, y = check_X_y(X, y)
        self.classes_ = unique_labels(y)
        y_binary = np.where(y == self.classes_[0], -1, 1)
        X_processed = self._prepare_features(X)
        n_features = X_processed.shape[1]

        XT = X_processed.T
        XTX = XT @ X_processed
        XTy = XT @ y_binary

        # regularization
        I = np.eye(n_features)
        if self.fit_intercept:
            I[0, 0] = 0  # Do not regularize intercept

        # ridge regularization: w = (X^T*X + αI)^{-1}X^T*y
        self.coef_ = np.linalg.inv(XTX + self.alpha * I) @ XTy

        # for compatibility
        self.n_features_in_ = X.shape[1]
        self.is_fitted_ = True

        return self

    def predict(self, X):
        decisions = self.decision_function(X)
        return np.where(decisions >= 0, self.classes_[1], self.classes_[0])

    # distance to the separating hyperplane
    def decision_function(self, X):
        check_is_fitted(self)
        X = check_array(X)
        X_processed = self._prepare_features(X)
        return X_processed @ self.coef_

    # Probabilities of belonging to classes
    def predict_proba(self, X):
        check_is_fitted(self)
        X = check_array(X)

        decisions = self.decision_function(X)
        # A sigmoid for converting to probabilities
        proba_positive = 1 / (1 + np.exp(-decisions))
        proba_negative = 1 - proba_positive
        return np.column_stack([proba_negative, proba_positive])

    # Classification accuracy
    def score(self, X, y):
        return accuracy_score(y, self.predict(X))

    def compute_loss(self, margin):
        losses = np.log1p(np.exp(-margin))  # loss: log(1 + exp(-margin))
        return np.mean(losses)




class GradientDescentLinearClassifier(BaseEstimator):
    # loss - function of loss
    # learning_rate - speed og learning
    # alpha - regularisation
    # l1_ratio - ratio: L1/L2
    # tol - criteria for stopping, changing the loss function
    # verbose - print processing of learning
    def __init__(self, loss='hinge', learning_rate=0.01, alpha=0.0001,
                 l1_ratio=0.5, max_iter=1000, tol=1e-4, random_state=None,
                 batch_size=None, verbose=False):
        self.loss = loss
        self.learning_rate = learning_rate
        self.alpha = alpha
        self.l1_ratio = l1_ratio
        self.max_iter = max_iter
        self.tol = tol
        self.random_state = random_state
        self.batch_size = batch_size
        self.verbose = verbose

    def _prepare_data(self, X, y):
        self.classes_ = unique_labels(y)
        if len(self.classes_) != 2:
            raise ValueError("Only binary classification")

        y_binary = np.where(y == self.classes_[0], -1, 1)

        n_samples, n_features = X.shape
        X_processed = np.column_stack([np.ones(n_samples), X])

        # initialise weights
        if self.random_state is not None:
            np.random.seed(self.random_state)
        self.coef_ = np.random.randn(n_features + 1) * 0.01

        return X_processed, y_binary

    def _compute_margin(self, X, y):
        return y * (X @ self.coef_)

    def _compute_loss(self, margin):
        if self.loss == 'hinge':
            losses = np.maximum(0, 1 - margin)
        elif self.loss == 'logistic':
            # Logistic loss: L = log(1 + exp(-margin))
            losses = np.log1p(np.exp(-margin))
        elif self.loss == 'exponential':
            losses = np.exp(-margin)
        else:
            raise ValueError(f"Неизвестная функция потерь: {self.loss}")

        return np.mean(losses)

    def _compute_data_gradient(self, X, y, margin):
        n_samples = X.shape[0]

        if self.loss == 'hinge':
            data_gradient = - (y * (margin < 1).astype(float)) @ X / n_samples
        elif self.loss == 'logistic':
            sigma = 1 / (1 + np.exp(margin))  # σ(-margin) = 1/(1+exp(margin))
            data_gradient = - (y * sigma.astype(float)) @ X / n_samples
        elif self.loss == 'exponential':
            data_gradient = - (y * np.exp(-margin)) @ X / n_samples
        else:
            raise ValueError(f"Неизвестная функция потерь: {self.loss}")

        return data_gradient

    # regularisation gradient Elastic Net
    def _compute_regularization_gradient(self):
        w = self.coef_[1:]  # Do not regularise intercept

        # L1 regularisation (subgradient)
        l1_grad = self.l1_ratio * np.where(w != 0, np.sign(w), 0)

        # L2 regularisation
        l2_grad = (1 - self.l1_ratio) * 2 * w

        return np.concatenate([[0], self.alpha * (l1_grad + l2_grad)])

    def _compute_total_gradient(self, X, y, margin):
        data_gradient = self._compute_data_gradient(X, y, margin)
        reg_gradient = self._compute_regularization_gradient()

        return data_gradient + reg_gradient

    def _get_batch_indices(self, n_samples):
        if self.batch_size is None or self.batch_size >= n_samples:
            return slice(None)  # Full batch

        indices = np.random.choice(n_samples, self.batch_size, replace=False)
        return indices

    def fit(self, X, y, X_val=None, y_val=None):
        X, y = check_X_y(X, y)
        X_processed, y_binary = self._prepare_data(X, y)

        n_samples = X_processed.shape[0]
        previous_loss = float('inf')
        self.loss_history_ = []
        self.coef_history_ = []
        self.val_loss_history_ = []

        for iteration in range(self.max_iter):
            batch_indices = self._get_batch_indices(n_samples)
            X_batch = X_processed[batch_indices]
            y_batch = y_binary[batch_indices]

            margin = self._compute_margin(X_batch, y_batch)
            data_loss = self._compute_loss(margin)

            w = self.coef_[1:]
            l1_penalty = self.l1_ratio * np.sum(np.abs(w))
            l2_penalty = (1 - self.l1_ratio) * np.sum(w ** 2)
            reg_loss = self.alpha * (l1_penalty + l2_penalty)

            total_loss = data_loss + reg_loss
            self.loss_history_.append(total_loss)
            self.coef_history_.append(self.coef_.copy())

            gradient = self._compute_total_gradient(X_batch, y_batch, margin)
            self.coef_ -= self.learning_rate * gradient

            if X_val is not None and y_val is not None:
                X_val_processed = np.column_stack([np.ones(X_val.shape[0]), X_val])
                y_val_binary = np.where(y_val == self.classes_[0], -1, 1)
                margin_val = y_val_binary * (X_val_processed @ self.coef_)
                val_loss = self._compute_loss(margin_val)
                self.val_loss_history_.append(val_loss)

            if iteration > 0 and abs(previous_loss - total_loss) < self.tol:
                if self.verbose:
                    print(f"Stop on iteration {iteration}")
                break

            previous_loss = total_loss

            if self.verbose and iteration % 100 == 0:
                print(f"Iteration {iteration}: Loss = {total_loss:.6f}")

        self.n_iter_ = iteration + 1
        self.n_features_in_ = X.shape[1]
        self.is_fitted_ = True

        return self #.loss_history_

    def decision_function(self, X):
        check_is_fitted(self)
        X = check_array(X)

        X_processed = np.column_stack([np.ones(X.shape[0]), X])
        return X_processed @ self.coef_

    def predict(self, X):
        decisions = self.decision_function(X)
        return np.where(decisions >= 0, self.classes_[1], self.classes_[0])

    # Probabilities of belonging to classes
    def predict_proba(self, X):
        check_is_fitted(self)
        X = check_array(X)

        decisions = self.decision_function(X)

        if self.loss == 'logistic':
            # For logistic loss use sigmoid
            proba_positive = 1 / (1 + np.exp(-decisions))
        else:
            #
            proba_positive = 1 / (1 + np.exp(-2 * decisions))

        proba_negative = 1 - proba_positive
        return np.column_stack([proba_negative, proba_positive])

    def score(self, X, y):
        y_pred = self.predict(X)
        return accuracy_score(y, y_pred)

    def get_loss_history(self):
        check_is_fitted(self)
        return self.loss_history_

