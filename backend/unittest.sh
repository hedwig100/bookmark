#/bin/bash

dir=$1

if [ "$dir" = "server" ]; then
    go test ./server/ -v -run TestUnit*
else 
    go test ./data/ -v
fi

code=$?
if [ "$code" = "0" ]; then 
    echo "OK"
else 
    echo "Failure"
fi