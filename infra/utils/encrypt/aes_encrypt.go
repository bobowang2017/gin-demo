package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

// AesEncrypt key(either 16, 24, or 32 bytes to select)
func AesEncrypt(orig string, key string) (string, error) {
	keyLen := len(key)
	if !(keyLen == 16 || keyLen == 24 || keyLen == 32) {
		return "", errors.New("AesEncrypt key(either 16, 24, or 32 bytes to select)")
	}
	// 转成字节数组
	origData := []byte(orig)
	k := []byte(key)

	// 分组秘钥
	block, _ := aes.NewCipher(k)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 补全码
	origData = PKCS7Padding(origData, blockSize)
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
	// 创建数组
	cryptArr := make([]byte, len(origData))
	// 加密
	blockMode.CryptBlocks(cryptArr, origData)

	return base64.StdEncoding.EncodeToString(cryptArr), nil

}

func AesDecrypt(cryptStr string, key string) (string, error) {
	// 转成字节数组
	cryptByte, _ := base64.StdEncoding.DecodeString(cryptStr)
	k := []byte(key)

	// 分组秘钥
	block, err := aes.NewCipher(k)
	if err != nil {
		return "", err
	}
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 加密模式
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
	// 创建数组
	orig := make([]byte, len(cryptByte))
	// 解密
	blockMode.CryptBlocks(orig, cryptByte)
	// 去补全码
	orig, err = PKCS7UnPadding(orig)
	if err != nil {
		return "", err
	}
	return string(orig), nil
}

// PKCS7Padding 补码
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

// PKCS7UnPadding 去码
func PKCS7UnPadding(origData []byte) ([]byte, error) {
	length := len(origData)
	if length == 0 {
		return nil, errors.New("加密字符串错误")
	}
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)], nil
}
