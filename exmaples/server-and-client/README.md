## Server-and-client
The second case here provides separate client and server instances.
Running this example requires starting the server and client separately.
- start server  
  The server provides some flags for you.
  If you run it without flags, it serves on `https://localhost:6121/` and 
  uses certificated files under [server/testdata](./server/testdata).

- start client  
  Use args to specify the destination. If no args, the client make a request to `https://localhost:6121/` by default.
```shell
  # under server-and-client/server
  go run main.go https://localhost:8080
```
The client ca file is stored under [client/testdata](./client/testdata).