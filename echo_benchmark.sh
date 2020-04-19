#!/bin/bash

set -e

echo ""
echo "--- BENCH ECHO START ---"
echo ""

cd $(dirname "${BASH_SOURCE[0]}")
function cleanup {
    echo "--- BENCH ECHO DONE ---"
    kill -9 $(jobs -rp)
    wait $(jobs -rp) 2>/dev/null
}
trap cleanup EXIT

mkdir -p bin
$(pkill -9 gnet-echo-server || printf "")
$(pkill -9 evio-echo-server || printf "")
$(pkill -9 getty-echo-server || printf "")

function gobench {
    echo "--- $1 ---"
    if [ "$3" != "" ]; then
        go build -o $2 $3
    fi
    GOMAXPROCS=1 $2 --port $4 &
    sleep 1
    echo "*** 10 connections, 30 seconds, 6 byte packets"
    nl=$'\r\n'
    tcpkali --workers 1 -c 10 -T 30s -m "PING{$nl}" 127.0.0.1:$4
    echo "--- DONE ---"
    echo ""
}

gobench "GNET" bin/gnet-echo-server gnet-echo-server/main.go 5001
gobench "EVIO" bin/evio-echo-server evio-echo-server/main.go 5002
gobench "GETTY" bin/getty-echo-server getty-echo-server/main.go 5003