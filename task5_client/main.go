package main

import (
	"flag"
	"task5_client/client"
)

func main() {
	listenAddrString := flag.String("port", "127.0.0.1:8080", "Input a port to listen:")
	serverAddrString := flag.String("destserver", "127.0.0.1:8081", "Input remote server address:")
	agreeMent := flag.String("agree", "socks5", "Input a agreement:")
	flag.Parse()
	client.ListenLocal(*listenAddrString, *serverAddrString, *agreeMent)
}
