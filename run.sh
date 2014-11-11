#!/bin/bash
while true; do
    /gopath/bin/app "$@"
    sleep 60
done
