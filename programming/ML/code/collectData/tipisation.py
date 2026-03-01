import json

import pandas as pd


def save_arff(df, filename):
    with open(filename, 'w', encoding='utf-8') as f:
        f.write('@RELATION stars\n\n')

        for column in df.columns:
            if column == 'object_id':
                f.write(f"@ATTRIBUTE {column} STRING\n")
            elif pd.api.types.is_categorical_dtype(df[column]):
                categories = "{" + ",".join([f"'{cat}'" for cat in df[column].cat.categories]) + "}"
                f.write(f"@ATTRIBUTE {column} {categories}\n")
            else:
                f.write(f"@ATTRIBUTE {column} NUMERIC\n")

        f.write('\n@DATA\n')
        for _, row in df.iterrows():
            row_data = []
            for value in row:
                if pd.isna(value):
                    row_data.append('?')
                elif isinstance(value, str):
                    row_data.append(f"'{value}'")
                else:
                    row_data.append(str(value))
            f.write(','.join(row_data) + '\n')


def typify_and_save_star_data():
    df = pd.read_csv('data.tsv', sep='\t')

    df_typed = df.copy()

    numerical_columns = [
        'v_magnitude', 'parallax', 'pm_ra', 'pm_dec', 'radial_velocity',
        'ra_degrees', 'dec_degrees', 'b_magnitude', 'r_magnitude', 'g_magnitude',
        'j_magnitude', 'h_magnitude', 'k_magnitude', 'i_magnitude', 'distance_pc'
    ]

    for col in numerical_columns:
        df_typed[col] = pd.to_numeric(df_typed[col], errors='coerce')

    categorical_columns = ['spectral_class', 'luminosity_class', 'motion_type']

    for col in categorical_columns:
        df_typed[col] = df_typed[col].astype('category')

    save_to_star_json_format(df_typed, 'data.json')

    df_arff = df_typed.copy()
    for col in categorical_columns:
        df_arff[col] = df_arff[col].astype(str)

    save_arff(df_arff, 'data.arff')

    print("Stage 3 complete: data.json and data.arff created")
    print(f"DataFrame shape: {df_typed.shape}")
    print(f"Columns: {list(df_typed.columns)}")

    return df_typed


def save_to_star_json_format(df, filename):
    feature_types = {
        'object_id': {'type': 'text'},
        'spectral_class': {
            'type': 'category',
            'values': ['O', 'B', 'A', 'F', 'G', 'K', 'M', 'Unknown']
        },
        'luminosity_class': {
            'type': 'category',
            'values': ['I', 'II', 'III', 'IV', 'V', 'VI', 'Unknown']
        },
        'is_variable': {
            'type': 'category',
            'values': ['Constant', 'Variable']
        },
        'motion_type': {
            'type': 'category',
            'values': ['Normal', 'High_Proper_Motion']
        },
        'v_magnitude': {'type': 'numeric'},
        'parallax': {'type': 'numeric'},
        'pm_ra': {'type': 'numeric'},
        'pm_dec': {'type': 'numeric'},
        'radial_velocity': {'type': 'numeric'},
        'ra_degrees': {'type': 'numeric'},
        'dec_degrees': {'type': 'numeric'},
        'b_magnitude': {'type': 'numeric'},
        'r_magnitude': {'type': 'numeric'},
        'g_magnitude': {'type': 'numeric'},
        'j_magnitude': {'type': 'numeric'},
        'h_magnitude': {'type': 'numeric'},
        'k_magnitude': {'type': 'numeric'},
        'i_magnitude': {'type': 'numeric'},
        'distance_pc': {'type': 'numeric'}
    }

    # Создаем header
    header = []
    for column in df.columns:
        if column in feature_types:
            feature_info = feature_types[column].copy()
            feature_info['feature_name'] = column

            if feature_info['type'] == 'category':
                unique_vals = df[column].dropna().unique()
                if len(unique_vals) > 0:
                    actual_values = sorted([str(x) for x in unique_vals])
                    feature_info['values'] = actual_values

            header.append(feature_info)
        else:
            # Автоматическое определение для неизвестных колонок
            dtype = df[column].dtype
            if pd.api.types.is_numeric_dtype(dtype):
                header.append({
                    "feature_name": column,
                    "type": "numeric"
                })
            else:
                unique_count = df[column].nunique()
                if unique_count <= 50:
                    unique_vals = df[column].dropna().unique()
                    header.append({
                        "feature_name": column,
                        "type": "category",
                        "values": sorted([str(x) for x in unique_vals if pd.notna(x)])
                    })
                else:
                    header.append({
                        "feature_name": column,
                        "type": "text"
                    })

    # Преобразуем данные в список словарей (заменяем NaN на None)
    data = []
    for _, row in df.iterrows():
        row_dict = {}
        for col in df.columns:
            value = row[col]
            if pd.isna(value):
                row_dict[col] = None
            elif hasattr(value, 'item'):  # Для numpy типов
                row_dict[col] = value.item() if hasattr(value, 'item') else value
            else:
                row_dict[col] = value
        data.append(row_dict)

    result = {
        "header": header,
        "data": data
    }

    with open(filename, 'w', encoding='utf-8') as f:
        json.dump(result, f, indent=2, ensure_ascii=False)

    print(f"Данные сохранены в {filename} в формате header/data")
    return result