package codec

import (
	"bytes"
	aes2 "crypto/aes"
	"crypto/cipher"
)

var AES aes

type aes struct {
}

func (ctr aes) AesEncryptCBC(origData []byte, key []byte) (encrypted []byte) {
	block, _ := aes2.NewCipher(key)
	blockSize := block.BlockSize()                              // 获取秘钥块的长度
	origData = ctr.pkcs5Padding(origData, blockSize)            // 补全码
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize]) // 加密模式
	encrypted = make([]byte, len(origData))                     // 创建数组
	blockMode.CryptBlocks(encrypted, origData)                  // 加密
	return encrypted
}
func (ctr aes) AesDecryptCBC(encrypted []byte, key []byte) (decrypted []byte) {
	block, _ := aes2.NewCipher(key)                             // 分组秘钥
	blockSize := block.BlockSize()                              // 获取秘钥块的长度
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize]) // 加密模式
	decrypted = make([]byte, len(encrypted))                    // 创建数组
	blockMode.CryptBlocks(decrypted, encrypted)                 // 解密
	decrypted = ctr.pkcs5UnPadding(decrypted)                   // 去除补全码
	return decrypted
}
func (aes) pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
func (aes) pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
