#!/bin/sh
git stash
git pull
sleep 3
docker restart infoblogserver_"$1"