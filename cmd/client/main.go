package main

import (
	"flag"
	"socks5proxy/client"
)

func main() {
	listenAddrString := flag.String("local", "127.0.0.1:8080", "Input a address to listen:")
	serverAddrString := flag.String("dest", "127.0.0.1:8081", "Input remote server address:")
	flag.Parse()
	client.ListenLocal(*listenAddrString, *serverAddrString)
}
