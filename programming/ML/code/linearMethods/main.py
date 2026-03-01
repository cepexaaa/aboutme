import pandas as pd
import numpy as np
from sklearn.preprocessing import LabelEncoder
from sklearn.model_selection import train_test_split
from sklearn.metrics import accuracy_score, classification_report, confusion_matrix, precision_score, recall_score, f1_score, roc_auc_score
import matplotlib.pyplot as plt
import seaborn as sns
from algorythm import RidgeLinearClassifier, GradientDescentLinearClassifier
from hyperparameters import HyperparameterOptimizer

df = pd.read_csv('data.csv')

X = df.iloc[:, 2:]  # another signs
y = df.iloc[:, 1]   # Spectral class

# Convert spectral_class in numerical
label_encoder = LabelEncoder()
y_encoded = label_encoder.fit_transform(y)

print("Spectral_class:", label_encoder.classes_)
print("Their codes:", np.unique(y_encoded))
print("Class distribution:")
print(pd.Series(y_encoded).value_counts())

# hot stars (A, B) vs cool stars (G, K, M, F)
hot_stars = ['A', 'B']

y_binary = np.where(y.isin(hot_stars), 1, 0)  # 1 - hot, 0 - cool

print("Distribution in binary classification:")
print(f"Hot stars (1): {np.sum(y_binary)}")
print(f"Coll stars (0): {len(y_binary) - np.sum(y_binary)}")

X_train, X_test, y_train, y_test = train_test_split(
    X, y_binary, test_size=0.3, random_state=42, stratify=y_binary
)

print(f"Size of train selection: {X_train.shape}")
print(f"Size of test selection: {X_test.shape}")

def evaluate_model(y_true, y_pred, y_pred_proba=None):
    accuracy = accuracy_score(y_true, y_pred) # the proportion of correct predictions
    precision = precision_score(y_true, y_pred) # how much of true is true = TP / (TP + FP)
    recall = recall_score(y_true, y_pred) # TP / (TP + FN)
    f1 = f1_score(y_true, y_pred) # harmonic mean

    print(f"Accuracy: {accuracy:.4f}")
    print(f"Precision: {precision:.4f}")
    print(f"Recall: {recall:.4f}")
    print(f"F1-Score: {f1:.4f}")

    if y_pred_proba is not None:
        auc = roc_auc_score(y_true, y_pred_proba)
        print(f"AUC-ROC: {auc:.4f}")

    return accuracy, precision, recall, f1

feature_df = X.copy()
feature_df['target'] = y_binary

correlations = []
for col in feature_df.columns[:-1]:  # all column except target
    if col != 'target':  # extra check
        corr_value = feature_df[col].corr(feature_df['target'])
        correlations.append((col, abs(corr_value), corr_value))  # (name, abs correlation, correlation)

correlations.sort(key=lambda x: x[1], reverse=True)

print("10 Most powerful correlation:")
for col_name, _, corr in correlations[:10]:
    print(f"{col_name}: corr = {corr:.4f}")

top_n = 10
selected_features = [col_name for col_name, _, _ in correlations[:top_n]]

print(f"\nWas chosen {len(selected_features)} signs: {selected_features}")

# Тут по-русски:
# Есть очень много признаков, которые как-то коррелируют с целевым признаком.
# Надо Взять самые весомые признаки, по которым лучше всего можно предсказывать.
# Выше я выбрал такие признаки и я использую выборку на обучение и тестирование
# из тех признаков, которые будут лучше помогать в предсказании.
X_train_selected = X_train[selected_features]
X_test_selected = X_test[selected_features]



optimizer = HyperparameterOptimizer(cv_folds=5, random_state=42)

# Optimize Ridge Classifier
ridge_params, ridge_score = optimizer.optimize_ridge_classifier(X_train_selected, y_train)
print(f"\nBest Ridge parameters: {ridge_params}")
print(f"Best Ridge CV score: {ridge_score:.4f}")

# Optimize Gradient Classifier
gradient_params, gradient_score = optimizer.optimize_gradient_classifier(X_train_selected, y_train)
print(f"\nBest Gradient parameters: {gradient_params}")
print(f"Best Gradient CV score: {gradient_score:.4f}")

model_best = ""

if ridge_score > gradient_score:
    print("Best model: RidgeLinearClassifier")
    best_model = RidgeLinearClassifier(**ridge_params)
    model_best = "RidgeLinearClassifier"
