#!/bin/bash
echo 'Compiling for Kazan'
GOOS=linux GOARCH=arm  go build
if [ $? -ne 0 ]; then
	echo 'An error has occurred! Aborting the script execution...'
	exit 1
fi
echo 'Copy potop to device Kazan'
tar -czvf potop.tar.gz potop
scp -P 222 potop.tar.gz root@185.27.195.194:/root 
scp -P 222 gopotop.sh root@185.27.195.194:/root 

#scp goirz.sh root@192.168.88.1:/root
# scp rc.local root@192.168.88.1:/etc
