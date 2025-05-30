#!/bin/bash

N=$1
K=$2

for ((i=1; i<=K; i++)); do
    echo "Script $i of $K with N=$N"
    ./scripts/newmem.bash "$N" &
    sleep 1
done

wait
echo "All scripts finished succesful."
