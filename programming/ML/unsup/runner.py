from validate import compare_feature_selection_methods_with_cv
from sklearn.linear_model import LogisticRegression
from feature_selection import FilterFeatureSelector, WrapperFeatureSelector, EmbeddedFeatureSelector
from sklearn.ensemble import RandomForestClassifier
from sklearn.feature_selection import VarianceThreshold, RFE, SelectFromModel

def run_my_methods(data):
    filter_selector = FilterFeatureSelector(k=500)
    wrapper_selector = WrapperFeatureSelector(
        estimator=LogisticRegression(random_state=42),
        n_features_to_select=300
    )
    embedded_selector = EmbeddedFeatureSelector(C=0.5, threshold=1e-5, max_iter=5000)
    selectors_dict = {
        'Filter': filter_selector,
        'Wrapper': wrapper_selector,
        'Embedded': embedded_selector
    }
    results = compare_feature_selection_methods_with_cv(
        original_data=data,
        selectors=selectors_dict,
        base_model=LogisticRegression(random_state=42),
        n_splits=5
    )

    return results


def run_lib_methods(data):
    selectors_dict = {
        'VarianceThreshold': VarianceThreshold(
            threshold=0.001
        ),
        'RFE_RandomForest': RFE(
            estimator=RandomForestClassifier(
                n_estimators=50,
                random_state=42
            ),
            n_features_to_select=150,
            step=30
        ),
        'RandomForest': SelectFromModel(
            RandomForestClassifier(
                n_estimators=100,
                random_state=42,
                n_jobs=-1
            ),
            threshold='1.25*median'
        ),
    }
    results = compare_feature_selection_methods_with_cv(
        original_data=data,
        selectors=selectors_dict,
        base_model=LogisticRegression(random_state=42),
        n_splits=5
    )

    return results
