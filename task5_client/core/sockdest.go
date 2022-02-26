package core

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

// 解析head获取真正的目标地址
func Socks5DestHandle(client *net.TCPConn) (string, error) {
	buffer := make([]byte, 256)
	// 先读四个字节VER,CMD,RSV,ATYP
	num, err := client.Read(buffer[:4])
	if num != 4 || err != nil {
		return "", errors.New("sockdest error :" + err.Error())
	}
	ver, cmd, _, atyp := buffer[0], buffer[1], buffer[2], buffer[3]
	// log.Print("desthandle:")
	// log.Println(buffer)
	// log.Printf("ver:%d,cmd:%d,atyp:%d", ver, cmd, atyp)
	if ver != 5 || cmd != 1 {
		// 目前只支持connect方式
		return "", errors.New("invalid socks5 version/cmd")
	}
	// 开始解析ipaddr
	// log.Printf("The atyp:%d\n", atyp)
	addr := ""
	switch atyp {
	case 1:
		// IPV4 address
		num, err = client.Read(buffer[:4])
		if num != 4 || err != nil {
			return "", errors.New("invalid ipv4 address :" + err.Error())
		}
		addr = fmt.Sprintf("%d.%d.%d.%d", buffer[0], buffer[1], buffer[2], buffer[3])
	case 3:
		// host域名解析,hostAddr和host
		num, err = client.Read(buffer[:1])
		if num != 1 || err != nil {
			return "", errors.New("sockdest error :" + err.Error())
		}
		hostLen := int(buffer[0])
		num, err = client.Read(buffer[:hostLen])
		if num != hostLen || err != nil {
			return "", errors.New("sockdest error :" + err.Error())
		}
		addr = string(buffer[:hostLen])
	case 4:
		// IPV6 address 暂时不支持
		return "", errors.New("IPV6:not support yet")
	default:
		return "", errors.New("invalid aytp")
	}
	// 开始解析port
	num, err = client.Read(buffer[:2])
	if num != 2 || err != nil {
		return "", errors.New("sockdest error :" + err.Error())
	}
	port := binary.BigEndian.Uint16(buffer[:2])
	// addr和port都已就绪，开始生成dest的conn
	destAddrPort := fmt.Sprintf("%s:%d", addr, port)
	// dest, err := net.Dial("tcp", destAddrPort)
	// if err != nil {
	// 	return nil, errors.New("dial error:" + err.Error())
	// }
	// log.Println("The destAddrPort:" + destAddrPort)
	// 回复浏览器
	_, err = client.Write([]byte{0x05, 0x00, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
	if err != nil {
		client.Close()
		return "", errors.New("write conn rsp: " + err.Error())
	}
	return destAddrPort, nil
}
