package client

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/url"
	"strings"
	"task5_client/config"
	"task5_client/core"
)

// type TCPClient struct {
// 	conn *net.TCPConn
// 	addr *net.TCPAddr
// }

func ListenLocal(listenAddrString, serverAddrString, agreeMentString string) error {
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
		go handleLocalClient(localClient, serverAddr, agreeMentString)
	}
}

func handleLocalClient(localClient *net.TCPConn, serverAddr *net.TCPAddr, argeeMentString string) {
	serverSocket, err := DialRemote(serverAddr)
	if err != nil {
		log.Fatal("远程服务器连接错误")
	}
	// serverSocket.SetKeepAlive(true)
	// defer serverSocket.Close()

	if argeeMentString == "http" {
		// // 开始建立socks5协议
		buffer := make([]byte, 256)
		// // 1.发送建立连接请求
		// _, err := serverSocket.Write([]byte{0x05, 0x01, 0x02})
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// // 读取回复
		// n, err := serverSocket.Read(buffer[:2])
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// if n == 0 {
		// 	log.Fatal("协议错误，服务器返回为空")
		// }
		// if buffer[1] == 0x02 && n == 2 {
		// 	log.Println("连接成功")
		// }
		// // 2.传输UNAME和PASSWD
		// uName := []byte(config.USERNAME)
		// passwd := []byte(config.PASSWORD)
		// uLen := len(uName)
		// pLen := len(passwd)
		// var temp []byte
		// temp = append(temp, 0x5, byte(uLen))
		// temp = append(temp, uName...)
		// temp = append(temp, byte(pLen))
		// temp = append(temp, passwd...)
		// serverSocket.Write(temp)
		// // 读取服务端的回复
		// _, err = serverSocket.Read(buffer[:2])
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// if buffer[1] != 0x0 {
		// 	log.Fatal("鉴权失败")
		// }

		// 开始读取请求消息
		// VER, CMD, RSV, ATYP, ADDR, PORT
		n, err := localClient.Read(buffer)
		if err != nil {
			log.Fatal(err)
		}
		localReq := buffer[:n]
		j := 0
		z := 0
		httpreq := []string{}
		for i := 0; i < n; i++ {
			if localReq[i] == 32 {
				httpreq = append(httpreq, string(buffer[j:i]))
				j = i + 1
			}
			if buffer[i] == 10 {
				z += 1
			}
		}
		dstURI, err := url.ParseRequestURI(httpreq[1])
		if err != nil {
			log.Fatal(err)
		}
		var dstAddr string
		var dstPort = "80"
		dstAddrPort := strings.Split(dstURI.Host, ":")
		if len(dstAddrPort) == 1 {
			dstAddr = dstAddrPort[0]
		} else if len(dstAddrPort) == 2 {
			dstAddr = dstAddrPort[0]
			dstPort = dstAddrPort[1]
		} else {
			log.Fatal("url parse error")
		}

		destAddrString := fmt.Sprintf("%s:%s", dstAddr, dstPort)
		destLen := len(destAddrString)
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
		// // 端口
		// dstPortBuff := bytes.NewBuffer(make([]byte, 0))
		// dstPortInt, err := strconv.ParseUint(dstPort, 10, 16)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// binary.Write(dstPortBuff, binary.BigEndian, dstPortInt)
		// dstPortBytes := dstPortBuff.Bytes() // int为8字节
		// resp = append(resp, dstPortBytes[len(dstPortBytes)-2:]...)

		_, err = core.EncodeWrite(serverSocket, buffer[:128])
		if err != nil {
			log.Println("client EncodeWrite err:" + err.Error())
		}
		_, err = core.EncodeWrite(serverSocket, buffer[128:])
		if err != nil {
			log.Println("client EncodeWrite err:" + err.Error())
		}
		// 读取回复消息
		// n, err = serverSocket.Read(resp)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// var targetResp [10]byte
		// copy(targetResp[:10], resp[:n])
		// specialResp := [10]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
		// if targetResp != specialResp {
		// 	log.Fatal("协议错误, 第二次协商返回出错")
		// }
		// log.Print("认证成功")
		// 转发消息
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
	} else {
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
	uName := config.USERNAME
	passwd := config.PASSWORD
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
