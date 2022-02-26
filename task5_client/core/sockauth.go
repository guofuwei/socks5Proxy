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
	// log.Print("sockauth:")
	// log.Println(buffer)
	// log.Printf("ver:%d,nMethods:%d\n", ver, nMethods)
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
	// // 开始解析username和password
	// // 1.首先读取VER和ULEN
	// num, err = client.Read(buffer[:2])
	// if num != 2 || err != nil {
	// 	return errors.New("sockauth error :" + err.Error())
	// }
	// ver, uLen := int(buffer[0]), int(buffer[1])
	// if ver != 5 {
	// 	return errors.New("invalid socks5 version")
	// }
	// // 2.读取UNAME并进行验证
	// num, err = client.Read(buffer[:uLen])
	// if num != uLen || err != nil {
	// 	return errors.New("sockauth error :" + err.Error())
	// }
	// uName := string(buffer[:uLen])
	// if uName != config.USERNMAE {
	// 	return errors.New("username or password not right")
	// }
	// // 3.读取PLEN
	// num, err = client.Read(buffer[:1])
	// if num != 1 || err != nil {
	// 	return errors.New("sockauth error :" + err.Error())
	// }
	// pLen := int(buffer[0])
	// // 4.读取PASSRD
	// num, err = client.Read(buffer[:pLen])
	// if num != pLen || err != nil {
	// 	return errors.New("sockauth error :" + err.Error())
	// }
	// passwd := string(buffer[:pLen])
	// if passwd != config.PASSWORD {
	// 	return errors.New("username or password not right")
	// }
	return nil
}
