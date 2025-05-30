#!/usr/bin/bash
value=1
operation="add"

if [[ ! -p pipe ]]; then
    mkfifo pipe
fi

line=""

(tail -f pipe & echo $!>.tail_pid) |
	while true; do
		read line;
		if [ -z "$line" ]; then
    			continue
		fi
		case "$line" in
			QUIT)
				echo "exit"
				kill $(cat .tail_pid)
      				rm .tail_pid
      				kill $(cat .pid5)
      				rm .pid5
				exit 0
				;;
			+)
				operation="add"
				echo "change to oparetion additional"
				;;
			\*)
				operation="mult"
				echo "change to operation multiplication"
				;;
			*)
				if [[ $line =~ ^-?[0-9]+$ ]]; then
					if [ "$operation" == "add" ]; then
						value=$((value + line))
					elif [ "$operation" == "mult" ]; then
						value=$((value * line))
					fi
					echo "current value: $value"
				else
					if [ -n "$line" ]; then
						echo "Uncorrect input data: '$line'"
						kill $(cat .tail_pid)
    						rm .tail_pid
    						kill $(cat .pid5)
    						rm .pid5
						exit 1
					fi
				fi
				;;


	esac
done

exit 0
