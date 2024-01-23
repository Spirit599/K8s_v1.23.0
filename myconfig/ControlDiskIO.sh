#!/bin/sh

speed=$1

while true; do
    start=$(date +%s.%N)
    dd if=/dev/zero of=/tmp/write_to bs=1M count=$speed status=none
    end=$(date +%s.%N)
    duration=$(echo "$end - $start" | bc)
    diff=$(echo "print 0; 1 - $duration" | bc)
    sleep $diff
done