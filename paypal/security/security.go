package security

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"math/rand"
	"os"
)

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

func Encrypt(textToEncrypt string, bytes []byte) (string, error) {
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

func Decrypt(text string, bytes []byte) (string, error) {
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

func CreateIV() []byte {
	max := 100
	min := 1
	vector := []byte{}
	for i := 0; i < 16; i++ {
		randInt := rand.Intn(max-min) + min
		vector = append(vector, byte(randInt))
	}
	return vector
}
