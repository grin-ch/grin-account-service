#!/bin/bash
# auther: grin

mode=$1

if [ mode=="dev" ]
then 
    go version
    go run ./cmd/main.go --cfgName=dev
elif [ mode=="serve" ]
then
    go run ./cmd/main.go --cfgName=service
fi