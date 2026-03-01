import os

import pandas as pd
from sklearn.feature_extraction.text import CountVectorizer
from sklearn.model_selection import train_test_split


def load_and_prepare_data(file_path='unsup-cepexaaa/data/castle-or-lock.tsv', test_size=0.2, random_state=42):
    if not os.path.exists(file_path):
        full_path = os.path.abspath(file_path)
        raise FileNotFoundError(f"File {file_path} (full path: {full_path}) not found!")

    df = pd.read_csv(file_path, sep='\t', encoding='utf-8')
    label_column = df.columns[0]
    text_column = df.columns[1]

    label_map = {'castle': 0, 'lock': 1}
    df['label_encoded'] = df[label_column].map(label_map)
    y = df['label_encoded'].astype(int).values

    print(f"\nClass distribution after processing:")
    print(f"  Castle (0): {(y == 0).sum()} records")
    print(f"  Lock (1): {(y == 1).sum()} records\n")

    vectorizer = CountVectorizer(
        max_features=10000,
        min_df=2,
        max_df=0.95,
        ngram_range=(1, 2),
        stop_words=['и', 'в', 'на', 'с', 'по', 'для', 'от', 'или']
    )

    X = vectorizer.fit_transform(df[text_column])

    print(f"Feature matrix shape: {X.shape}")
    print(f"Vocabulary size: {len(vectorizer.get_feature_names_out())}\n")


    X_train, X_test, y_train, y_test, idx_train, idx_test = train_test_split(
        X, y, df.index,
        test_size=test_size,
        random_state=random_state,
        stratify=y
    )

    print(f"Training set: {X_train.shape[0]} records")
    print(f"Test set: {X_test.shape[0]} records\n")
    result = {
        'X': X,
        'y': y,
        'X_train': X_train,
        'X_test': X_test,
        'y_train': y_train,
        'y_test': y_test,
        'vectorizer': vectorizer,
        'df': df,
        'df_train': df.loc[idx_train],
        'df_test': df.loc[idx_test],
        'feature_names': vectorizer.get_feature_names_out(),
        'label_mapping': {0: 'castle', 1: 'lock'},
        'original_file_path': file_path
    }
    return result
