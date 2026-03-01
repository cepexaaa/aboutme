import matplotlib.pyplot as plt
import pandas as pd
import numpy as np
from statsmodels.tsa.seasonal import seasonal_decompose
from statsmodels.tsa.stattools import acf
from paint import toImg
from scipy import stats

def analyze_seasonality(data, filename_prefix="seasonality"):
    df = data.copy()
    if 'date' in df.columns:
        df['date'] = pd.to_datetime(df['date'])
        df.set_index('date', inplace=True)

    if 'usd_rub' in df.columns:
        value_col = 'usd_rub'
    else:
        value_col = df.columns[0]
    months = ['Янв', 'Фев', 'Мар', 'Апр', 'Май', 'Июн', 'Июл', 'Авг', 'Сен', 'Окт', 'Ноя', 'Дек']

    try:
        df['year'] = df.index.year
        df['month'] = df.index.month

        heatmap_data = df.pivot_table(
            values=value_col,
            index='year',
            columns='month',
            aggfunc='mean'
        )

        plt.figure(figsize=(12, 8))
        plt.imshow(heatmap_data, aspect='auto', cmap='RdYlGn_r')
        plt.colorbar(label='RUB')
        plt.title('heatmap data of season USD/RUB', fontweight='bold')
        plt.xlabel('Month')
        plt.ylabel('Year')
        plt.xticks(range(12), months, rotation=45)
        plt.yticks(range(len(heatmap_data.index)), heatmap_data.index)

        for i in range(len(heatmap_data.index)):
            for j in range(len(heatmap_data.columns)):
                plt.text(j, i, f'{heatmap_data.iloc[i, j]:.1f}',
                         ha='center', va='center', fontsize=8)

        plt.tight_layout()
        plt.savefig(f"./ds-cepexaaa/code/timeSequences/img/{filename_prefix}_heatmap.png")
        plt.close()

    except Exception as e:
        print(f"Error while make heat map: {e}")


    try:
        years = sorted(df['year'].unique())
        plt.figure(figsize=(12, 8))

        for year in years:
            year_data = df[df['year'] == year]
            monthly_avg = year_data.groupby('month')[value_col].mean()
            plt.plot(monthly_avg.index, monthly_avg.values, marker='o', label=str(year))

        plt.title('Seasons patterns by years', fontweight='bold')
        plt.xlabel('Month')
        plt.ylabel('Rub')
        plt.xticks(range(1, 13), months, rotation=45)
        plt.legend()
        plt.grid(True, alpha=0.3)
        plt.tight_layout()
        plt.savefig(f"./ds-cepexaaa/code/timeSequences/img/{filename_prefix}_yearly_patterns.png")
        plt.close()

    except Exception as e:
        print(f"Error while make seasons graphics: {e}")

# =============================================================================================================================================================


def calculate_seasonal_decomposition(data, value_col='u-r'):
    periods_to_try = [7, 30, 90, 365]
    best_decomposition = None
    best_period = None
    best_strength = 0

    for period in periods_to_try:
        try:
            if len(data) < 2 * period:
                print(f"   Period {period}: missed (few data count)")
                continue

            decomposition = seasonal_decompose(
                data[value_col].dropna(),
                model='additive',
                period=period
            )

            strength = np.std(decomposition.seasonal.dropna()) / np.std(decomposition.observed)
            print(f"   Period {period}: strength season = {strength:.4f}")

            if strength > best_strength:
                best_strength = strength
                best_decomposition = decomposition
                best_period = period

        except Exception as e:
            print(f"   Period {period}: error - {e}")

    if best_decomposition is not None:
        print(f"The best Period: {best_period} days (strength: {best_strength:.4f})")
        return {
            'observed': best_decomposition.observed,
            'trend': best_decomposition.trend,
            'seasonal': best_decomposition.seasonal,
            'residual': best_decomposition.resid,
            'seasonal_strength': best_strength,
            'period_used': best_period,
            'all_periods_tested': periods_to_try
        }
    else:
        print("No available decompositions")
        return None

def calculate_autocorrelation_analysis(data, value_col='usd_rub', max_lag=400):
    data_clean = data[value_col].dropna()
    autocorr_values = acf(data_clean, nlags=min(max_lag, len(data_clean) // 4))

    significant_periods = []
    threshold = 0.15  # Порог значимости

    for lag in range(1, len(autocorr_values)):
        if abs(autocorr_values[lag]) > threshold:
            significant_periods.append({
                'lag': lag,
                'correlation': autocorr_values[lag]
            })

    return {
        'autocorrelation': autocorr_values,
        'significant_periods': significant_periods[:10],
        'max_correlation': max(autocorr_values[1:]) if len(autocorr_values) > 1 else 0
    }


def calculate_monthly_seasonality(data, value_col='usd_rub'):
    df = data.copy()
    if 'date' in df.columns:
        df['date'] = pd.to_datetime(df['date'])
        df.set_index('date', inplace=True)

    df['month'] = df.index.month
    df['year'] = df.index.year

    # mean in months
    monthly_stats = df.groupby('month')[value_col].agg(['mean', 'std', 'count'])

    monthly_amplitude = monthly_stats['mean'].max() - monthly_stats['mean'].min()

    return {
        'monthly_means': monthly_stats['mean'],
        'monthly_std': monthly_stats['std'],
        'amplitude': monthly_amplitude,
        'relative_amplitude': monthly_amplitude / monthly_stats['mean'].mean()
    }


def analyze_seasonality_complete(data, value_col='usd_rub'):
    results = {'decomposition': calculate_seasonal_decomposition(data, value_col),
               'autocorrelation': calculate_autocorrelation_analysis(data, value_col),
               'monthly': calculate_monthly_seasonality(data, value_col)}

    return results
