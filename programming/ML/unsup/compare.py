from sklearn.metrics import confusion_matrix
import numpy as np
import pandas as pd
import matplotlib.pyplot as plt
from sklearn.model_selection import cross_val_score
from sklearn.metrics import accuracy_score, f1_score, silhouette_score, adjusted_rand_score
from sklearn.decomposition import PCA
from sklearn.manifold import TSNE
from sklearn.linear_model import LogisticRegression
from sklearn.preprocessing import Normalizer
from sklearn.svm import SVC
from sklearn.ensemble import RandomForestClassifier
from feature_selection import FilterFeatureSelector, WrapperFeatureSelector, EmbeddedFeatureSelector
from clasterisation import CustomKMeans

def get_top_features(selector, feature_names, X, y, method_name, top_n=30):
    if hasattr(selector, 'fit'):
        selector.fit(X, y)

    indices = []
    scores = []

    if hasattr(selector, 'selected_indices_'):  # Custom Filter
        indices = selector.selected_indices_[:top_n]
        if hasattr(selector, 'scores_'):
            scores = selector.scores_[indices]
        else:
            scores = [1.0] * len(indices)

    elif hasattr(selector, 'get_support'):  # Library selectors
        support = selector.get_support()
        if hasattr(selector, 'scores_'):
            scores = selector.scores_
            indices = np.argsort(scores)[-top_n:][::-1]
        elif hasattr(selector, 'feature_importances_'):
            importances = selector.feature_importances_
            indices = np.argsort(importances)[-top_n:][::-1]
            scores = importances[indices]
        else:
            if hasattr(selector, 'ranking_'):
                # best rank == last deleted
                indices = np.argsort(selector.ranking_)[:top_n]
                scores = 1.0 / selector.ranking_[indices]
            else:
                indices = np.where(support)[0][:top_n]
                scores = [1.0] * len(indices)

    elif hasattr(selector, 'coef_'):  # Embedded with coef
        # importances = np.abs(selector.coef_).flatten()
        if selector.coef_.ndim > 1:
            importances = np.mean(np.abs(selector.coef_), axis=0)
        else:
            importances = np.abs(selector.coef_)
        indices = np.argsort(importances)[-top_n:][::-1]
        scores = importances[indices]

    else:
        print(f"Cannot extract features from {method_name}")
        return [], []

    top_features = []
    for i, (idx, score) in enumerate(zip(indices, scores), 1):
        if idx < len(feature_names):
            feature = feature_names[idx]
            top_features.append(feature)
            print(f"{i:2d}. {feature:30} | score = {score:.4f}")

    return top_features, indices


def compare_feature_selection_methods(data):
    filter_selector = FilterFeatureSelector(k=500)
    wrapper_selector = WrapperFeatureSelector(
        estimator=LogisticRegression(random_state=42),
        n_features_to_select=300
    )
    embedded_selector = EmbeddedFeatureSelector(C=0.5, threshold=1e-5, max_iter=5000)
    methods = {}
    methods['Filter'] = filter_selector
    methods['Wrapper'] = wrapper_selector
    methods['Embedded'] = embedded_selector

    all_top_features = {}
    for method_name, selector in methods.items():
        print("\nTop 30 features in " + method_name)
        top_features, indices = get_top_features(selector, data['feature_names'], data['X_train'], data['y_train'],  method_name, 30)
        all_top_features[method_name] = set(top_features)

    _compare_feature_sets(all_top_features)

    return methods


def _compare_feature_sets(feature_sets):
    methods = list(feature_sets.keys())
    n_methods = len(methods)

    overlap_matrix = np.zeros((n_methods, n_methods), dtype=int)

    for i in range(n_methods):
        for j in range(n_methods):
            overlap = len(feature_sets[methods[i]] & feature_sets[methods[j]])
            overlap_matrix[i, j] = overlap

    print("\nOverlap matrix (number of common features):")
    for method in methods:
        print(f"{method[:8]:>8}", end=" ")
    print()

    for i, method_i in enumerate(methods):
        print(f"{method_i[:15]:<15}", end="")
        for j in range(n_methods):
            print(f"{overlap_matrix[i, j]:>8}", end=" ")
        print()

    common_features = set.intersection(*feature_sets.values())
    print(f"\nFeatures common to ALL methods ({len(common_features)}):")
    for feature in sorted(list(common_features))[:20]:
        print(f"  - {feature}")

    print("\nUnique features per method:")
    for method, features in feature_sets.items():
        other_features = set.union(*[fs for m, fs in feature_sets.items() if m != method])
        unique = features - other_features
        print(f"\n{method} - {len(unique)} unique features:")
        for feature in sorted(list(unique))[:10]:
            print(f"  - {feature}")


