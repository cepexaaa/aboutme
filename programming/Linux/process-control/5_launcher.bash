#!/usr/bin/bash

./5_handler.bash &
handler_pid=$!

./5_producer.bash
producer_pid=$!

wait $handler_pid $producer_pid

exit 0
