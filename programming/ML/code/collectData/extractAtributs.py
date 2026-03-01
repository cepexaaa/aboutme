import re
import numpy as np
from bs4 import BeautifulSoup


def parse_star_attributes(html_content, object_id):
    try:
        soup = BeautifulSoup(html_content, 'html.parser')

        if "Star" not in html_content and "star" not in html_content.lower():
            return None

        data_row = {'object_id': object_id}

        # 1. Спектральный класс
        spectral_type = extract_value_by_label(soup, "Spectral type:")
        data_row['spectral_class'] = clean_spectral_type(spectral_type)

        # 2. Класс светимости - извлекаем из спектрального типа
        luminosity_class = extract_luminosity_class(spectral_type)
        data_row['luminosity_class'] = luminosity_class

        # 4. Тип движения
        motion_type = "Normal"
        if "High Proper Motion" in html_content:
            motion_type = "High_Proper_Motion"
        data_row['motion_type'] = motion_type

        # 5. Видимая звездная величина V
        v_magnitude = extract_flux_value(soup, "V")
        data_row['v_magnitude'] = parse_float(v_magnitude)

        # 6. Параллакс
        parallax = extract_value_by_label(soup, "Parallaxes")
        data_row['parallax'] = parse_float(clean_parallax_value(parallax))

        # 7-8. Собственное движение
        pm_ra, pm_dec = extract_proper_motion(soup)
        data_row['pm_ra'] = parse_float(pm_ra)
        data_row['pm_dec'] = parse_float(pm_dec)

        # 9. Лучевая скорость
        radial_velocity = extract_radial_velocity(soup)
        data_row['radial_velocity'] = parse_float(radial_velocity)

        # 10-11. Координаты RA и DEC
        ra_deg, dec_deg = extract_coordinates(soup)
        data_row['ra_degrees'] = ra_deg
        data_row['dec_degrees'] = dec_deg

        # 12-17. Другие фотометрические величины
        photometry_bands = ['B', 'R', 'G', 'J', 'H', 'K', 'g', 'r', 'i']
        for band in photometry_bands:
            magnitude = extract_flux_value(soup, band)
            data_row[f'{band.lower()}_magnitude'] = parse_float(magnitude)

        # 18. Расстояние (вычисляем из параллакса)
        if data_row['parallax'] and data_row['parallax'] > 0:
            data_row['distance_pc'] = round(1000 / data_row['parallax'], 2)
        else:
            data_row['distance_pc'] = np.nan

        valid_fields = sum(1 for key, value in data_row.items()
                           if value not in [np.nan, "N/A", "Unknown", ""] and key != 'object_id')

        if valid_fields >= 13:
            return data_row
        else:
            print(f"✗ {object_id} - insufficient data ({valid_fields} fields)")
            return None

    except Exception as e:
        print(f"Error parsing {object_id}: {e}")
        return None


def extract_value_by_label(soup, label):
    try:
        label_elements = soup.find_all(string=lambda text: text and label.strip() in text)
        for label_element in label_elements:
            label_td = label_element.find_parent('td')
            if not label_td:
                continue
            if label_td.find('td'):
                continue
            parent_tr = label_td.find_parent('tr')
            if not parent_tr:
                continue
            all_tds = parent_tr.find_all('td')
            if len(all_tds) < 2:
                continue
            try:
                label_index = all_tds.index(label_td)
                if label_index + 1 < len(all_tds):
                    value_td = all_tds[label_index + 1]

                    value_text = value_td.get_text(strip=True)

                    if value_text and value_text != "N/A":
                        # Очищаем от bibcode и качества и Убираем буквы качества (A, B, C, D, E) в конце
                        value_text = re.sub(r'[ABCDE]\s*$', '', value_text).strip()
                        # Убираем bibcode (формат года + Journal)
                        value_text = re.sub(r'\d{4}[A-Za-z]+\..*$', '', value_text).strip()

                        if value_text:
                            return value_text
            except ValueError:
                continue

    except Exception as e:
        print(f"Error extracting {label}: {e}")

    return "N/A"