def evaluate_classifiers_with_feature_selection(data, methods):
    classifiers = {
        'Logistic Regression': LogisticRegression(random_state=42, max_iter=1000),
        'SVM': SVC(kernel='linear', random_state=42, probability=True),
        'Random Forest': RandomForestClassifier(n_estimators=100, random_state=42)
    }

    results = []

    print(f"\n{'=' * 60}")
    print("CLASSIFIER PERFORMANCE - ORIGINAL FEATURES")

    for clf_name, clf in classifiers.items():
        clf.fit(data['X_train'], data['y_train'])
        y_pred = clf.predict(data['X_test'])
        acc = accuracy_score(data['y_test'], y_pred)
        f1 = f1_score(data['y_test'], y_pred, average='weighted')

        cv_scores = cross_val_score(clf, data['X_train'], data['y_train'], cv=5, scoring='f1_weighted')

        results.append({
            'method': 'Original',
            'classifier': clf_name,
            'n_features': data['X_train'].shape[1],
            'test_accuracy': acc,
            'test_f1': f1,
            'cv_f1_mean': cv_scores.mean(),
            'cv_f1_std': cv_scores.std()
        })

        print(f"{clf_name:20} | Acc: {acc:.4f} | F1: {f1:.4f} | CV F1: {cv_scores.mean():.4f} ± {cv_scores.std():.4f}")

    for fs_name, selector in methods.items():
        print(f"CLASSIFIER PERFORMANCE - {fs_name}")

        X_train_selected = selector.fit_transform(data['X_train'], data['y_train'])
        X_test_selected = selector.transform(data['X_test'])
        n_features = X_train_selected.shape[1]

        print(f"Selected {n_features} features")

        for clf_name, clf in classifiers.items():
            clf.fit(X_train_selected, data['y_train'])
            y_pred = clf.predict(X_test_selected)
            acc = accuracy_score(data['y_test'], y_pred)
            f1 = f1_score(data['y_test'], y_pred, average='weighted')

            cv_scores = cross_val_score(clf, X_train_selected, data['y_train'], cv=5, scoring='f1_weighted')

            results.append({
                'method': fs_name,
                'classifier': clf_name,
                'n_features': n_features,
                'test_accuracy': acc,
                'test_f1': f1,
                'cv_f1_mean': cv_scores.mean(),
                'cv_f1_std': cv_scores.std()
            })

            print(
                f"{clf_name:20} | Acc: {acc:.4f} | F1: {f1:.4f} | CV F1: {cv_scores.mean():.4f} ± {cv_scores.std():.4f}")

    results_df = pd.DataFrame(results)

    print(f"\n{'=' * 80}")
    print("SUMMARY - BEST METHOD SELECTION")

    summary = results_df.groupby('method').agg({
        'test_f1': 'mean',
        'cv_f1_mean': 'mean',
        'n_features': 'first'
    }).round(4)

    summary['compression'] = (data['X_train'].shape[1] / summary['n_features']).round(1)

    print("\nComparison of feature selection methods:")
    print(f"{'Method':<20} {'Test F1':<10} {'CV F1':<10} {'Features':<10} {'Compression':<10}")
    print("-" * 60)
    for idx, row in summary.iterrows():
        print(f"{idx:<20} {row['test_f1']:<10.4f} {row['cv_f1_mean']:<10.4f} "
              f"{int(row['n_features']):<10} {row['compression']:<10.1f}x")

    best_method = summary['cv_f1_mean'].idxmax()
    print(f"\n RECOMMENDED METHOD: {best_method}")
    print(f"   - Average CV F1: {summary.loc[best_method, 'cv_f1_mean']:.4f}")
    print(f"   - Feature compression: {summary.loc[best_method, 'compression']:.1f}x")

    return results_df, best_method


