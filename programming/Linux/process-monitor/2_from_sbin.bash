#!/bin/bash

if [ -f "fileToTask.txt" ]; then 
	rm fileToTask.txt
fi

ps -eo pid,cmd | grep '^ *[0-9]\+ /sbin/' | awk '{print $1}' >> fileToTask.txt

exit 0
