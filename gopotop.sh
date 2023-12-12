#!/bin/ash
mkdir /tmp/rura
cp config.json /tmp/rura/
cd /cache/rura
tar -xzvf potop.tar.gz -C "/tmp/rura"
cd /tmp/rura
while true
do
    ./potop > /dev/null 2>/dev/null
    cp config.json /root
done 
