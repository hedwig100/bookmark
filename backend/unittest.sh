#/bin/bash

ok() {
    code=$?
    if [ "$code" = "0" ]; then 
        echo "OK"
    else 
        echo "Failure"
    fi
}

dir=$1

if [ "$dir" = "" ]; then 
    go test ./server/ -v -run TestUnit*; ok
    go test ./data/ -v; ok
elif [ "$dir" = "server" ]; then
    go test ./server/ -v -run TestUnit*; ok
else 
    go test ./$dir/ -v; ok
fi