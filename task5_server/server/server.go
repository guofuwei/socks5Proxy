package server

import (
	"errors"
	"log"
	"net"
	"task5_server/config"
	"task5_server/core"
)

func ListenServer(listenAddrString string) {
	listenAddr, err := net.ResolveTCPAddr("tcp", listenAddrString)
	if err != nil {
		log.Println(err)
		return
	}
	serverListener, err := net.ListenTCP("tcp", listenAddr)
	if err != nil {
		log.Fatal("服务器监听端口错误")
	}
	defer serverListener.Close()
	log.Println("server start successed!")
	for {
		serverClient, err := serverListener.AcceptTCP()
		// println("a new connect establish")
		if err != nil {
			log.Println(err)
			continue
		}
		serverClient.SetLinger(0)
		go handleServerClient(serverClient)
	}
}

func handleServerClient(serverClient *net.TCPConn) {
	// 先进行权限的认证
	err := localAuthHandle(serverClient)
	if err != nil {
		log.Println(err)
		serverClient.Close()
		return
	}
	// 获取目标地址
	destSocket, err := localDestHandle(serverClient)
	if err != nil {
		log.Println(err)
		serverClient.Close()
		return
	}
	// 双向转发
	go func() {
		err = core.EncodeCopy(destSocket, serverClient)
		if err != nil {
			log.Println("server EncodeCopy err:" + err.Error())
		}
		destSocket.Close()
	}()
	go func() {
		err = core.DecodeCopy(serverClient, destSocket)
		if err != nil {
			log.Println("server DecodeCpoy err:" + err.Error())
		}
		serverClient.Close()
	}()
}

func localAuthHandle(serverClient *net.TCPConn) error {
	// 开始进行身份验证
	// 读取uLen,uName,pLen,passwd
	buffer := make([]byte, 128)
	_, err := serverClient.Read(buffer[:1])
	if err != nil {
		return err
	}
	uLen := int(buffer[0])
	n, err := serverClient.Read(buffer[:uLen])
	if err != nil {
		return err
	}
	if n != uLen {
		return errors.New("read uname error")
	}
	uName := string(buffer[:uLen])

	_, err = serverClient.Read(buffer[:1])
	if err != nil {
		return err
	}
	pLen := int(buffer[0])
	n, err = serverClient.Read(buffer[:pLen])
	if err != nil {
		return err
	}
	if n != pLen {
		return errors.New("read passwd error")
	}
	passwd := string(buffer[:pLen])

	// 读取uName和passwd完成，开始比对uName和passwd并回复客户端
	if uName == config.USERNAME && passwd == config.PASSWORD {
		serverClient.Write([]byte{0x0})
	} else {
		serverClient.Write([]byte{0x1})
	}
	return nil
}

func localDestHandle(serverClient *net.TCPConn) (*net.TCPConn, error) {
	// 读取destLen,destAddr
	buffer := make([]byte, 128)
	_, err := serverClient.Read(buffer[:1])
	if err != nil {
		return nil, errors.New("get destLen error :" + err.Error())
	}
	destLen := int(buffer[0])
	// log.Printf("destLen:%d", destLen)
	n, err := serverClient.Read(buffer[:destLen])
	if err != nil {
		return nil, errors.New("read destAddr error :" + err.Error())
	}
	if n != destLen {
		return nil, errors.New("read destAddr length error")
	}
	destAddrString := string(buffer[:destLen])
	destAddr, err := net.ResolveTCPAddr("tcp", destAddrString)
	log.Printf("destAddr is:%v", destAddr)
	if err != nil {
		return nil, err
	}
	destSocket, err := net.DialTCP("tcp", nil, destAddr)
	if err != nil {
		return nil, err
	}
	// 返回拿取destAddr成功
	_, err = serverClient.Write([]byte{0x0})
	if err != nil {
		return nil, errors.New("get destAddr error:" + err.Error())
	}
	return destSocket, nil
}
