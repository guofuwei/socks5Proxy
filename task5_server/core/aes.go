package core

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"log"
	"task5_server/config"
)

const BlockSize = 128

var key = []byte(config.REQ_KEY)
var iv = []byte(config.REQ_IV)

func Encrypt(text []byte) ([]byte, error) {
	//生成cipher.Block 数据块
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println("错误 -" + err.Error())
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
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func Decrypt(decode_data []byte, sign int) ([]byte, error) {
	// decode_data, err := base64.StdEncoding.DecodeString(text)
	// if err != nil {
	// 	return nil, nil
	// }
	//生成密码数据块cipher.Block
	block, _ := aes.NewCipher(key)
	//解密模式
	blockMode := cipher.NewCBCDecrypter(block, iv)
	//输出到[]byte数组
	origin_data := make([]byte, len(decode_data))
	blockMode.CryptBlocks(origin_data, decode_data)
	// log.Println(origin_data)
	// log.Println(unpad(origin_data))
	//去除填充,并返回
	return unpad(origin_data, sign), nil
}

func unpad(ciphertext []byte, sign int) []byte {
	// log.Println("func unpad:")
	// log.Println(ciphertext)
	// length := len(ciphertext)
	// log.Printf("length:%d", length)
	//去掉最后一次的padding
	// unpadding := int(ciphertext[length-1])
	// return ciphertext[:(length - unpadding)]
	unpadSign := BlockSize - sign
	return ciphertext[:unpadSign]
}
