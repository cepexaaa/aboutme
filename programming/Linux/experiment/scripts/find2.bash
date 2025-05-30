#!/bin/bash

NEWMEM_SCRIPT="./scripts/newmem.bash"

K=30

LOW=1000000 
HIGH=100000000

check_success() {
    local N=$1
    echo "Check  N=$N..."

    for ((i=1; i<=K; i++)); do
        $NEWMEM_SCRIPT "$N" &
	sleep 1
    done
 
    wait
  
    if dmesg | grep -q "Out of memory"; then
        echo "OOM Killer by  N=$N"
        return 1
    else
        echo "All launches were completed successfully N=$N"
        return 0 
    fi
}


while [ "$LOW" -lt "$HIGH" ]; do

    MID=$(( (LOW + HIGH + 1) / 2 ))
 
    if check_success "$MID"; then
       
        LOW=$MID
    else
        HIGH=$((MID - 1))
    fi
done

echo "The maximum value of N at which all launches were completed successfully: $LOW"
