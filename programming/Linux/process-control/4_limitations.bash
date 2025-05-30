#!/usr/bin/bash

first_pid=$(pgrep -f for_4_task.bash | head -n 1)
cpulimit -p $first_pid -l 10 &
third_pid=$(pgrep -f for_4_task.bash | tail -n 1)

kill $third_pid 

top -p $first_pid

exit 0
