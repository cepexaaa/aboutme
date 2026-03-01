from download_data import download, aggregation
from paint import toImg, plot_seasonality_analysis, create_regression_plots, create_feature_analysis_plot
from season_analyze import analyze_seasonality, analyze_seasonality_complete
from stationarity import analyze_transformations, plot_transformations, TimeSeriesTransformer
import pandas as pd
from linera_regression import RidgeLinearRegressor, prepare_regression_data, evaluate_regression_model
from timeSeriesModels import TimeSeriesModelComparator, NaiveMeanModel, NaiveLastValueModel, HoltWintersModel, SARIMAModel


# Download once, because I don't want to wait while it will download data again
# data_frame = download()

data_frame = pd.read_csv("./ds-cepexaaa/code/timeSequences/data/rub-usd.csv")
aggregated_data = pd.read_csv("./ds-cepexaaa/code/timeSequences/data/aggregated_rub-usd.csv") # aggregation(data_frame)

toImg(data_frame, "courses", show_moving_average=True)
toImg(aggregated_data, "aggregated", show_moving_average=True)


print("===================== Step 2 ===========================")

toImg(aggregated_data, "trend", show_trend_line=True)
analyze_seasonality(data_frame)
results = analyze_seasonality_complete(data_frame)
plot_seasonality_analysis(results, "usd_rub_seasonality")

# Плавное убывание автокорреляционной функции к нулю без выраженных пиков говорит о том,
# что в динамике курса USD/RUB преобладают трендовые компоненты и случайные колебания над периодичными изменениями.
# Мы должны были увидеть на графике значительные пики на определенных лагах, если бы имела место сезонность.
# Такие пики указывали бы на то, что значения курса через эти промежутки времени статистически связаны между собой, формируя повторяющиеся циклы.
#
# Однако отсутствие таких выраженных пиков означает, что корреляция между значениями курса быстро затухает со временем,
# и не существует устойчивых периодических паттернов, которые бы систематически повторялись.
# Это характерно для финансовых временных рядов.
#
# Слабая годовая сезонность с силой 0.27, хоть и является наиболее заметной среди всех проверенных периодов,
# всё равно недостаточно сильна, чтобы считать, что есть годовая сезонность.


print("===================== Step 3 ===========================")

results = analyze_transformations(aggregated_data)

plot_transformations(results, "usd_rub_stationarity")

transformer = TimeSeriesTransformer()
stationary_series = transformer.apply_transformations(aggregated_data)

transformation_info = transformer.get_transformation_info()
print(f"   Transforms: {transformation_info['transformations']}")

if transformation_info['needs_inverse']:
    test_data = stationary_series.head(10)
    restored_data = transformer.inverse_transform(test_data)
    print(f"   Converted data: {test_data.values[:3]}...")
    print(f"   Recovered data: {restored_data.values[:3]}...")


print("===================== Step 4 ===========================")

regression_data = prepare_regression_data(aggregated_data, 'usd_rub')

X_train = regression_data['X_train']
X_test = regression_data['X_test']
y_train = regression_data['y_train']
y_test = regression_data['y_test']
feature_names = regression_data['feature_names']

print(f"Dataset Info:")
print(f"Training samples: {X_train.shape[0]}")
print(f"Test samples: {X_test.shape[0]}")
print(f"Features: {len(feature_names)}")
print(f"Feature names: {feature_names}")

model = RidgeLinearRegressor()
model.fit(X_train, y_train)

results = evaluate_regression_model(model, X_test, y_test)

create_regression_plots(actual_values=y_test.values, predictions=results['predictions'], feature_names=feature_names, coefficients=model.coef_, filename_prefix="regression")
create_feature_analysis_plot(feature_names=feature_names,coefficients=model.coef_,actual_values=y_test.values,predictions=results['predictions'],filename_prefix="feature_analysis")

print("\nFeature Insights:")
feature_coefs = model.coef_[1:] if len(model.coef_) > len(feature_names) else model.coef_
for name, coef in zip(feature_names, feature_coefs):
    direction = "increases" if coef > 0 else "decreases"
    print(f"  {name}: {coef:.4f} ({direction} predicted rate)")

print(f"Model performance (Regression analysis): R² = {results['r2']:.4f}")


print("===================== Step 5+6 ===========================")

series_data = aggregated_data['usd_rub']
split_idx = int(len(series_data) * 0.8)
train_data = series_data[:split_idx]
test_data = series_data[split_idx:]

comparator = TimeSeriesModelComparator()

# Add naive models (baselines)
comparator.add_model("Naive_Mean", NaiveMeanModel())
comparator.add_model("Naive_Last", NaiveLastValueModel())
comparator.add_model("Holt_Winters", HoltWintersModel(seasonal_periods=52))
comparator.add_model("SARIMA", SARIMAModel(order=(1, 1, 1), seasonal_order=(1, 1, 1, 52)))
results = comparator.train_and_evaluate(train_data, test_data)
comparator.plot_comparison("time_series_forecasting_comparison")

comparator.print_summary()

print(f"\nData Characteristics:")
print(f"  Training mean: {train_data.mean():.2f}")
print(f"  Test mean: {test_data.mean():.2f}")
print(f"  Overall trend: {'Increasing' if series_data.iloc[-1] > series_data.iloc[0] else 'Decreasing'}")
print(f"  Volatility (std): {series_data.std():.2f}")














