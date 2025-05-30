#!/usr/bin/bash

handler_pid=$(cat .pid)

while true; do
    read LINE
    case $LINE in
        "SHOW")
            kill -SIGHUP $handler_pid
            ;;
        "QUIT")
            kill -SIGTERM $handler_pid
            exit 0
            ;;
        *)
            if [[ $LINE =~ ^[0-9]+$ ]]; then
                if [ $LINE -lt 1 ] || [ $LINE -gt 20 ]; then
                    echo "The number must to be between 1 and 20"
                else
                    if [ $LINE -lt $(cat .secret) ]; then
                        kill -SIGUSR2 $handler_pid
                    elif [ $LINE -gt $(cat .secret) ]; then
                        kill -SIGUSR1 $handler_pid
                    else
			kill -SIGCONT $handler_pid	   
                    fi
                fi
            else
                echo "Uncorrect input date"
            fi
            ;;
    esac
done
exit 0