def cluster_and_evaluate(data, best_selector):

    print(f"\n{'=' * 60}")
    print("CLUSTERING EVALUATION")

    best_selector.fit(data['X_train'], data['y_train'])

    X_all_original = data['X']
    X_all_selected = best_selector.transform(data['X'])

    datasets = {
        'Original': X_all_original,
        f'Selected ({X_all_selected.shape[1]} features)': X_all_selected
    }

    results = []

    for name, X_data in datasets.items():
        print(f"\nClustering on {name} data...")

        if hasattr(X_data, 'toarray'):
            X_dense = X_data.toarray()
        else:
            X_dense = X_data

        normalizer = Normalizer()
        X_norm = normalizer.fit_transform(X_dense)

        kmeans = CustomKMeans(random_state=9)
        cluster_labels = kmeans.fit_predict(X_norm)

        # Внутренняя мера: Силуэтный коэффициент
        silhouette = silhouette_score(X_dense, cluster_labels)

        # Внешняя мера: Adjusted Rand Index
        ari = adjusted_rand_score(data['y'], cluster_labels)

        results.append({
            'dataset': name,
            'n_features': X_dense.shape[1],
            'silhouette': silhouette,
            'ari': ari,
            'cluster_labels': cluster_labels,
            'true_labels': data['y'],
            'X_data': X_dense
        })

        print(f"  Silhouette score: {silhouette:.4f}")
        print(f"  Adjusted Rand Index: {ari:.4f}")

        # Analysis of cluster correspondence to classes
        cm = confusion_matrix(data['y'], cluster_labels)
        print(f"  Confusion matrix:\n{cm}")

    return results


def dimensionality_reduction_and_visualization(data, clustering_results):
    fig, axes = plt.subplots(2, 4, figsize=(20, 10))

    for idx, result in enumerate(clustering_results):
        dataset_name = result['dataset']
        true_labels = result['true_labels']
        cluster_labels = result['cluster_labels']
        X_dense = result['X_data']

        # PCA
        pca = PCA(n_components=2, random_state=42)
        X_pca = pca.fit_transform(X_dense)

        # t-SNE
        tsne = TSNE(n_components=2, random_state=42, perplexity=30)
        X_tsne = tsne.fit_transform(X_dense)

        ax = axes[0, idx * 2]
        scatter = ax.scatter(X_pca[:, 0], X_pca[:, 1], c=true_labels,
                             cmap='tab10', alpha=0.6, edgecolors='k', s=50)
        ax.set_title(f'{dataset_name}\nPCA - True Classes')
        ax.set_xlabel(f'PC1 ({pca.explained_variance_ratio_[0]:.1%})')
        ax.set_ylabel(f'PC2 ({pca.explained_variance_ratio_[1]:.1%})')

        ax = axes[0, idx * 2 + 1]
        scatter = ax.scatter(X_pca[:, 0], X_pca[:, 1], c=cluster_labels,
                             cmap='tab10', alpha=0.6, edgecolors='k', s=50)
        ax.set_title(f'{dataset_name}\nPCA - Clusters')
        ax.set_xlabel(f'PC1 ({pca.explained_variance_ratio_[0]:.1%})')
        ax.set_ylabel(f'PC2 ({pca.explained_variance_ratio_[1]:.1%})')

        ax = axes[1, idx * 2]
        scatter = ax.scatter(X_tsne[:, 0], X_tsne[:, 1], c=true_labels,
                             cmap='tab10', alpha=0.6, edgecolors='k', s=50)
        ax.set_title(f'{dataset_name}\nt-SNE - True Classes')
        ax.set_xlabel('t-SNE 1')
        ax.set_ylabel('t-SNE 2')

        ax = axes[1, idx * 2 + 1]
        scatter = ax.scatter(X_tsne[:, 0], X_tsne[:, 1], c=cluster_labels,
                             cmap='tab10', alpha=0.6, edgecolors='k', s=50)
        ax.set_title(f'{dataset_name}\nt-SNE - Clusters')
        ax.set_xlabel('t-SNE 1')
        ax.set_ylabel('t-SNE 2')

        print(f"\n{dataset_name}:")
        print(f"  PCA explained variance: {pca.explained_variance_ratio_.sum():.2%}")

    plt.tight_layout()
    plt.savefig("./unsup-cepexaaa/img/reduction.png")
    plt.close()

    _visualize_cluster_comparison(clustering_results)


def _visualize_cluster_comparison(clustering_results):
    fig, axes = plt.subplots(1, 2, figsize=(15, 6))

    for idx, result in enumerate(clustering_results):
        dataset_name = result['dataset']
        unique, counts = np.unique(result['cluster_labels'], return_counts=True)

        axes[idx].bar(unique, counts, color=['skyblue', 'lightcoral'])
        axes[idx].set_title(f'Cluster Sizes - {dataset_name}')
        axes[idx].set_xlabel('Cluster')
        axes[idx].set_ylabel('Number of Texts')
        axes[idx].set_xticks(unique)

        for i, count in enumerate(counts):
            axes[idx].text(unique[i], count + 0.5, str(count),
                           ha='center', va='bottom')

    plt.tight_layout()
    plt.savefig("./unsup-cepexaaa/img/comparison.png")
    plt.close()

