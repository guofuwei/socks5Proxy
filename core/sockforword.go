package core

import (
	"io"
	"net"
)

// 转发client<->server<->truedest
func Socks5Forward(client net.Conn, target net.Conn) {
	forward := func(src net.Conn, dest net.Conn) {
		defer src.Close()
		defer dest.Close()
		io.Copy(src, dest)
	}
	// 实现双向转发
	go forward(client, target)
	go forward(target, client)
}
