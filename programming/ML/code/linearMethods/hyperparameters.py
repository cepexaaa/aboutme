import numpy as np
from sklearn.model_selection import cross_val_score, StratifiedKFold
from sklearn.metrics import accuracy_score
from algorythm import GradientDescentLinearClassifier, RidgeLinearClassifier
from itertools import product


class HyperparameterOptimizer:
    def __init__(self, cv_folds=5, random_state=42, n_jobs=-1):
        self.cv_folds = cv_folds
        self.random_state = random_state
        self.n_jobs = n_jobs
        self.best_params_ = {}
        self.best_score_ = {}

    def optimize_ridge_classifier(self, X, y, param_grid=None):
        best_score = -np.inf
        best_params = {}

        alpha_values = [0.001, 0.01, 0.1, 1.0, 10.0, 100.0]
        fit_intercept_options = [True, False]

        cv = StratifiedKFold(n_splits=self.cv_folds, shuffle=True, random_state=self.random_state)

        for alpha in alpha_values:
            for fit_intercept in fit_intercept_options:
                model = RidgeLinearClassifier(alpha=alpha, fit_intercept=fit_intercept)
                scores = cross_val_score(model, X, y, cv=cv, scoring='accuracy')
                mean_score = np.mean(scores)

                if mean_score > best_score:
                    best_score = mean_score
                    best_params = {'alpha': alpha, 'fit_intercept': fit_intercept}

                print(f"alpha={alpha:.3f}, intercept={fit_intercept}: CV Accuracy = {mean_score:.4f}")

        self.best_params_['ridge'] = best_params
        self.best_score_['ridge'] = best_score
        return best_params, best_score

    def optimize_gradient_classifier(self, X, y, param_grid=None):
        best_score = -np.inf
        best_params = {}

        param_grid = {
            'loss': ['hinge', 'logistic', 'exponential'],
            'learning_rate': [0.001, 0.01, 0.1],
            'alpha': [0.0001, 0.001, 0.01, 0.1],
            'l1_ratio': [0.0, 0.3, 0.5, 0.7, 1.0]
        }
        param_combinations = list(product(
            param_grid['loss'],
            param_grid['learning_rate'],
            param_grid['alpha'],
            param_grid['l1_ratio']
        ))

        # param_combinations = [
        #     ('hinge', 0.01, 0.001, 0.5),
        #     ('logistic', 0.01, 0.001, 0.5),
        #     ('exponential', 0.01, 0.001, 0.5),
        #     ('hinge', 0.1, 0.01, 0.3),
        #     ('logistic', 0.1, 0.01, 0.7),
        # ]

        cv = StratifiedKFold(n_splits=3, shuffle=True, random_state=self.random_state)

        for loss, lr, alpha, l1_ratio in param_combinations:
            model = GradientDescentLinearClassifier(
                loss=loss, learning_rate=lr, alpha=alpha,
                l1_ratio=l1_ratio, max_iter=1000, random_state=self.random_state, verbose=False
            )

            scores = cross_val_score(model, X, y, cv=cv, scoring='accuracy')
            mean_score = np.mean(scores)

            if mean_score > best_score:
                best_score = mean_score
                best_params = {
                    'loss': loss, 'learning_rate': lr,
                    'alpha': alpha, 'l1_ratio': l1_ratio
                }

            print(f"loss={loss}, lr={lr}, alpha={alpha}, l1_ratio={l1_ratio}: CV Accuracy = {mean_score:.4f}")

        self.best_params_['gradient'] = best_params
        self.best_score_['gradient'] = best_score
        return best_params, best_score