# Extract photometry to certain filter
def extract_flux_value(soup, band):
    try:
        band_elements = soup.find_all('b')

        for element in band_elements:
            text = element.get_text(strip=True)

            if text.startswith(f"{band} "):
                parts = text.split()
                if len(parts) >= 2:
                    value = parts[1].strip('[]')
                    if re.match(r'^[-+]?\d*\.?\d+$', value):
                        return value

    except Exception as e:
        print(f"Error extracting flux for band {band}: {e}")

    return "N/A"


# Self moving
def extract_proper_motion(soup):
    try:
        pm_text = extract_value_by_label(soup, "Proper motions")
        if pm_text != "N/A":
            matches = re.findall(r'[-+]?\d*\.?\d+', pm_text)
            if len(matches) >= 2:
                return matches[0], matches[1]
    except:
        pass
    return "N/A", "N/A"


# Это проекция скорости звезды вдоль луча зрения наблюдателя, то есть скорость приближения или удаления звезды относительно Земли.
def extract_radial_velocity(soup):
    try:
        rv_text = extract_value_by_label(soup, "Radial velocity")
        if rv_text != "N/A":
            # V(km/s) number
            match = re.search(r'V\(km/s\)\s*([-+]?\d*\.?\d+)', rv_text)
            if match:
                return match.group(1)
            matches = re.findall(r'[-+]?\d*\.?\d+', rv_text)
            if matches:
                return matches[0]
    except:
        pass
    return "N/A"


def extract_coordinates(soup):
    try:
        coords_text = extract_value_by_label(soup, "ICRS")
        if coords_text != "N/A":
            match = re.search(r'(\d{2})\s+(\d{2})\s+(\d{2}\.\d+)\s+([-+]?\d{2})\s+(\d{2})\s+(\d{2}\.\d+)', coords_text)
            if match:
                ra_h, ra_m, ra_s, dec_d, dec_m, dec_s = match.groups()

                ra_deg = 15 * (float(ra_h) + float(ra_m) / 60 + float(ra_s) / 3600)

                dec_sign = -1 if dec_d.startswith('-') else 1
                dec_deg = dec_sign * (abs(float(dec_d)) + float(dec_m) / 60 + float(dec_s) / 3600)

                return round(ra_deg, 6), round(dec_deg, 6)
    except:
        pass
    return np.nan, np.nan


def extract_luminosity_class(spectral_type):
    if not spectral_type or spectral_type == "N/A":
        return "V"

    luminosity_class = "V"  # dwarf
    roman_pattern = r'[IV]+'
    match = re.search(roman_pattern, spectral_type)
    if match:
        luminosity_class = match.group()

    return luminosity_class


def clean_spectral_type(sp_type):
    if not sp_type or sp_type == "N/A":
        return "Unknown"

    # first letter of spectral class
    sp_class = sp_type[0].upper() if sp_type and sp_type[0].isalpha() else "Unknown"
    return sp_class if sp_class in "OBAFGKM" else "Unknown"


def clean_parallax_value(parallax_text):
    if not parallax_text or parallax_text == "N/A":
        return "N/A"
    try:
        clean_text = re.sub(r'\[[^\]]*\]', '', parallax_text)
        clean_text = re.sub(r'[ABCDE]\s*$', '', clean_text).strip()
        clean_text = clean_text.strip()
        match = re.search(r'[-+]?\d*\.?\d+', clean_text)
        if match:
            return match.group()
        else:
            return "N/A"
    except Exception as e:
        print(f"Error cleaning parallax value '{parallax_text}': {e}")
        return "N/A"


def parse_float(value_str):
    if not value_str or value_str == "N/A":
        return np.nan
    try:
        clean_str = re.sub(r'[^\d\.\-+]', '', value_str)
        return float(clean_str) if clean_str else np.nan
    except:
        return np.nan

