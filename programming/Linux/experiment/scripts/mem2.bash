#!/bin/bash

> logs/report2.log

array=()
step=0


while true; do
    array+=(1 2 3 4 5 6 7 8 9 10)
    step=$((step + 1))

    if [ $((step % 100000)) -eq 0 ]; then
        echo "Step: $step, Array size: ${#array[@]}" >> logs/report2.log
    fi
done
