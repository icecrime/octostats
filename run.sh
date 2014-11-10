#!/bin/bash
while true; do
    /gopath/bin/app "$@" | nc -q0 $CARBON_PORT_2003_TCP_ADDR $CARBON_PORT_2003_TCP_PORT
    sleep 60
done
