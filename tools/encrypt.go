package tools

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

const sKey = "PfHRR48%sFhw8*K1"
const sIv = "6nrcjzRDN:f8,C?f"

func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

func pkcs7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("加密字符串错误！")
	}
	unPadding := int(data[length-1])
	return data[:(length - unPadding)], nil
}

func AAesEncrypt(data []byte) ([]byte, error) {
	key := []byte(sKey)
	iv := []byte(sIv)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	encryptBytes := pkcs7Padding(data, blockSize)
	crypted := make([]byte, len(encryptBytes))
	blockMode := cipher.NewCBCEncrypter(block, iv)
	blockMode.CryptBlocks(crypted, encryptBytes)
	return crypted, nil
}

func AAesDecrypt(data []byte) ([]byte, error) {
	key := []byte(sKey)
	iv := []byte(sIv)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	crypted := make([]byte, len(data))
	blockMode.CryptBlocks(crypted, data)
	crypted, err = pkcs7UnPadding(crypted)
	if err != nil {
		return nil, err
	}
	return crypted, nil
}
