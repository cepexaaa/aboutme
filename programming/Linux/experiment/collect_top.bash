#!/bin/bash

OUTPUT_FILE="top_data2.log"
REPORT_LOG="report.log"

> "$OUTPUT_FILE"

collect_top_data() {
    echo "--- Collecting top data at $(date) ---" >> "$OUTPUT_FILE"
    top -b -n 1 | grep -E "Mem|Swap|^Tasks:|^ *PID|%CPU|%MEM|COMMAND" >> "$OUTPUT_FILE"

    echo "--- Top 5 processes ---" >> "$OUTPUT_FILE"
    top -b -n 1 | head -n 12 | tail -n 5 >> "$OUTPUT_FILE"
    echo "---" >> "$OUTPUT_FILE"
}

collect_syslog() {
    echo "--- Collecting syslog data at $(date) ---" >> "$OUTPUT_FILE"
    dmesg | grep "Out of memory" >> "$OUTPUT_FILE"
    echo "---" >> "$OUTPUT_FILE"
}

collect_report_log() {
    echo "--- Collecting last line from report.log at $(date) ---" >> "$OUTPUT_FILE"
    if [ -f "$REPORT_LOG" ]; then
        tail -n 1 "$REPORT_LOG" >> "$OUTPUT_FILE"
    else
        echo "File $REPORT_LOG not found." >> "$OUTPUT_FILE"
    fi
    echo "---" >> "$OUTPUT_FILE"
}

echo "Results is writing in $OUTPUT_FILE..."
while true; do
    collect_top_data
    echo "" >> "$OUTPUT_FILE"
    sleep 1
done
