#!/bin/sh
git stash
git pull
docker restart infoblogserver_"$1"