import requests
import time
import random
import pandas as pd
import numpy as np
from bs4 import BeautifulSoup
from concurrent.futures import ThreadPoolExecutor, as_completed
import os
import extractAtributs as ea

# source venv/bin/activate

def generate_mass_star_identifiers():
    identifiers = []

    # 1. HIP каталог (120000 звезд)
    print("Generating HIP identifiers...")
    hip_ids = [f"HIP {i}" for i in random.sample(range(1, 120000), 6000)]
    identifiers.extend(hip_ids)

    # 2. HD каталог (400000+ звезд)
    print("Generating HD identifiers...")
    hd_ids = [f"HD {i}" for i in random.sample(range(1, 400000), 6000)]
    identifiers.extend(hd_ids)

    # 3. TYC каталог (миллионы звезд)
    print("Generating TYC identifiers...")
    tyc_ids = []
    for _ in range(3000):
        tyc1 = random.randint(1, 9537)
        tyc2 = random.randint(1, 12121)
        tyc3 = random.randint(1, 5)
        tyc_ids.append(f"TYC {tyc1}-{tyc2}-{tyc3}")
    identifiers.extend(tyc_ids)

    # 4. HR каталог (9110 ярких звезд)
    print("Generating HR identifiers...")
    hr_ids = [f"HR {i}" for i in random.sample(range(1, 9110), 2000)]
    identifiers.extend(hr_ids)

    random.shuffle(identifiers)
    print(f"Generated {len(identifiers)} total identifiers")
    return identifiers


# Extract data from html
def fetch_and_parse_star(object_id):
    base_url = "https://simbad.cds.unistra.fr/simbad/sim-id?Ident="
    headers = {
        'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36',
        'Accept-Language': 'en-US,en;q=0.5',
    }

    try:
        encoded_id = requests.utils.quote(object_id)
        response = requests.get(base_url + encoded_id, headers=headers, timeout=15)
        response.raise_for_status()

        star_data = ea.parse_star_attributes(response.text, object_id)

        time.sleep(random.uniform(1, 3))

        return star_data

    except Exception as e:
        print(f"Error fetching {object_id}: {e}")
        return None


def mass_star_data_collection():
    all_identifiers = generate_mass_star_identifiers()

    all_star_data = []
    processed_count = 0
    success_count = 0

    print("Starting mass data collection (Using thread pool)...")

    with ThreadPoolExecutor(max_workers=3) as executor:
        # run tasks
        future_to_id = {executor.submit(fetch_and_parse_star, obj_id): obj_id for obj_id in all_identifiers[:15000]}

        for future in as_completed(future_to_id):
            obj_id = future_to_id[future]
            processed_count += 1

            try:
                result = future.result()
                if result is not None:
                    all_star_data.append(result)
                    success_count += 1

                    # save progress
                    if success_count % 100 == 0:
                        print(f"Progress: {processed_count}/{15000} processed, {success_count} successful")
                        temp_df = pd.DataFrame(all_star_data)
                        temp_df.to_csv('temp_star_data.tsv', sep='\t', index=False)

            except Exception as e:
                print(f"Error processing {obj_id}: {e}")

            if success_count >= 10000: # 10000:
                print(f"Reached target of {success_count} stars! Stopping...")
                for f in future_to_id:
                    f.cancel()
                break

    print(f"Collection complete: {success_count} stars collected")
    return all_star_data



print("=== MASS STAR DATA COLLECTION ===")
print("This will collect ~10000 stars (may take several hours)")
print("Progress will be saved periodically to temp_star_data.tsv")

star_data = mass_star_data_collection()

# save raw data
if star_data:
    df_raw = pd.DataFrame(star_data)
    df_raw.to_csv('data.tsv', sep='\t', index=False)
    print(f"Raw data saved: {len(df_raw)} stars")

    if os.path.exists('temp_star_data.tsv'):
        os.remove('temp_star_data.tsv')
else:
    print("No data collected!")

