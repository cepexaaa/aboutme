#!/bin/sh

echo "Testing multithreading performance..."

# Функция для измерения времени выполнения
measure_time() {
    threads=$1
    output_file="test_output_${threads}_threads.png"
    echo "Running with $threads threads..."

    # Измерение времени выполнения
    START_TIME=$(date +%s)
    "$EXECUTABLE_PATH" -w 1920 -h 1080 -t "$threads" -o "$output_file"
    EXIT_CODE=$?
    END_TIME=$(date +%s)

    # Вычисление длительности
    DURATION=$((END_TIME - START_TIME))
    if [ $EXIT_CODE -eq 0 ]; then
        echo "✓ Completed with $threads threads in ${DURATION} seconds"
        echo "$threads,$DURATION" >> performance_results.csv
        return 0
    else
        echo "✗ Failed with $threads threads (exit code: $EXIT_CODE)"
        return 1
    fi
}

# Проверка наличия EXECUTABLE_PATH
if [ -z "$1" ]; then
    echo "✗ Executable file path not provided."
    echo "Usage: $0 <path_to_executable_file>"
    exit 1
fi

EXECUTABLE_PATH="$1"

# Проверка существования исполняемого файла
if [ ! -f "$EXECUTABLE_PATH" ]; then
    echo "✗ Executable file '$EXECUTABLE_PATH' does not exist."
    exit 1
fi

# Инициализация файла с результатами
echo "threads,duration_seconds" > performance_results.csv

# Тестирование с разным количеством потоков
for threads in 1 2 4; do
    if ! measure_time "$threads"; then
        echo "Performance test failed for $threads threads"
        exit 1
    fi
done

echo ""
echo "Performance test results:"
echo "------------------------"
cat performance_results.csv
echo ""
echo "Multithreading performance test completed!"
