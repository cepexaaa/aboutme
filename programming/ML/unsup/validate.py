import numpy as np
from sklearn.model_selection import KFold, cross_val_score, cross_validate
from sklearn.metrics import f1_score, precision_score, recall_score, accuracy_score, make_scorer
from sklearn.base import clone
import pandas as pd
from typing import Dict, Any




def evaluate_model_with_metrics(model, X_test: Any, y_test: np.ndarray) -> Dict[str, float]:
    y_pred = model.predict(X_test)
    metrics = {
        'accuracy': accuracy_score(y_test, y_pred),
        'precision': precision_score(y_test, y_pred, average='weighted'),
        'recall': recall_score(y_test, y_pred, average='weighted'),
        'f1': f1_score(y_test, y_pred, average='weighted')
    }
    return metrics


def kfold_cross_validation(
        model,
        X: Any,
        y: np.ndarray,
        n_splits: int = 5,
        random_state: int = 42,
        verbose: bool = True
) -> Dict[str, Any]:
    kf = KFold(n_splits=n_splits, shuffle=True, random_state=random_state)

    scoring = {
        'accuracy': 'accuracy',
        'precision': make_scorer(precision_score, average='weighted', zero_division=0),
        'recall': make_scorer(recall_score, average='weighted', zero_division=0),
        'f1': make_scorer(f1_score, average='weighted', zero_division=0)
    }

    cv_results = cross_validate(
        model, X, y,
        cv=kf,
        scoring=scoring,
        return_train_score=True,
        n_jobs=-1
    )

    results = {
        'cv_splits': n_splits,
        'test_accuracy_mean': cv_results['test_accuracy'].mean(),
        'test_accuracy_std': cv_results['test_accuracy'].std(),
        'test_precision_mean': cv_results['test_precision'].mean(),
        'test_precision_std': cv_results['test_precision'].std(),
        'test_recall_mean': cv_results['test_recall'].mean(),
        'test_recall_std': cv_results['test_recall'].std(),
        'test_f1_mean': cv_results['test_f1'].mean(),
        'test_f1_std': cv_results['test_f1'].std(),
        'train_accuracy_mean': cv_results['train_accuracy'].mean(),
        'train_f1_mean': cv_results['train_f1'].mean(),
        'fold_details': []
    }

    # Info for each folds
    for fold_idx, (train_idx, val_idx) in enumerate(kf.split(X)):
        X_train_fold = X[train_idx] if hasattr(X, 'shape') else X.iloc[train_idx]
        X_val_fold = X[val_idx] if hasattr(X, 'shape') else X.iloc[val_idx]
        y_train_fold = y[train_idx]
        y_val_fold = y[val_idx]

        fold_model = clone(model)
        fold_model.fit(X_train_fold, y_train_fold)

        y_pred_train = fold_model.predict(X_train_fold)
        y_pred_val = fold_model.predict(X_val_fold)

        fold_metrics = {
            'fold': fold_idx + 1,
            'train_size': len(y_train_fold),
            'val_size': len(y_val_fold),
            'train_accuracy': accuracy_score(y_train_fold, y_pred_train),
            'val_accuracy': accuracy_score(y_val_fold, y_pred_val),
            'val_precision': precision_score(y_val_fold, y_pred_val, average='weighted', zero_division=0),
            'val_recall': recall_score(y_val_fold, y_pred_val, average='weighted', zero_division=0),
            'val_f1': f1_score(y_val_fold, y_pred_val, average='weighted', zero_division=0)
        }

        results['fold_details'].append(fold_metrics)

    if verbose:
        print(f"K-FOLD CROSS-VALIDATION (K={n_splits})\n")

        print(f"\n{'Metric':<20} {'Mean':<10} {'Std':<10}")
        print("-" * 40)
        print(f"{'Test Accuracy':<20} {results['test_accuracy_mean']:.4f}     ±{results['test_accuracy_std']:.4f}")
        print(f"{'Test Precision':<20} {results['test_precision_mean']:.4f}     ±{results['test_precision_std']:.4f}")
        print(f"{'Test Recall':<20} {results['test_recall_mean']:.4f}     ±{results['test_recall_std']:.4f}")
        print(f"{'Test F1-score':<20} {results['test_f1_mean']:.4f}     ±{results['test_f1_std']:.4f}")
        print(f"{'Train Accuracy':<20} {results['train_accuracy_mean']:.4f}")
        print(f"{'Train F1-score':<20} {results['train_f1_mean']:.4f}")

        print(f"\nDetails per fold:")
        print("-" * 60)
        print(f"{'Fold':<6} {'Train':<8} {'Val':<8} {'Val F1':<10} {'Val Acc':<10}")
        print("-" * 60)
        for fold in results['fold_details']:
            print(f"{fold['fold']:<6} {fold['train_size']:<8} {fold['val_size']:<8} "
                  f"{fold['val_f1']:<10.4f} {fold['val_accuracy']:<10.4f}")

    return results


