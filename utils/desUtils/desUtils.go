package desUtils

import (
	"bytes"
	"crypto/aes"
	"fmt"
)

var DefaultKey = []byte(`arthur#0409.5cc98&3008E#CD^E542@`)

func DefaultECBDecrypt(crypted []byte) ([]byte, error) {
	return ECBDecrypt(crypted, DefaultKey)
}

// ECB模式解密
func ECBDecrypt(crypted, key []byte) ([]byte, error) {
	if !validKey(key) {
		return nil, fmt.Errorf("秘钥长度错误,当前传入长度为 %d", len(key))
	}
	if len(crypted) < 1 {
		return nil, fmt.Errorf("源数据长度不能为0")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(crypted)%block.BlockSize() != 0 {
		return nil, fmt.Errorf("源数据长度必须是 %d 的整数倍，当前长度为：%d", block.BlockSize(), len(crypted))
	}
	var dst []byte
	tmpData := make([]byte, block.BlockSize())

	for index := 0; index < len(crypted); index += block.BlockSize() {
		block.Decrypt(tmpData, crypted[index:index+block.BlockSize()])
		dst = append(dst, tmpData...)
	}
	dst, err = PKCS5UnPadding(dst)
	if err != nil {
		return nil, err
	}
	return dst, nil
}
func DefaultECBEncrypt(src []byte) ([]byte, error) {
	return ECBEncrypt(src, DefaultKey)
}

// ECB模式加密
func ECBEncrypt(src, key []byte) ([]byte, error) {
	if !validKey(key) {
		return nil, fmt.Errorf("秘钥长度错误, 当前传入长度为 %d", len(key))
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(src) < 1 {
		return nil, fmt.Errorf("源数据长度不能为0")
	}
	src = PKCS5Padding(src, block.BlockSize())
	if len(src)%block.BlockSize() != 0 {
		return nil, fmt.Errorf("源数据长度必须是 %d 的整数倍，当前长度为：%d", block.BlockSize(), len(src))
	}
	var dst []byte
	tmpData := make([]byte, block.BlockSize())
	for index := 0; index < len(src); index += block.BlockSize() {
		block.Encrypt(tmpData, src[index:index+block.BlockSize()])
		dst = append(dst, tmpData...)
	}
	return dst, nil
}

// PKCS5填充
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// 去除PKCS5填充
func PKCS5UnPadding(origData []byte) ([]byte, error) {
	length := len(origData)
	unpadding := int(origData[length-1])

	if length < unpadding {
		return nil, fmt.Errorf("invalid unpadding length")
	}
	return origData[:(length - unpadding)], nil
}

// 秘钥长度验证
func validKey(key []byte) bool {
	k := len(key)
	switch k {
	default:
		return false
	case 16, 24, 32:
		return true
	}
}
