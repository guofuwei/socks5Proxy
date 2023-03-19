package core

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"log"

	"socks5proxy"
)

const BlockSize = 128

var key = []byte(socks5proxy.REQ_KEY)
var iv = []byte(socks5proxy.REQ_IV)

func Encrypt(text []byte) ([]byte, error) {
	//生成cipher.Block 数据块
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println("error -" + err.Error())
		return nil, err
	}
	//填充内容，如果不足16位字符
	originData := pad(text, BlockSize)
	//加密方式
	blockMode := cipher.NewCBCEncrypter(block, iv)
	//加密，输出到[]byte数组
	crypted := make([]byte, len(originData))
	blockMode.CryptBlocks(crypted, originData)
	// log.Println(crypted)
	return crypted, nil
}

func pad(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)
	padtext := bytes.Repeat([]byte{0x0}, padding)
	return append(ciphertext, padtext...)
}

func Decrypt(decode_data []byte, sign int) ([]byte, error) {
	//生成密码数据块cipher.Block
	block, _ := aes.NewCipher(key)
	//解密模式
	blockMode := cipher.NewCBCDecrypter(block, iv)
	//输出到[]byte数组
	// log.Printf("Len:%d", len(decode_data))
	origin_data := make([]byte, len(decode_data))
	blockMode.CryptBlocks(origin_data, decode_data)
	// log.Println(origin_data)
	// log.Println(unpad(origin_data))
	//去除填充,并返回
	return unpad(origin_data, sign), nil
}

func unpad(ciphertext []byte, sign int) []byte {
	unpadSign := BlockSize - sign
	return ciphertext[:unpadSign]
}