def compare_feature_selection_methods_with_cv(
        original_data: Dict[str, Any],
        selectors: Dict[str, Any],
        base_model,
        n_splits: int = 5
) -> pd.DataFrame:
    X_train = original_data['X_train']
    y_train = original_data['y_train']
    X_test = original_data['X_test']
    y_test = original_data['y_test']

    results = []

    print("\n" + "=" * 60)
    print("EVALUATING ORIGINAL FEATURES")

    cv_results_original = kfold_cross_validation(
        clone(base_model), X_train, y_train,
        n_splits=n_splits, verbose=True
    )

    model_original = clone(base_model)
    model_original.fit(X_train, y_train)
    test_metrics = evaluate_model_with_metrics(model_original, X_test, y_test)

    results.append({
        'method': 'Original',
        'n_features': X_train.shape[1],
        'cv_f1_mean': cv_results_original['test_f1_mean'],
        'cv_f1_std': cv_results_original['test_f1_std'],
        'cv_accuracy_mean': cv_results_original['test_accuracy_mean'],
        'test_f1': test_metrics['f1'],
        'test_accuracy': test_metrics['accuracy'],
        'test_precision': test_metrics['precision'],
        'test_recall': test_metrics['recall']
    })

    for method_name, selector in selectors.items():
        print(f"\n" + "=" * 60)
        print(f"EVALUATING {method_name.upper()} METHOD")

        if hasattr(selector, 'fit_transform'):
            X_train_selected = selector.fit_transform(X_train, y_train)
            X_test_selected = selector.transform(X_test)
        else:
            selector.fit(X_train, y_train)
            X_train_selected = selector.transform(X_train)
            X_test_selected = selector.transform(X_test)

        n_features = X_train_selected.shape[1]
        print(f"Selected {n_features} features")

        cv_results = kfold_cross_validation(
            clone(base_model), X_train_selected, y_train,
            n_splits=n_splits, verbose=True
        )

        model_selected = clone(base_model)
        model_selected.fit(X_train_selected, y_train)
        test_metrics = evaluate_model_with_metrics(model_selected, X_test_selected, y_test)

        results.append({
            'method': method_name,
            'n_features': n_features,
            'cv_f1_mean': cv_results['test_f1_mean'],
            'cv_f1_std': cv_results['test_f1_std'],
            'cv_accuracy_mean': cv_results['test_accuracy_mean'],
            'test_f1': test_metrics['f1'],
            'test_accuracy': test_metrics['accuracy'],
            'test_precision': test_metrics['precision'],
            'test_recall': test_metrics['recall']
        })

    comparison_df = pd.DataFrame(results)

    print("\n" + "=" * 80)
    print("FEATURE SELECTION METHODS COMPARISON")
    print("=" * 80)
    print(f"\n{'Method':<15} {'Features':<10} {'CV F1':<12} {'Test F1':<12} {'Test Acc':<12}")
    print("-" * 80)

    for _, row in comparison_df.iterrows():
        print(f"{row['method']:<15} {row['n_features']:<10} "
              f"{row['cv_f1_mean']:.4f}±{row['cv_f1_std']:.4f}  "
              f"{row['test_f1']:.4f}        {row['test_accuracy']:.4f}")

    return comparison_df


