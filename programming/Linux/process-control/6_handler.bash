#!/usr/bin/bash

echo $$ > .pid

A=1

usr1() {
    let A=$A+2
    echo "Current value: $A (incremented by 2)"
}

usr2() {
    let A=$A*2
    echo "Current value: $A (multiplied by 2)"
}

term() {
    echo "Handler stopped by SIGTERM"
    rm .pid
    exit 0
}

trap 'usr1' USR1
trap 'usr2' USR2
trap 'term' SIGTERM

while true; do
    sleep 1
done
exit 0
