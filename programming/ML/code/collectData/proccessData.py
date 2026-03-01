from sklearn.impute import SimpleImputer
from sklearn.preprocessing import StandardScaler, LabelEncoder, OneHotEncoder
from sklearn.compose import ColumnTransformer
from sklearn.pipeline import Pipeline
import pandas as pd
import numpy as np

import tipisation


def preprocess_star_data_manual(df):

    df_processed = df.copy()

    target_column = 'spectral_class'

    print("Original columns:", list(df_processed.columns))
    print("Missing values before imputation:")
    print(df_processed.isnull().sum())

    numerical_columns = [
        'v_magnitude', 'parallax', 'pm_ra', 'pm_dec', 'radial_velocity',
        'ra_degrees', 'dec_degrees', 'b_magnitude', 'r_magnitude', 'g_magnitude',
        'j_magnitude', 'h_magnitude', 'k_magnitude', 'i_magnitude', 'distance_pc'
    ]

    numerical_columns = [col for col in numerical_columns if col in df_processed.columns]

    for col in numerical_columns:
        if df_processed[col].isnull().any():
            median_val = df_processed[col].median()
            df_processed[col] = df_processed[col].fillna(median_val)
            print(f"Filled missing values in {col} with median: {median_val}")

    categorical_columns = ['luminosity_class', 'motion_type']
    categorical_columns = [col for col in categorical_columns if col in df_processed.columns and col != target_column]

    for col in categorical_columns:
        if df_processed[col].isnull().any():
            df_processed[col] = df_processed[col].fillna('Unknown')
            print(f"Filled missing values in {col} with 'Unknown'")

    # target sign - most frequency value
    if target_column in df_processed.columns and df_processed[target_column].isnull().any():
        if not df_processed[target_column].mode().empty:
            mode_val = df_processed[target_column].mode()[0]
            df_processed[target_column] = df_processed[target_column].fillna(mode_val)
            print(f"Filled missing values in {target_column} with mode: {mode_val}")
        else:
            df_processed[target_column] = df_processed[target_column].fillna('Unknown')
            print(f"Filled missing values in {target_column} with 'Unknown'")

    categorical_for_encoding = [col for col in categorical_columns if col in df_processed.columns]

    print("Applying one-hot encoding to:", categorical_for_encoding)

    # one-hot encoding
    for col in categorical_for_encoding:
        dummies = pd.get_dummies(df_processed[col], prefix=col)
        dummies = dummies.astype(int)
        print(f"Created {len(dummies.columns)} dummy variables for {col}: {list(dummies.columns)}")

        df_processed = pd.concat([df_processed, dummies], axis=1)
        # delete natural column
        df_processed = df_processed.drop(columns=[col])

    # 4. normalisation
    scaler = StandardScaler()
    numerical_cols_to_scale = [col for col in numerical_columns if col in df_processed.columns]

    print("Scaling numerical columns:", numerical_cols_to_scale)

    if numerical_cols_to_scale:
        df_processed[numerical_cols_to_scale] = scaler.fit_transform(df_processed[numerical_cols_to_scale])
        print("Applied StandardScaler to numerical columns")

    # if 'object_id' in df_processed.columns:
    #     df_processed = df_processed.drop(columns=['object_id'])
    #     print("Removed object_id column")

    df_processed.to_csv('data.csv', index=False)

    print("Stage 4 complete: Preprocessed data saved to data.csv")
    print(f"Final dataset shape: {df_processed.shape}")
    print(f"Missing values: {df_processed.isnull().sum().sum()}")
    print(f"Columns after preprocessing: {list(df_processed.columns)}")

    return df_processed


print("Loading and preprocessing data...")
df_typed = tipisation.typify_and_save_star_data()
final_dataset = preprocess_star_data_manual(df_typed)

print("\n=== FINAL DATASET INFO ===")
print(final_dataset.info())
print(f"\nTarget distribution:")
print(final_dataset['spectral_class'].value_counts())
print(f"\nFirst few rows:")
print(final_dataset.head())