else:
    print("Best model: GradientDescentLinearClassifier")
    best_model = GradientDescentLinearClassifier(**gradient_params)
    model_best = "GradientDescentLinearClassifier"

# Train and evaluate best model
print("\nTraining best model...")
best_model.fit(X_train_selected, y_train)

y_pred = best_model.predict(X_test_selected)
y_pred_proba = best_model.predict_proba(X_test_selected)[:, 1]

print("\nBest Model Performance on Test Set:")
evaluate_model(y_test, y_pred, y_pred_proba)

# Plot confusion matrix
cm = confusion_matrix(y_test, y_pred)
plt.figure(figsize=(8, 6))
sns.heatmap(cm, annot=True, fmt='d', cmap='Blues')
plt.title(f'Confusion Matrix - {model_best}')
plt.ylabel('True classes')
plt.xlabel('Predicted classes')
plt.savefig("./grafics/confusion_matrix_best_model.png")
plt.close()

# Learning curves for gradient descent classifier
print("\n" + "="*50)
print("LEARNING CURVES ANALYSIS")


# =====================



# Use gradient descent for learning curves since it has iterative training
gd_model = GradientDescentLinearClassifier(
    loss='logistic', learning_rate=0.01, alpha=0.001,
    l1_ratio=0.5, max_iter=1000, random_state=42, verbose=False
)

train_losses = gd_model.fit(X_train_selected, y_train, X_val=X_test_selected, y_val=y_test).loss_history_
test_losses = gd_model.val_loss_history_

def smooth_curve(losses, window_size=10):
    if len(losses) < window_size:
        return losses

    smoothed = []
    for i in range(len(losses)):
        if i < window_size:
            smoothed.append(np.mean(losses[:i + 1]))
        else:
            smoothed.append(np.mean(losses[i - window_size + 1:i + 1]))
    return smoothed


smoothed_train_losses = smooth_curve(train_losses, window_size=20)

plt.figure(figsize=(12, 5))

plt.subplot(1, 2, 1)
plt.plot(smoothed_train_losses, label='Training Loss', linewidth=2)
plt.xlabel('Iteration')
plt.ylabel('Empirical Risk (Logistic Loss)')
plt.title('Training Learning Curve\n(Smoothed Empirical Risk)')
plt.legend()
plt.grid(True, alpha=0.3)

plt.subplot(1, 2, 2)
plt.plot(test_losses, label='Test Loss', linewidth=2, color='orange')
plt.xlabel('Iteration')
plt.ylabel('Test Loss')
plt.title('Test Learning Curve\n(Generalization Performance)')
plt.legend()
plt.grid(True, alpha=0.3)

plt.tight_layout()
plt.savefig("./grafics/learning_curves.png", dpi=300, bbox_inches='tight')
plt.close()




print("\n" + "=" * 50)
print("LEARNING CURVES FOR DIFFERENT SPLITS")

plt.figure(figsize=(15, 10))

all_curves = {0.2: [], 0.3: [], 0.4: []}

for test_s in [0.2, 0.3, 0.4]:
    for split_idx, random_state in enumerate([42, 123, 456, 789, 999]):
        X_temp_train, X_temp_test, y_temp_train, y_temp_test = train_test_split(
            X[selected_features], y_binary,
            test_size=test_s, random_state=random_state, stratify=y_binary
        )

        model = GradientDescentLinearClassifier(
            loss='logistic', learning_rate=0.01, alpha=0.001,
            l1_ratio=0.5, max_iter=500, random_state=42, verbose=False
        )

        model.fit(X_temp_train, y_temp_train)
        train_losses = model.loss_history_

        all_curves[test_s].append(train_losses)

