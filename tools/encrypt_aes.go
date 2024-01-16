package tools

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

func DefaultEncryptString(plainText string) string {
	cipherText, _ := EncryptString(plainText, "PfHRR48%sFhw8*K1")
	return cipherText
}

func DefaultDecryptString(plainText string) string {
	cipherText, _ := DecryptString(plainText, "PfHRR48%sFhw8*K1")
	return cipherText
}

/// EncryptString 使用 AES 加密字符串
func EncryptString(plainText, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	plainTextBytes := []byte(plainText)
	blockSize := block.BlockSize()
	plainTextBytes = PKCS7Padding(plainTextBytes, blockSize)
	cipherText := make([]byte, len(plainTextBytes))
	mode := cipher.NewCBCEncrypter(block, []byte(key[:blockSize]))
	mode.CryptBlocks(cipherText, plainTextBytes)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// DecryptString 使用 AES 解密字符串
func DecryptString(cipherText, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	cipherTextBytes, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	if len(cipherTextBytes) < blockSize {
		return "", errors.New("cipherText too short")
	}

	mode := cipher.NewCBCDecrypter(block, []byte(key[:blockSize]))
	mode.CryptBlocks(cipherTextBytes, cipherTextBytes)
	cipherTextBytes = PKCS7UnPadding(cipherTextBytes)
	return string(cipherTextBytes), nil
}

// PKCS7Padding 填充字符串
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS7UnPadding 删除填充
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
