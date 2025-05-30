#!/bin/bash

if [ ! -f "fileToTask.txt" ]; then
    echo "File fileToTask.txt didn't find."
    exit 1
fi

temp_file=$(mktemp)

current_ppid=0
sum_art=0
count=0

while IFS= read -r line; do
    ppid=$(echo "$line" | awk -F'[= ]' '{print $5}')
    art=$(echo "$line" | awk -F'[= ]' '{print $8}')

    if [ "$ppid" != "$current_ppid" ]; then
        avg_art=$(echo "scale=6; $sum_art / $count" | bc | xargs printf "%.6f")
        echo "Average_Running_Children_of_ParentID=$current_ppid is $avg_art" >> "$temp_file"
        sum_art=0
        count=0
    fi

    current_ppid=$ppid
    sum_art=$(echo "scale=6; $sum_art + $art" | bc | xargs printf "%.6f")
    count=$((count + 1))

    echo "$line" >> "$temp_file"
done < fileToTask.txt

if [ "$count" -gt 0 ]; then
    avg_art=$(echo "scale=6; $sum_art / $count" | bc | xargs printf "%.6f")
    echo "Average_Running_Children_of_ParentID=$current_ppid is $avg_art" >> "$temp_file"
fi

mv "$temp_file" fileToTask.txt

exit 0



