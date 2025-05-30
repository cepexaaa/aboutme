#!/bin/bash

INPUT_FILE="top_data1.log"

OUTPUT_RAM="ram_data_first.txt"
OUTPUT_SWAP="swap_data_first.txt"
OUTPUT_MEM_BASH="mem_bash_data_first.txt"
OUTPUT_MEM2_BASH="mem2_bash_data.txt"

> "$OUTPUT_RAM"
> "$OUTPUT_SWAP"
> "$OUTPUT_MEM_BASH"
> "$OUTPUT_MEM2_BASH"

extract_ram_data() {
    local line=$1
#    echo "$line" | awk '{print $3}' >> "$OUTPUT_RAM"  # total
#    echo "$line" | awk '{print $5}' >> "$OUTPUT_RAM"  # free
    echo "$line" | awk '{print $8}' >> "$OUTPUT_RAM"  # used
#    echo "$line" | awk '{print $9}' >> "$OUTPUT_RAM"  # buff/cache
}

extract_swap_data() {
    local line=$1
 #   echo "$line" | awk '{print $3}' >> "$OUTPUT_SWAP"  # total
#    echo "$line" | awk '{print $5}' >> "$OUTPUT_SWAP"  # free
    echo "$line" | awk '{print $7}' >> "$OUTPUT_SWAP"  # used
}

extract_process_data() {
    local line=$1
    local process_name=$(echo "$line" | awk '{print $12}')
    local output_file="$OUTPUT_MEM_BASH"

    if [[ "$process_name" == *"mem.bash"* ]]; then
        output_file="$OUTPUT_MEM_BASH"
    elif [[ "$process_name" == *"mem2.ba"* ]]; then
	output_file="$OUTPUT_MEM2_BASH"
    fi

    echo "$line" | awk '{print $5}' >> "$output_file"  # %CPU
#    echo "$line" | awk '{print $7}' >> "$output_file"  # %MEM
#    echo "$line" | awk '{print $8}' >> "$output_file"  # VIRT
#    echo "$line" | awk '{print $9}' >> "$output_file"  # RES
}

while IFS= read -r line; do
    if echo "$line" | grep -qE "MiB Mem"; then
        extract_ram_data "$line"
    fi

    if echo "$line" | grep -qE "MiB Swap"; then
        extract_swap_data "$line"
    fi

    if echo "$line" | grep -qE "mem.bash|mem2.ba"; then
        extract_process_data "$line"
    fi
done < "$INPUT_FILE"

echo "Данные успешно извлечены в файлы: $OUTPUT_RAM, $OUTPUT_SWAP, $OUTPUT_MEM_BASH, $OUTPUT_MEM2_BASH"
