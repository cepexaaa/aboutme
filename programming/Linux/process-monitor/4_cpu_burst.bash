#!/bin/bash

if [ -f "fileToTask.txt" ]; then
	rm fileToTask.txt
fi

for pid in /proc/[0-9]*; do
	if [ -d "$pid" ]; then
		pid_value=$(basename "$pid")
		ppid=$(awk '/^PPid:/ {print $2}' "$pid/status")
		sum_exec_runtime=$(awk '/^se.sum_exec_runtime/ {print $3}' "$pid/sched")
		nr_switches=$(awk '/^nr_switches/ {print $3}' "$pid/sched")

		if [ -n "$sum_exec_runtime" ] && [ -n "$nr_switches" ] && [ "$nr_switches" -ne 0 ]; then
			art=$(echo "scale=6; $sum_exec_runtime / $nr_switches" | bc | xargs printf "%.6f")
		else
			art="0.00000"
		fi

		echo "ProcessID=$pid_value : Parent_ProcessID=$ppid : Average_Running_Time=$art" >> fileToTask.txt
	fi
done

sort -t= -k3,3 fileToTask.txt -o fileToTask.txt

exit 0
