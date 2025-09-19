package security

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"os"
)

// This is IV, TODO: Change this to be unique for all different objects and save it to DB
var bytes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

func Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func Decode(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}

func Encrypt(textToEncrypt string) (string, error) {
	block, err := aes.NewCipher([]byte(os.Getenv("AES_SECRET")))
	if err != nil {
		return "", err
	}
	plainText := []byte(textToEncrypt)
	cfb := cipher.NewCFBEncrypter(block, bytes)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return Encode(cipherText), nil
}

func Decrypt(text string) (string, error) {
	block, err := aes.NewCipher([]byte(os.Getenv("AES_SECRET")))
	if err != nil {
		return "", err
	}
	cipherText := Decode(text)
	cfb := cipher.NewCFBDecrypter(block, bytes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}
