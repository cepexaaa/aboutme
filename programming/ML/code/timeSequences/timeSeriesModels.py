import numpy as np
import pandas as pd
from sklearn.metrics import mean_squared_error, mean_absolute_error
import matplotlib.pyplot as plt
import os
from statsmodels.tsa.holtwinters import ExponentialSmoothing
from statsmodels.tsa.statespace.sarimax import SARIMAX


class NaiveMeanModel:
    def __init__(self):
        self.mean_value = None
        self.name = "Naive_Mean"

    def fit(self, data):
        self.mean_value = np.mean(data)
        return self

    def predict(self, steps):
        return np.full(steps, self.mean_value)


class NaiveLastValueModel:
    def __init__(self):
        self.last_value = None
        self.name = "Naive_Last"

    def fit(self, data):
        self.last_value = data.iloc[-1] if hasattr(data, 'iloc') else data[-1]
        return self

    def predict(self, steps):
        return np.full(steps, self.last_value)


class HoltWintersModel:
    def __init__(self, seasonal_periods=52, trend='add', seasonal='add'):
        self.seasonal_periods = seasonal_periods
        self.trend = trend
        self.seasonal = seasonal
        self.model = None
        self.name = "Holt_Winters"

    def fit(self, data):
        self.model = ExponentialSmoothing(
            data,
            seasonal_periods=self.seasonal_periods,
            trend=self.trend,
            seasonal=self.seasonal,
            initialization_method="estimated"
        ).fit()
        return self

    def predict(self, steps):
        if self.model is None:
            raise ValueError("Model not fitted. Call fit() first.")
        return self.model.forecast(steps)


class SARIMAModel:
    def __init__(self, order=(1, 1, 1), seasonal_order=(1, 1, 1, 52)):
        self.order = order  # (p,d,q)
        self.seasonal_order = seasonal_order  # (P,D,Q,s)
        self.model = None
        self.name = "SARIMA"

    def fit(self, data):
        # self.model = SARIMAX(
        #     data,
        #     order=self.order,
        #     seasonal_order=self.seasonal_order,
        #     enforce_stationarity=False,
        #     enforce_invertibility=False
        # ).fit(disp=False)
        self.model = SARIMAX(
            data,
            order=(1, 0, 0),
            seasonal_order=(0, 0, 0, 0),
            enforce_stationarity=False,
            enforce_invertibility=False
        ).fit(disp=False)
        return self

    def predict(self, steps):
        if self.model is None:
            raise ValueError("Model not fitted. Call fit() first.")
        return self.model.forecast(steps)


