#!/usr/bin/bash

handler_pid=$(cat .pid)

while true; do
    read LINE
    case $LINE in
        "+")
            kill -USR1 $handler_pid
            ;;
        "*")
            kill -USR2 $handler_pid
            ;;
        "TERM")
            kill -SIGTERM $handler_pid
            exit 0
            ;;
        *)
            :
            ;;
    esac
done
exit 0
