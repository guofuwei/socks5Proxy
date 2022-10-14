package client

import (
	"errors"
	"log"
	"net"
	"socks5proxy"
	"socks5proxy/core"
)

func ListenLocal(listenAddrString, serverAddrString string) error {
	listenAddr, err := net.ResolveTCPAddr("tcp", listenAddrString)
	if err != nil {
		log.Fatal(err)
	}
	serverAddr, err := net.ResolveTCPAddr("tcp", serverAddrString)
	if err != nil {
		log.Fatal(err)
	}

	localListener, err := net.ListenTCP("tcp", listenAddr)
	if err != nil {
		return err
	}
	defer localListener.Close()
	log.Println("local start successed!")

	for {
		localClient, err := localListener.AcceptTCP()
		if err != nil {
			return err
		}
		go handleLocalClient(localClient, serverAddr)
	}
}

func handleLocalClient(localClient *net.TCPConn, serverAddr *net.TCPAddr) {
	serverSocket, err := DialRemote(serverAddr)
	if err != nil {
		log.Fatal("远程服务器连接错误")
	}
	//这里是socks5协议，应该先截取socks5头部然后读取destAddr
	buffer := make([]byte, 128)
	err = core.Socks5AuthHandle(localClient)
	if err != nil {
		log.Fatal("auth error" + err.Error())
	}
	destAddrString, err := core.Socks5DestHandle(localClient)
	if err != nil {
		log.Fatal("get destAddr error:" + err.Error())
	}
	//开始使用自己的协议与server通信
	// 1.发送host长度帧，发送host帧
	destLen := len(destAddrString)
	// log.Printf("destLen:%d", destLen)
	// log.Printf("destAddrString:%s", destAddrString)
	_, err = serverSocket.Write([]byte{byte(destLen)})
	if err != nil {
		log.Println("send destLen error:" + err.Error())
		return
	}
	_, err = serverSocket.Write([]byte(destAddrString))
	if err != nil {
		log.Println("send destAddr error:" + err.Error())
		return
	}
	_, err = serverSocket.Read(buffer[:1])
	if err != nil {
		log.Println("server can't get the destSocket:" + err.Error())
		return
	}
	// 2.转发消息
	go func() {
		err = core.EncodeCopy(localClient, serverSocket)
		if err != nil {
			log.Println("client EncodeCpoy err:" + err.Error())
		}
		localClient.Close()
	}()
	go func() {
		err = core.DecodeCopy(serverSocket, localClient)
		if err != nil {
			log.Println("client DecodeCpoy err:" + err.Error())
		}
		serverSocket.Close()
	}()
}

func DialRemote(serverAddr *net.TCPAddr) (serverSocket *net.TCPConn, err error) {
	// 1.先与服务器建立tcp连接
	serverSocket, err = net.DialTCP("tcp", nil, serverAddr)
	if err != nil {
		log.Fatal("远程服务器连接错误")
	}
	// defer serverSocket.Close()
	// 2.开始进行身份验证
	// 发送uLen,uName,pLen,passwd
	uName := socks5proxy.USERNAME
	passwd := socks5proxy.PASSWORD
	uLen := len(uName)
	pLen := len(passwd)
	serverSocket.Write([]byte{byte(uLen)})
	serverSocket.Write([]byte(uName))
	serverSocket.Write([]byte{byte(pLen)})
	serverSocket.Write([]byte(passwd))
	// 开始接受身份确认信息
	buffer := make([]byte, 128)
	serverSocket.Read(buffer[:1])
	if buffer[0] != 0x0 {
		return nil, errors.New("remote server auth fail")
	}
	return serverSocket, nil
}