for test_s, curves in all_curves.items():
    min_length = min(len(curve) for curve in curves)
    curves_aligned = [curve[:min_length] for curve in curves]

    mean_curve = np.mean(curves_aligned, axis=0)
    std_curve = np.std(curves_aligned, axis=0)
    n_curves = len(curves)

    ci_lower = mean_curve - 1.96 * std_curve / np.sqrt(n_curves)
    ci_upper = mean_curve + 1.96 * std_curve / np.sqrt(n_curves)

    plt.figure(figsize=(10, 6))

    for i, curve in enumerate(curves_aligned):
        plt.plot(curve, alpha=0.3, color='blue', linewidth=0.5)

    plt.plot(mean_curve, 'r-', linewidth=2, label='Mean training loss')

    plt.fill_between(range(min_length), ci_lower, ci_upper,
                     alpha=0.3, color='red', label='95% Confidence Interval')


    # target val of curve of ridge linear classifier
    ridge_model = RidgeLinearClassifier(alpha=1.0, fit_intercept=True)
    ridge_model.fit(X_train_selected, y_train)

    X_test_processed = np.column_stack([np.ones(X_test_selected.shape[0]), X_test_selected])
    y_test_binary = np.where(y_test == ridge_model.classes_[0], -1, 1)
    ridge_test_loss = ridge_model.compute_loss(y_test_binary * (X_test_processed @ ridge_model.coef_))

    plt.axhline(y=ridge_test_loss, color='purple', linestyle='--', linewidth=2,
                label=f'Ridge Test Loss: {ridge_test_loss:.4f}')


    plt.xlabel('Iteration')
    plt.ylabel('Training Loss')
    plt.title(f'Learning Curves with Confidence Interval\n(Test size: {test_s}, {n_curves} splits)')
    plt.legend()
    plt.grid(True, alpha=0.3)
    plt.savefig(f"./grafics/learning_curve_ci_test_size_{test_s}.png", dpi=300, bbox_inches='tight')
    plt.close()





# Regularization path analysis
print("\n" + "="*50)
print("REGULARIZATION PATH ANALYSIS")

# L1 regularization path (L2 coefficient = 0)
l1_alphas = np.logspace(-3, 2, 20)  # logarithmic space
l1_coef_path = []

X_processed = np.column_stack([np.ones(X_train_selected.shape[0]), X_train_selected])
y_binary_train = np.where(y_train == np.unique(y_train)[0], -1, 1)

for alpha in l1_alphas:
    # For L1, we need to use gradient descent with l1_ratio=1.0
    model = GradientDescentLinearClassifier(
        loss='logistic', learning_rate=0.0001, alpha=alpha,
        l1_ratio=1.0, max_iter=1000, tol=1e10, random_state=42, verbose=False
    )
    model.fit(X_train_selected, y_train)
    l1_coef_path.append(model.coef_[1:].copy())  # exclude intercept

# L2 regularization path (L1 coefficient = 0)
l2_alphas = np.logspace(-3, 2, 20)
l2_coef_path = []

for alpha in l2_alphas:
    model = GradientDescentLinearClassifier(
        loss='logistic', learning_rate=0.0001, alpha=alpha,
        l1_ratio=0.0, max_iter=1000, random_state=42, verbose=False
    )
    model.fit(X_train_selected, y_train)
    l2_coef_path.append(model.coef_[1:].copy())

# Plot regularization paths
plt.figure(figsize=(15, 6))

plt.subplot(1, 2, 1)
for i in range(min(10, len(l1_coef_path[0]))):
    coef_values = [coef[i] for coef in l1_coef_path]
    plt.plot(l1_alphas, coef_values, linewidth=2, label=selected_features[i])

plt.xscale('log')
plt.xlabel('L1 Regularization Coefficient (alpha)')
plt.ylabel('Coefficient Value')
plt.title('L1 Regularization Path\n(L2 coefficient = 0)')
plt.legend(bbox_to_anchor=(1.05, 1), loc='upper left')
plt.grid(True, alpha=0.3)

plt.subplot(1, 2, 2)
for i in range(min(10, len(l2_coef_path[0]))):
    coef_values = [coef[i] for coef in l2_coef_path]
    plt.plot(l2_alphas, coef_values, linewidth=2, label=selected_features[i])

plt.xscale('log')
plt.xlabel('L2 Regularization Coefficient (alpha)')
plt.ylabel('Coefficient Value')
plt.title('L2 Regularization Path\n(L1 coefficient = 0)')
plt.legend(bbox_to_anchor=(1.05, 1), loc='upper left')
plt.grid(True, alpha=0.3)

plt.tight_layout()
plt.savefig("./grafics/regularization_paths.png", dpi=300, bbox_inches='tight')
plt.close()





print(f"Best model: {type(best_model).__name__}")
print(f"Best parameters: {optimizer.best_params_['ridge' if ridge_score > gradient_score else 'gradient']}")
print(f"Test accuracy: {accuracy_score(y_test, y_pred):.4f}")







