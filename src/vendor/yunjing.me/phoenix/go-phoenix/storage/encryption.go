package storage

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"log"
)

var (
	iv  = []byte{0x07, 0x0f, 0x04, 0x0d, 0x06, 0x0b, 0x0a, 0x00, 0x0e, 0x09, 0x02, 0x05, 0x08, 0x03, 0x0c, 0x01}
	key = []byte("Msn4liow84c2kz4lns9qbo42iz4nj2na")
)

func Encrypt(text string) string {
	c, err := aes.NewCipher(key)
	if err != nil {
		log.Panicf("生成数据加密工具时出错: %v", err)
	}

	dec := cipher.NewCFBEncrypter(c, iv)
	dst := make([]byte, len(text))
	dec.XORKeyStream(dst, []byte(text))

	return base64.StdEncoding.EncodeToString(dst)
}

func Decrypt(text string) string {
	c, err := aes.NewCipher(key)
	if err != nil {
		log.Panicf("生成数据解密工具时出错: %v", err)
	}

	raw, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		log.Panicf("数据格式错误: %v", err)
	}

	dec := cipher.NewCFBDecrypter(c, iv)
	dst := make([]byte, len(raw))
	dec.XORKeyStream(dst, raw)
	return string(dst)
}
