package main

import (
	"flag"
	"task5_server/server"
)

func main() {
	listenAddrString := flag.String("port", "127.0.0.1:8081", "Input a port to listen")
	flag.Parse()
	server.ListenServer(*listenAddrString)
}
