#!/bin/ash
while true
do
    echo "start potop" >> start
    ./potop > /dev/null 2>/dev/null
    echo "need restart " >> start
done 
