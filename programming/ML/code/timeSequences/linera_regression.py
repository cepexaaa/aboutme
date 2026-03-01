import pandas as pd
import numpy as np
from sklearn.preprocessing import StandardScaler
from sklearn.base import BaseEstimator
from sklearn.utils.validation import check_X_y, check_array, check_is_fitted
from sklearn.metrics import mean_squared_error, mean_absolute_error, r2_score

class RidgeLinearRegressor(BaseEstimator):
    def __init__(self, alpha=1.0, fit_intercept=True):
        self.alpha = alpha
        self.fit_intercept = fit_intercept

    def _prepare_features(self, X):
        if self.fit_intercept:
            X = np.column_stack([np.ones(X.shape[0]), X])
        return X

    def fit(self, X, y):
        X, y = check_X_y(X, y)
        X_processed = self._prepare_features(X)
        n_features = X_processed.shape[1]

        XT = X_processed.T
        XTX = XT @ X_processed
        XTy = XT @ y

        # Regularization
        I = np.eye(n_features)
        if self.fit_intercept:
            I[0, 0] = 0  # Do not regularize intercept

        # Ridge regression: w = (X^T*X + αI)^{-1}X^T*y
        self.coef_ = np.linalg.inv(XTX + self.alpha * I) @ XTy

        self.n_features_in_ = X.shape[1]
        self.is_fitted_ = True

        return self

    def predict(self, X):
        check_is_fitted(self)
        X = check_array(X)
        X_processed = self._prepare_features(X)
        return X_processed @ self.coef_

    def score(self, X, y):
        y_pred = self.predict(X)
        return r2_score(y, y_pred)




def create_time_features(data, date_column='date'):
    df = data.copy()

    if date_column in df.columns:
        df[date_column] = pd.to_datetime(df[date_column])
        dates = df[date_column]
    else:
        dates = pd.to_datetime(df.index)

    representative_date = dates

    df['week_of_year'] = representative_date.dt.isocalendar().week
    df['month'] = representative_date.dt.month
    df['quarter'] = representative_date.dt.quarter
    df['year'] = representative_date.dt.year

    df['week_of_year_sin'] = np.sin(2 * np.pi * df['week_of_year'] / 52)
    df['week_of_year_cos'] = np.cos(2 * np.pi * df['week_of_year'] / 52)

    df['month_sin'] = np.sin(2 * np.pi * df['month'] / 12)
    df['month_cos'] = np.cos(2 * np.pi * df['month'] / 12)

    return df


def prepare_regression_data(data, target_col='usd_rub', test_size=0.2):
    df_with_features = create_time_features(data)

    feature_columns = [col for col in df_with_features.columns
                       if col not in ['date', target_col] and not col.startswith('_')]

    X = df_with_features[feature_columns]
    y = df_with_features[target_col]

    split_idx = int(len(X) * (1 - test_size))

    X_train = X.iloc[:split_idx]
    X_test = X.iloc[split_idx:]
    y_train = y.iloc[:split_idx]
    y_test = y.iloc[split_idx:]

    scaler = StandardScaler()
    X_train_scaled = scaler.fit_transform(X_train)
    X_test_scaled = scaler.transform(X_test)

    return {
        'X_train': X_train_scaled,
        'X_test': X_test_scaled,
        'y_train': y_train,
        'y_test': y_test,
        'feature_names': feature_columns,
        'scaler': scaler
    }


def evaluate_regression_model(model, X_test, y_test, model_name="Ridge Regression"):
    y_pred = model.predict(X_test)

    mse = mean_squared_error(y_test, y_pred)
    mae = mean_absolute_error(y_test, y_pred)
    r2 = r2_score(y_test, y_pred)

    print(f"\n{model_name} Evaluation:")
    print(f"Mean Squared Error: {mse:.4f}")
    print(f"Mean Absolute Error: {mae:.4f}")
    print(f"R^2 Score: {r2:.4f}")
    print(f"Root MSE: {np.sqrt(mse):.4f}")

    return {
        'mse': mse,
        'mae': mae,
        'r2': r2,
        'predictions': y_pred,
        'rmse': np.sqrt(mse)
    }