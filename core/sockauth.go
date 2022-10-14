package core

import (
	"errors"
	"net"
)

// 该函数用来完成socks5的权限认证问题

func Socks5AuthHandle(client *net.TCPConn) error {
	buffer := make([]byte, 256)
	// 处理第一个建立连接请求头
	num, err := client.Read(buffer[:2])
	if err != nil || num != 2 {
		return errors.New("sockauth error :" + err.Error())
	}
	// 1.首先读取VER和NMETHODS
	ver, nMethods := int(buffer[0]), int(buffer[1])
	if ver != 5 {
		return errors.New("invalid socks5 version")
	}
	// 2.根据NMENTHODS读取NMETHOD
	num, err = client.Read(buffer[:nMethods])
	if num != nMethods {
		return errors.New("sockauth error :" + err.Error())
	}

	// 开始响应浏览器不需要权限认证
	num, err = client.Write([]byte{0x05, 0x00})
	if num != 2 || err != nil {
		return errors.New("sockauth error :" + err.Error())
	}
	return nil
}
