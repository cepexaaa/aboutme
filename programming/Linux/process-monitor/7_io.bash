#!/bin/bash



collect_io_data() {
	for pid in /proc/[0-9]*; do
		if [ -f "$pid/io" ]; then
#	for pid in $(ps -u $(id -u) -o pid=); do
#		if [ -d "/proc/$pid" ]; then
			pid_value=$(basename "$pid")
			read_bytes=$(sudo awk '/^read_bytes:/ {print $2}' "/proc/$pid_value/io")
			cmdline=$(sudo tr -d '\0' < "/proc/$pid_value/cmdline" | sed 's/ *$//')
			echo "$pid_value:$cmdline:$read_bytes"
		fi
	done
}

start_data=$(collect_io_data)

sleep 60

end_data=$(collect_io_data)

echo "$start_data" | while IFS=: read -r pid cmdline start_bytes; do
	end_bytes=$(echo "$end_data" | awk -F: -v pid="$pid" '$1 == pid {print $3}')
	if [ -n "$end_bytes" ]; then
		bytes_read=$((end_bytes - start_bytes))
		echo "$pid:$cmdline:$bytes_read"
	fi
done | sort -t: -k3,3nr | head -n 3

exit 0
