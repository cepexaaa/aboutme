#!/bin/bash

echo "Total MEM:"
free -h | grep "Mem:" | awk '{print $2}'

echo "Total Swap:"
free -h | grep "Swap:" | awk '{print $2}'

echo "Size of page:"
getconf PAGE_SIZE

echo "Free MEM:"
free -h | grep "Mem:" | awk '{print $4}'

echo "Free Swap:"
free -h | grep "Swap:" | awk '{print $4}'
