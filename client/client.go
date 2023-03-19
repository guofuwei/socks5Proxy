package client

import (
	"errors"
	"log"
	"net"
	"socks5proxy"
	"socks5proxy/core"
	"sync"
)

func ListenLocal(listenAddrString, serverAddrString string) error {
	listenAddr, err := net.ResolveTCPAddr("tcp", listenAddrString)
	if err != nil {
		log.Fatal("resolve local address error")
	}
	serverAddr, err := net.ResolveTCPAddr("tcp", serverAddrString)
	if err != nil {
		log.Fatal("resolve server address error")
	}

	localListener, err := net.ListenTCP("tcp", listenAddr)
	if err != nil {
		return err
	}
	defer localListener.Close()
	log.Printf("local start successed,listen on %s\n", listenAddrString)

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
		// log.Print("remote server connect error")
		localClient.Close()
		return
	}
	//这里是socks5协议，应该先截取socks5头部然后读取destAddr
	buffer := make([]byte, 128)
	err = core.Socks5AuthHandle(localClient)
	if err != nil {
		// log.Print("auth error " + err.Error())
		serverSocket.Close()
		localClient.Close()
		return
	}
	destAddrString, err := core.Socks5DestHandle(localClient)
	if err != nil {
		// log.Print("get destAddr error:" + err.Error())
		serverSocket.Close()
		localClient.Close()
		return
	}
	//开始使用自己的协议与server通信
	// 1.发送host长度帧，发送host帧
	destLen := len(destAddrString)
	_, err = serverSocket.Write([]byte{byte(destLen)})
	if err != nil {
		// log.Println("send destLen error:" + err.Error())
		return
	}
	_, err = serverSocket.Write([]byte(destAddrString))
	if err != nil {
		// log.Println("send destAddr error:" + err.Error())
		return
	}
	_, err = serverSocket.Read(buffer[:1])
	if err != nil {
		// log.Println("server can't get the destSocket:" + err.Error())
		return
	}
	// 2.转发消息
	wg := new(sync.WaitGroup)
	wg.Add(2)
	go func() {
		defer wg.Done()
		defer localClient.Close()
		core.EncodeCopy(localClient, serverSocket)
	}()
	go func() {
		defer wg.Done()
		defer serverSocket.Close()
		core.DecodeCopy(serverSocket, localClient)
	}()
	wg.Wait()
}

func DialRemote(serverAddr *net.TCPAddr) (serverSocket *net.TCPConn, err error) {
	// 1.先与服务器建立tcp连接
	serverSocket, err = net.DialTCP("tcp", nil, serverAddr)
	if err != nil {
		log.Fatal("remote server connect error")
	}
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
