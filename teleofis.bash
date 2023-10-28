#!/bin/bash
echo 'Compiling for teleofis'
GOOS=linux GOARCH=arm  go build
if [ $? -ne 0 ]; then
	echo 'An error has occurred! Aborting the script execution...'
	exit 1
fi
echo 'Copy potop to device'
scp potop root@192.168.88.1:/root
#scp goirz.sh root@192.168.88.1:/root
# scp rc.local root@192.168.88.1:/etc