class TimeSeriesModelComparator:
    def __init__(self):
        self.models = {}
        self.results = {}
        self.train_data = None
        self.test_data = None

    def add_model(self, name, model):
        self.models[name] = model

    def train_and_evaluate(self, train_data, test_data):
        self.train_data = train_data
        self.test_data = test_data
        self.results = {}

        print("Training time series models...")
        print(f"Train size: {len(train_data)}, Test size: {len(test_data)}")

        for name, model in self.models.items():
            try:
                print(f"\n--- Training {name} ---")

                model.fit(train_data)
                predictions = model.predict(len(test_data))

                mse = mean_squared_error(test_data, predictions)
                mae = mean_absolute_error(test_data, predictions)
                rmse = np.sqrt(mse)

                mape = np.mean(np.abs((test_data - predictions) / test_data)) * 100

                self.results[name] = {
                    'predictions': predictions,
                    'mse': mse,
                    'mae': mae,
                    'rmse': rmse,
                    'mape': mape,
                    'model': model
                }

                print(f"{name} Results:")
                print(f"  MSE: {mse:.4f}")
                print(f"  MAE: {mae:.4f}")
                print(f"  RMSE: {rmse:.4f}")
                print(f"  MAPE: {mape:.2f}%")

            except Exception as e:
                print(f"Error training {name}: {e}")
                continue

        return self.results

    def plot_comparison(self, filename="time_series_models_comparison"):
        if not self.results:
            print("No results to plot. Run train_and_evaluate first.")
            return

        fig, axes = plt.subplots(2, 2, figsize=(15, 12))

        ax = axes[0, 0]
        time_index = range(len(self.test_data))
        ax.plot(time_index, self.test_data.values, 'k-',
                label='Actual', linewidth=3, alpha=0.8)

        colors = ['red', 'blue', 'green', 'orange', 'purple', 'brown']
        for i, (name, result) in enumerate(self.results.items()):
            if i < len(colors):
                ax.plot(time_index, result['predictions'],
                        color=colors[i], linestyle='--',
                        label=name, linewidth=2, alpha=0.8)

        ax.set_title('Time Series Models: Actual vs Predicted', fontsize=14, fontweight='bold')
        ax.set_xlabel('Time Index')
        ax.set_ylabel('USD/RUB Exchange Rate')
        ax.legend()
        ax.grid(True, alpha=0.3)

        # Plot 2: Model errors (RMSE)
        ax = axes[0, 1]
        model_names = list(self.results.keys())
        rmse_values = [self.results[name]['rmse'] for name in model_names]

        bars = ax.bar(model_names, rmse_values, color='lightcoral', alpha=0.7)
        ax.set_title('Model RMSE Comparison', fontsize=14, fontweight='bold')
        ax.set_ylabel('Root Mean Squared Error')
        ax.tick_params(axis='x', rotation=45)

        # Add value labels on bars
        for bar, value in zip(bars, rmse_values):
            ax.text(bar.get_x() + bar.get_width() / 2, bar.get_height() + 0.1,
                    f'{value:.2f}', ha='center', va='bottom', fontweight='bold')

        ax.grid(True, alpha=0.3)

        # Plot 3: Prediction errors over time
        ax = axes[1, 0]
        for i, (name, result) in enumerate(self.results.items()):
            if i < len(colors):
                errors = np.abs(self.test_data.values - result['predictions'])
                ax.plot(time_index, errors, color=colors[i],
                        label=name, alpha=0.7)

        ax.set_title('Absolute Prediction Errors Over Time', fontsize=14, fontweight='bold')
        ax.set_xlabel('Time Index')
        ax.set_ylabel('Absolute Error')
        ax.legend()
        ax.grid(True, alpha=0.3)

        # Plot 4: MAPE comparison
        ax = axes[1, 1]
        mape_values = [self.results[name]['mape'] for name in model_names]

        bars = ax.bar(model_names, mape_values, color='lightblue', alpha=0.7)
        ax.set_title('Mean Absolute Percentage Error (MAPE)', fontsize=14, fontweight='bold')
        ax.set_ylabel('MAPE (%)')
        ax.tick_params(axis='x', rotation=45)

        # Add value labels on bars
        for bar, value in zip(bars, mape_values):
            ax.text(bar.get_x() + bar.get_width() / 2, bar.get_height() + 0.1,
                    f'{value:.1f}%', ha='center', va='bottom', fontweight='bold')

        ax.grid(True, alpha=0.3)

        plt.tight_layout()
        plt.savefig(f"./ds-cepexaaa/code/timeSequences/img/{filename}.png",
                    dpi=300, bbox_inches='tight')
        plt.close()

        print(f"\nComparison plots saved: {filename}.png")

    def get_best_model(self):
        if not self.results:
            return None

        best_name = min(self.results.keys(),
                        key=lambda x: self.results[x]['rmse'])
        return best_name, self.results[best_name]

    def print_summary(self):
        if not self.results:
            print("No results available.")
            return

        print("\n" + "=" * 60)
        print("TIME SERIES MODELS COMPARISON SUMMARY")
        print("=" * 60)

        # Sort models by RMSE
        sorted_models = sorted(self.results.items(),
                               key=lambda x: x[1]['rmse'])

        for i, (name, result) in enumerate(sorted_models, 1):
            print(f"\n{i}. {name}:")
            print(f"   RMSE:  {result['rmse']:.4f}")
            print(f"   MAE:   {result['mae']:.4f}")
            print(f"   MAPE:  {result['mape']:.2f}%")

        best_name, best_result = self.get_best_model()
        print(f"\nBEST MODEL: {best_name} (RMSE: {best_result['rmse']:.4f})")
        print("=" * 60)

