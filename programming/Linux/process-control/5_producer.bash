#!/usr/bin/bash

LINE=""

echo $$ > .pid5

if [[ ! -p wire ]]; then
    mkfifo pipe
fi

while true; do
    read LINE
    echo "$LINE" > pipe
done


