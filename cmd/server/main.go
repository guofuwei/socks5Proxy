package main

import (
	"flag"
	"socks5proxy/server"
)

func main() {
	listenAddrString := flag.String("local", "127.0.0.1:8081", "Input a address to listen")
	flag.Parse()
	server.ListenServer(*listenAddrString)
}
