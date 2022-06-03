#!/bin/bash
pushd server
go run main.go &
popd
pushd client
go run main.go &
echo "finish to run the server and client"
