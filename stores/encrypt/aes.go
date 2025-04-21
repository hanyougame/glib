package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

// AESEncrypt AES加密
func AESEncrypt(plainText string, secretKey []byte) (string, error) {
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}

	// 填充数据
	padding := aes.BlockSize - len(plainText)%aes.BlockSize
	plainText = plainText + string(bytes.Repeat([]byte{byte(padding)}, padding))

	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]

	// 随机生成初始化向量
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[aes.BlockSize:], []byte(plainText))

	// 将加密后的数据转为base64
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// AESDecrypt AES解密
func AESDecrypt(cipherTextBase64 string, secretKey []byte) (string, error) {
	cipherText, err := base64.StdEncoding.DecodeString(cipherTextBase64)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherText, cipherText)

	// 去除填充
	padding := int(cipherText[len(cipherText)-1])
	cipherText = cipherText[:len(cipherText)-padding]

	return string(cipherText), nil
}
