import requests
import pandas as pd
from datetime import datetime, timedelta
import time
import os

# собираю данные по курсу рубль-доллар за 5 лет
def download(start_date='2020-01-01', end_date='2025-01-31'):
    start = datetime.strptime(start_date, '%Y-%m-%d')
    end = datetime.strptime(end_date, '%Y-%m-%d')

    all_data = []

    current_date = start
    while current_date <= end:
        date_str = current_date.strftime('%d/%m/%Y')
        url = f"https://www.cbr.ru/scripts/XML_daily.asp?date_req={date_str}"

        try:
            response = requests.get(url, timeout=10)
            response.encoding = 'utf-8'

            if response.status_code == 200:
                if 'USD' in response.text:
                    lines = response.text.split('<Valute ID=')
                    for i, line in enumerate(lines):
                        if 'USD' in line and 'NumCode' in line:
                            for j in range(i, min(i + 10, len(lines))):
                                if 'Value' in lines[j]:
                                    value_line = lines[j]
                                    value = value_line.split("<Value>")[1].split('<')[0]
                                    value = float(value.replace(',', '.'))
                                    all_data.append({
                                        'date': current_date.strftime('%Y-%m-%d'),
                                        'usd_rub': value
                                    })
                                    print(f"OK: {current_date.strftime('%Y-%m-%d')}: {value} руб.")
                                    break
                            break

            time.sleep(0.5)

        except Exception as e:
            print(f"Fail: Error for data {current_date.strftime('%Y-%m-%d')}: {e}")
        current_date += timedelta(days=1)
    df = pd.DataFrame(all_data)
    os.makedirs("./ds-cepexaaa/code/timeSequences/data", exist_ok=True)
    df.to_csv("./ds-cepexaaa/code/timeSequences/data/rub-usd.csv", index=False, encoding='utf-8')
    return df


def aggregation(data):
    df = data.copy()
    df['date'] = pd.to_datetime(df['date'])

    df_agg = df.resample('W', on='date').agg({'usd_rub': 'mean'}).reset_index()
    # df_agg = df.resample('ME', on='date').agg({'usd_rub': 'mean'}).reset_index()

    df_agg.to_csv("./ds-cepexaaa/code/timeSequences/data/aggregated_rub-usd.csv", index=False, encoding='utf-8')
    print(f"aggregation data was saved. Items: {len(df_agg)}")

    return df_agg