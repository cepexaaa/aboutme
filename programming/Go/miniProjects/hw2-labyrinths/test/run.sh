#!/bin/bash

set -euo pipefail
#set -x

# Цвета ANSI для вывода
RED='\033[1;31m'
GREEN='\033[1;32m'
YELLOW='\033[1;33m'
RESET='\033[0m'

# Функция для вывода подробного сравнения файлов с описанием различий
diff_with_details() {
  local expected_file=$1
  local actual_file=$2
  local prefix=$3

  if diff -u "$expected_file" "$actual_file" > /tmp/diff_output.txt ; then
    return 0
  else
    printf "${RED}${prefix}Ожидалось:${RESET}\n"
    sed -n '1,10p' "$expected_file" | sed 's/^/  /'
    printf "${RED}${prefix}Получено:${RESET}\n"
    sed -n '1,10p' "$actual_file" | sed 's/^/  /'
    printf "${RED}${prefix}Подробности различий (diff -u):${RESET}\n"
    sed 's/^/  /' /tmp/diff_output.txt
    return 1
  fi
}

# Сравнение файлов из expected.files.txt
compare_expected_files() {
  local test_dir=$1
  local expected_files_txt="$test_dir/expected.files.txt"
  local failed=false

  [[ -f "$expected_files_txt" ]] || return 0

  local current_file=""
  local current_content=""
  # Временный файл для ожидаемого содержимого
  local temp_expected=$(mktemp)
  # Временный файл для фактического содержимого
  local temp_actual=$(mktemp)

  while IFS= read -r line || [[ -n "$line" ]]; do
    if [[ -z "$current_file" ]]; then
      current_file="$line"
      file_path="$current_file"
      current_content=""
    elif [[ "$line" == "===EOF===" ]]; then
      # Записать expected содержимое во временный файл
      echo "$current_content" > "$temp_expected"
      if [[ ! -f "$file_path" ]]; then
        printf "${RED}Ожидаемый файл '$file_path' не найден${RESET}\n"
        failed=true
      else
        cp "$file_path" "$temp_actual"
        if ! diff_with_details "$temp_expected" "$temp_actual" "Файл '$current_file': "; then
          failed=true
        fi
      fi
      current_file=""
      current_content=""
    else
      if [[ -z "$current_content" ]]; then
        current_content="$line"
      else
        current_content="$current_content
$line"
      fi
    fi
  done < "$expected_files_txt"

  rm -f "$temp_expected" "$temp_actual"

  $failed && return 1 || return 0
}

# Запуск и проверка одного теста с run_before.txt и run_after.txt
run_test_case() {
  local test_dir=$1
  local run_cmd_file="$test_dir/run.txt"
  local run_before_file="$test_dir/run_before.txt"
  local run_after_file="$test_dir/run_after.txt"
  local input_file="$test_dir/input.txt"
  local expected_output_file="$test_dir/expected.output.txt"

  if [[ ! -f "$run_cmd_file" ]]; then
    printf "${YELLOW}Пропущен run.txt в $test_dir${RESET}\n"
    return 1
  fi

  # Запуск run_before.txt если есть
  if [[ -f "$run_before_file" ]]; then
    sh -c "$(cat "$run_before_file")"
    local before_status=$?
    if [[ $before_status -ne 0 ]]; then
      printf "${RED}Ошибка при выполнении run_before.txt в $test_dir${RESET}\n"
      return 1
    fi
  fi

  local run_cmd
  run_cmd=$(cat "$run_cmd_file")

  local output_file
  output_file=$(mktemp)

  if [[ -f "$input_file" ]]; then
    sh -c "$run_cmd" < "$input_file" > "$output_file" 2>&1
  else
    sh -c "$run_cmd" > "$output_file" 2>&1
  fi

  local passed=true

  # Проверка вывода
  if [[ -f "$expected_output_file" ]]; then
    if ! diff_with_details "$expected_output_file" "$output_file" "Вывод: "; then
      passed=false
    fi
  fi

  # Проверка файлов из expected.files.txt
  if ! compare_expected_files "$test_dir"; then
    passed=false
  fi

  rm -f "$output_file"

  # Запуск run_after.txt если есть, только если тест прошел
  if $passed && [[ -f "$run_after_file" ]]; then
    sh -c "$(cat "$run_after_file")"
    local after_status=$?
    if [[ $after_status -ne 0 ]]; then
      printf "${RED}Ошибка при выполнении run_after.txt в $test_dir${RESET}\n"
      return 1
    fi
  fi

  if $passed; then
    return 0
  else
    return 1
  fi
}

# Запуск всех тестов из указанного каталога
run_all_tests() {
  local base_dir=$1
  local passed=0
  local failed=0

  if [[ ! -d "$base_dir" ]]; then
    printf "${RED}Каталог $base_dir не найден${RESET}\n"
    exit 1
  fi

  local output=""
  for test_dir in "$base_dir"/*/; do
    [ -d "$test_dir" ] || continue
    set +e
    output=$(run_test_case "$test_dir" 2>&1)
    ret_code=$?
    set -e
    if [ $ret_code -eq 0 ]; then
      printf "${GREEN}Тест $test_dir: ПРОЙДЕН${RESET}\n"
      passed=$(expr "$passed" + 1)
    else
      printf "${RED}Тест $test_dir: НЕ ПРОЙДЕН${RESET}\n"
      printf "$output\n"
      failed=$(expr "$failed" + 1)
    fi
  done

  printf "\n"
  printf "ИТОГИ ПРОВЕРКИ:\n"
  printf "${GREEN}Успешно: $passed${RESET}, ${RED}Неудачи: $failed${RESET}\n"

  if [ "$failed" -gt 0 ]; then
    exit 1
  else
    exit 0
  fi
}

# Проверка параметров и запуск
if [ "$#" -ne 1 ]; then
  printf "${YELLOW}Использование: %s <каталог_тестов>${RESET}\n" "$0"
  exit 1
fi

run_all_tests "$1"
