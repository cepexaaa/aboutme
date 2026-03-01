import math
import os

import matplotlib.pyplot as plt
import pandas as pd
import numpy as np

def toImg(data, filename,
          show_moving_average=False,
          show_confidence_interval=False,
          show_trend_line=False,
          show_grid=True,
          show_legend=True,
          title=None,
          xlabel=None,
          ylabel=None,
          figsize=(12, 6)):

    df = data.copy()
    if 'date' in df.columns:
        df['date'] = pd.to_datetime(df['date'])
        df.set_index('date', inplace=True)
    elif df.index.dtype == 'object':
        df.index = pd.to_datetime(df.index)

    if 'usd_rub' in df.columns:
        value_col = 'usd_rub'
    else:
        value_col = df.columns[0]

    plt.figure(figsize=figsize)

    plt.plot(df.index, df[value_col], linewidth=1.5, color='blue', alpha=0.8, label='USD ~ RUB')
    if show_moving_average:
        window = int(math.sqrt(len(data)))
        ma_values = df[value_col].rolling(window=window).mean()
        plt.plot(df.index, ma_values, linewidth=2, color='red',
                            alpha=0.7, label=f'Moving average ({window} days.)')

    if show_confidence_interval:
        window = min(30, len(df) // 10)
        ma_values = df[value_col].rolling(window=window).mean()
        std_values = df[value_col].rolling(window=window).std()

        plt.fill_between(df.index,
                         ma_values - 1.96 * std_values,
                         ma_values + 1.96 * std_values,
                         color='red', alpha=0.2, label='95% confidence interval')

    if show_trend_line:
        x_numeric = np.arange(len(df))
        y_values = df[value_col].values

        coefficients = np.polyfit(x_numeric, y_values, 1)
        trend_line = np.polyval(coefficients, x_numeric)

        plt.plot(df.index, trend_line, '--', linewidth=2,
                               color='green', alpha=0.7, label='trend line')

        slope = coefficients[0]
        intercept = coefficients[1]
        trend_eq = f'y = {slope:.4f}x + {intercept:.2f}'

        plt.annotate(trend_eq, xy=(0.02, 0.02), xycoords='axes fraction',
                     bbox=dict(boxstyle='round,pad=0.3', facecolor='white', alpha=0.8),
                     fontsize=9, verticalalignment='bottom')

    if title:
        plt.title(title, fontsize=14, fontweight='bold')
    else:
        plt.title('USD ~ RUB', fontsize=14, fontweight='bold')

    if xlabel:
        plt.xlabel(xlabel)
    else:
        plt.xlabel('Date')

    if ylabel:
        plt.ylabel(ylabel)
    else:
        plt.ylabel('Rub == 1$')

    if show_grid:
        plt.grid(True, alpha=0.3)

    if show_legend and (show_moving_average or show_trend_line or show_confidence_interval):
        plt.legend()

    plt.xticks(rotation=45)
    plt.tight_layout()

    plt.savefig(f"./ds-cepexaaa/code/timeSequences/img/{filename}.png", dpi=300, bbox_inches='tight')
    plt.close()


# ===========================================================================================================================================


def plot_seasonality_analysis(results, filename_prefix="seasonality"):
    plot_autocorrelation(results, f"./ds-cepexaaa/code/timeSequences/img/{filename_prefix}_autocorrelation.png")
    plot_monthly_seasonality(results, f"./ds-cepexaaa/code/timeSequences/img/{filename_prefix}_monthly.png")




def plot_autocorrelation(results, filename):
    plt.figure(figsize=(12, 6))

    autocorr = results['autocorrelation']['autocorrelation']
    significant_periods = results['autocorrelation']['significant_periods']

    plt.stem(range(len(autocorr)), autocorr, basefmt=" ")
    plt.axhline(y=0.15, color='red', linestyle='--', alpha=0.7, label='Threshold of significance (0.15)')
    plt.axhline(y=-0.15, color='red', linestyle='--', alpha=0.7)

    for period in significant_periods[:5]:
        plt.axvline(x=period['lag'], color='orange', linestyle=':', alpha=0.5)
        plt.annotate(f'лаг {period["lag"]}', xy=(period['lag'], period['correlation']),
                     xytext=(5, 5), textcoords='offset points', fontsize=8)

    plt.title('The autocorrelation function (ACF) - search for seasonal periods.', fontweight='bold')
    plt.xlabel('Lag (days)')
    plt.ylabel('autocorrelation')
    plt.legend()
    plt.grid(True, alpha=0.3)
    plt.tight_layout()

    plt.savefig(filename)
    plt.close()



def plot_monthly_seasonality(results, filename):
    plt.figure(figsize=(10, 6))

    monthly_means = results['monthly']['monthly_means']
    months = ['Янв', 'Фев', 'Мар', 'Апр', 'Май', 'Июн',
              'Июл', 'Авг', 'Сен', 'Окт', 'Ноя', 'Дек']

    plt.bar(monthly_means.index, monthly_means.values, alpha=0.7)
    plt.axhline(y=monthly_means.mean(), color='red', linestyle='--',
                label=f'Mean: {monthly_means.mean():.2f}')

    plt.title('Mean value for month', fontweight='bold')
    plt.xlabel('Month')
    plt.ylabel('Course (RUB.)')
    plt.xticks(range(1, 13), months, rotation=45)
    plt.legend()
    plt.grid(True, alpha=0.3)
    plt.tight_layout()

    plt.savefig(filename)
    plt.close()


def create_regression_plots(actual_values, predictions, feature_names, coefficients,
                            filename_prefix="regression"):

    # Plot 1: Actual vs Predicted values
    plt.figure(figsize=(12, 10))

    # Subplot 1: Time series comparison
    plt.subplot(2, 2, 1)
    time_index = range(len(actual_values))
    plt.plot(time_index, actual_values, 'b-', label='Actual', linewidth=2, alpha=0.8)
    plt.plot(time_index, predictions, 'r--', label='Predicted', linewidth=2, alpha=0.8)
    plt.title('USD/RUB Exchange Rate: Actual vs Predicted')
    plt.xlabel('Time Index')
    plt.ylabel('Exchange Rate')
    plt.legend()
    plt.grid(True, alpha=0.3)

    # Subplot 2: Scatter plot of actual vs predicted
    plt.subplot(2, 2, 2)
    plt.scatter(actual_values, predictions, alpha=0.6)
    plt.plot([actual_values.min(), actual_values.max()],
             [actual_values.min(), actual_values.max()], 'r--', linewidth=2)
    plt.title('Predicted vs Actual Values')
    plt.xlabel('Actual Values')
    plt.ylabel('Predicted Values')
    plt.grid(True, alpha=0.3)

    # Subplot 3: Residuals plot
    residuals = actual_values - predictions
    plt.subplot(2, 2, 3)
    plt.scatter(predictions, residuals, alpha=0.6)
    plt.axhline(y=0, color='r', linestyle='--', linewidth=2)
    plt.title('Residuals Plot')
    plt.xlabel('Predicted Values')
    plt.ylabel('Residuals')
    plt.grid(True, alpha=0.3)

    # Subplot 4: Feature importance
    plt.subplot(2, 2, 4)
    # Exclude intercept from feature importance
    feature_importance = np.abs(coefficients[1:]) if len(coefficients) > len(feature_names) else np.abs(coefficients)
    sorted_idx = np.argsort(feature_importance)[::-1]

    # Take top features for readability
    top_n = len(feature_names)
    plt.barh(range(top_n), feature_importance[sorted_idx[:top_n]])
    plt.yticks(range(top_n), [feature_names[i] for i in sorted_idx[:top_n]])
    plt.title('Top Feature Importance (Ridge)')
    plt.xlabel('Absolute Coefficient Value')
    plt.tight_layout()

    plt.savefig(f"./ds-cepexaaa/code/timeSequences/img/{filename_prefix}_analysis.png",
                dpi=300, bbox_inches='tight')
    plt.close()


def create_feature_analysis_plot(feature_names, coefficients, actual_values, predictions,
                                 filename_prefix="feature_analysis"):
    """
    Create detailed feature analysis plot
    """
    plt.figure(figsize=(15, 8))

    # Plot 1: Feature coefficients
    plt.subplot(2, 2, 1)
    feature_coefs = coefficients[1:] if len(coefficients) > len(feature_names) else coefficients
    colors = ['red' if coef < 0 else 'blue' for coef in feature_coefs]
    plt.bar(range(len(feature_coefs)), feature_coefs, color=colors, alpha=0.7)
    plt.axhline(y=0, color='black', linestyle='-', linewidth=1)
    plt.title('Feature Coefficients')
    plt.xlabel('Feature Index')
    plt.ylabel('Coefficient Value')
    plt.xticks(rotation=45)
    plt.grid(True, alpha=0.3)

    # # Plot 2: Prediction error distribution
    # plt.subplot(2, 3, 2)
    errors = actual_values - predictions
    # plt.hist(errors, bins=20, alpha=0.7, edgecolor='black')
    # plt.axvline(x=0, color='red', linestyle='--', linewidth=2)
    # plt.title('Prediction Error Distribution')
    # plt.xlabel('Prediction Error')
    # plt.ylabel('Frequency')
    # plt.grid(True, alpha=0.3)

    # Plot 3: Cumulative prediction accuracy
    plt.subplot(2, 2, 2)
    cumulative_accuracy = 1 - np.abs(errors.cumsum()) / (np.arange(len(errors)) + 1)
    plt.plot(cumulative_accuracy)
    plt.title('Cumulative Prediction Accuracy')
    plt.xlabel('Time Index')
    plt.ylabel('Cumulative Accuracy')
    plt.grid(True, alpha=0.3)

    # Plot 4: Top positive features
    plt.subplot(2, 2, 3)
    positive_coefs = [(name, coef) for name, coef in zip(feature_names, feature_coefs) if coef > 0]
    positive_coefs.sort(key=lambda x: x[1], reverse=True)
    top_positive = positive_coefs[:5]
    if top_positive:
        names, values = zip(*top_positive)
        plt.barh(range(len(names)), values, color='blue', alpha=0.7)
        plt.yticks(range(len(names)), names)
        plt.title('Top Positive Features')
        plt.xlabel('Coefficient Value')

    # Plot 5: Top negative features
    plt.subplot(2, 2, 4)
    negative_coefs = [(name, coef) for name, coef in zip(feature_names, feature_coefs) if coef < 0]
    negative_coefs.sort(key=lambda x: x[1])
    top_negative = negative_coefs[:5]
    if top_negative:
        names, values = zip(*top_negative)
        plt.barh(range(len(names)), values, color='red', alpha=0.7)
        plt.yticks(range(len(names)), names)
        plt.title('Top Negative Features')
        plt.xlabel('Coefficient Value')

    plt.tight_layout()
    plt.savefig(f"./ds-cepexaaa/code/timeSequences/img/{filename_prefix}.png",
                dpi=300, bbox_inches='tight')
    plt.close()
