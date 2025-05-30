#!/bin/bash

max_memory=0
max_pid=0
max_cmd=""
# add information about user and
max_user=""

for pid in /proc/[0-9]*; do
	if [ -d "$pid" ]; then
		pid_value=$(basename "$pid")
		rss=$(awk '/^VmRSS:/ {print $2}' "$pid/status")
		if [ -n "$rss" ] && [ "$rss" -gt "$max_memory" ]; then
			max_memory=$rss
			max_pid=$pid_value
			max_cmd=$(cat "$pid/cmdline" | tr '\0' ' ' | sed 's/ *$//')
			max_user=$(awk '/^Uid:/ {print $2}' "$pid/status")
            		max_user=$(getent passwd "$max_user" | cut -d: -f1)
		fi
	fi
done

# getent - get data from DB
#passwd - it is DataBase
# we take only $max_user

echo "PID: $max_pid"
echo "RSS: $max_memory KB"
#echo "Command: $max_cmd"
echo "User: $max_user"

echo -e "\nCompare with command top:"
#top -b -n 1 | grep "^ *[0-9]" | sort -k6,6n | tail -n 1
top -b -n 1 | grep "^ *[0-9]" | sort -k6,6n | tail -n 1 | awk '{print "PID: "$1 "\nRSS: " $6 " KB"}'

# find user with the biggest count of process starts
most_processes_user=$(ps -eo user= | sort | uniq -c | sort -nr | head -n 1 | awk '{print $2}')

echo -e "\nUser with the most processes: $most_processes_user"
 
if [ "$max_user" = "$most_processes_user" ]; then
    echo "This is the same user"
else
    echo "These are two different users"
fi

exit 0
