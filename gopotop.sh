#!/bin/ash
mkdir /tmp/rura
tar -xzvf potop.tar.gz -C "/tmp/rura"
cp config.json /tmp/rura/
cd /tmp/rura
while true
do
    echo "start potop" >> start
    ./potop > /dev/null 2>/dev/null
    cp config.json /root
    echo "need restart " >> start
done 
