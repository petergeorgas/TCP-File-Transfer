# TCP File Transfer
This repository houses the code demonstrating a basic file transfer client and server. 
I mostly created this project just to get more comfortable with the `net` Go package, and also to develop a better understanding of the TCP protocol.

## Usage
Start the server
```bash
go run server/main.go
```

To send file(s) to the server, navigate to `client/` and run `go run main.go tcp-transfer <file_path(s)>`, `<file_paths` may be realtive or absolute. The file(s) will appear in the `server/transferred` directory once the transfer is complete. 

## The Custom Protocol
A (simple) custom protocol was developed for this application. When a file is transferred, the following are sent (in order):
1. Length of the file name
2. The file name
3. The size of the file being sent
4. The (buffered) contents of the file

If the server does not get the above information, it would not be able to properly transfer the file. If you would like to use just the server portion of this application and write your own client, feel free to do so. 

