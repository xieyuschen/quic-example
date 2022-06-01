# Echo cli
The echo-cli provides a more controllable way for you to try it. You will learn some quic behaviors during 
try it. I will also discuss some details found in programming.
Running `go run main.go` both in server and client could get a echo demo of quic.

## The key pair and certification
// todo

## A good wrapper with io.copy
// todo

## io.ReadFull or io.Read when reading a stream?
// todo

## Difference between OpenStream and OpenStreamSync?
// todo

## What happened after the server side finish write data in a stream?
// todo 
## Why the server could only serve one client request?
The first client request gets a echo message, but the second one blocks until gets a "timeout: 
no recent network activity" error.  
So why the server doesn't exit but cannot serve a new client request?
