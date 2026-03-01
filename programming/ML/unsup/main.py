from unpack_data import load_and_prepare_data
from runner import run_my_methods, run_lib_methods
from compare import compare_feature_selection_methods, evaluate_classifiers_with_feature_selection, cluster_and_evaluate, dimensionality_reduction_and_visualization

data = load_and_prepare_data()
run_my_methods(data)
run_lib_methods(data)

methods = compare_feature_selection_methods(data)

print("\n3. EVALUATING CLASSIFIERS...")
results_df, best_method = evaluate_classifiers_with_feature_selection(data, methods)

print("\n4. CLUSTERING ANALYSIS...")
best_selector = methods.get(best_method, list(methods.values())[0])
clustering_results = cluster_and_evaluate(data, best_selector)

print("\n5. DIMENSIONALITY REDUCTION AND VISUALIZATION...")
dimensionality_reduction_and_visualization(data, clustering_results)


# Для валидации алгоритма буду использовать *K-Fold CV*, K-блочная кросс-валидация
# В качество классификации буду оценивать на метрике F1-score, так как это гармоническое среднее между Precision и Recall. Это будет достаточно полно отображать результаты

# Работу выполнил Кубеш Сергей из M3335