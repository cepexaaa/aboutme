#!/bin/bash

process_count=$(( $(ps -u $USER | wc -l) - 1))

if [ -f "fileToTask1.txt" ]; then
	rm fileToTask.txt
fi

echo "$process_count" >> fileToTask.txt

ps -u $USER -o pid,comm --no-headers | awk '{print $1 ":" $2}' >> fileToTask.txt

exit 0
