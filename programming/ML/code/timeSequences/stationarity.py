import pandas as pd
import numpy as np
import matplotlib.pyplot as plt
from statsmodels.tsa.stattools import adfuller
from statsmodels.graphics.tsaplots import plot_acf, plot_pacf


def test_stationarity(series, title="Time sequence"):
    result = adfuller(series.dropna())

    print(f" {title}:")
    print(f"   ADF stats: {result[0]:.6f}")
    print(f"   p-value: {result[1]:.10f}")
    print(f"   Critical values:")
    for key, value in result[4].items():
        print(f"     {key}: {value:.6f}")

    is_stationary = result[1] < 0.05
    if is_stationary:
        print("    STATION")
    else:
        print("    NO STATION")

    return is_stationary, result[1]


def tolog(data, value_col='usd_rub'):
    if isinstance(data, pd.DataFrame):
        df = data.copy()
        if 'date' in df.columns:
            df.set_index('date', inplace=True)
        return np.log(df[value_col])
    else:  # Series
        return np.log(data)


def todiff(data, value_col='usd_rub', order=1):
    if isinstance(data, pd.DataFrame):
        df = data.copy()
        if 'date' in df.columns:
            df.set_index('date', inplace=True)
        return df[value_col].diff(order).dropna()
    else:  # Series
        return data.diff(order).dropna()


def analyze_transformations(data, value_col='usd_rub'):
    transformations = {
        'original': data[value_col],
        'log': tolog(data, value_col),
        'diff_1': todiff(data, value_col, 1),
        'diff_2': todiff(data, value_col, 2),
        'log_diff': todiff(tolog(data))
    }

    results = {}

    print(" Analise of station transformations")

    for name, transformed_series in transformations.items():
        is_stationary, p_value = test_stationarity(transformed_series, f"{name.upper()} series")
        results[name] = {
            'series': transformed_series,
            'is_stationary': is_stationary,
            'p_value': p_value,
            'transformation': name
        }

    return results


def plot_transformations(results, filename_prefix="stationarity"):
    fig, axes = plt.subplots(3, 2, figsize=(15, 12))
    axes = axes.flatten()

    transformations = list(results.keys())

    for i, (name, result) in enumerate(results.items()):
        if i >= len(axes)-1:
            break

        ax = axes[i]
        series = result['series']

        # main plot
        ax.plot(series.index, series.values, linewidth=1)
        ax.set_title(f'{name.upper()} (p-value: {result["p_value"]:.4f})', fontweight='bold')
        ax.grid(True, alpha=0.3)

        if result['is_stationary']:
            ax.set_facecolor('#f0fff0')  # green = good

    # delete extra subplots
    for i in range(len(transformations) * 2, len(axes)):
        fig.delaxes(axes[i])

    plt.tight_layout()
    plt.savefig(f"./ds-cepexaaa/code/timeSequences/img/{filename_prefix}_transformations.png")
    plt.close()


class TimeSeriesTransformer:
    def __init__(self):
        self.transformations_applied = []
        self.original_data = None
        self.last_value = 0
        self.first_value = 0

    def apply_transformations(self, data, value_col='usd_rub', target_p_value=0.05):
        self.original_data = data.copy()
        self.first_value = self.original_data[value_col].iloc[0]
        if 'date' in self.original_data.columns:
            self.original_data.set_index('date', inplace=True)

        current_data = self.original_data[value_col]
        self.transformations_applied = []

        # Testing original series
        is_stationary, p_value = test_stationarity(current_data, "original series")

        if current_data.min() > 0 and p_value > target_p_value:
            print(" Using logarithm...")
            current_data = np.log(current_data)
            self.transformations_applied.append('log')
            is_stationary, p_value = test_stationarity(current_data, "After logarithm")

        diff_order = 0
        while p_value > target_p_value and diff_order < 3:
            diff_order += 1
            print(f"  Using differentiation {diff_order}...")
            current_data = current_data.diff().dropna()
            self.transformations_applied.append(f'diff_{diff_order}')
            is_stationary, p_value = test_stationarity(current_data, f"After {diff_order} differentiation")

        self.last_value = self.original_data[value_col].iloc[-1] if diff_order > 0 else None

        print(f"RESULT: Transformations: {self.transformations_applied}")
        print(f" Final p-value: {p_value:.6f}")

        return current_data

    def inverse_transform(self, transformed_data):
        if not self.transformations_applied:
            return transformed_data

        result = transformed_data.copy()

        for transformation in reversed(self.transformations_applied):
            if transformation.startswith('diff'):
                result.iloc[0] = result.iloc[0] + self.first_value
                result = result.cumsum()

            elif transformation == 'log':
                result = np.exp(result)

        return result

    def get_transformation_info(self):
        return {
            'transformations': self.transformations_applied,
            'needs_inverse': len(self.transformations_applied) > 0
        }