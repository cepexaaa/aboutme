#!/bin/bash

./scripts/mem.bash &

./scripts/mem2.bash &

./collect_top.bash &

wait
