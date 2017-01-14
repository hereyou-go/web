package aes

import (
	stdaes "crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

const DEFAULT_KEY = "abcdefghijklmnopkrstuvwsyz012345"

func Key(key ...string) []byte {
	str := DEFAULT_KEY
	if len(key) >= 1 {
		str = key[0]
	}

	keyLen := len(str)
	if keyLen < 16 {
		// padding this key
		for i := 0; i < 16-keyLen; i++ {
			str += "0"
		}
		keyLen = 16
	}

	arr := []byte(str)
	if keyLen >= 32 {
		return arr[:32]
	} else if keyLen >= 24 {
		return arr[:24]
	} else {
		return arr[:16]
	}

}

func Encrypt(plantText, key []byte) (encrypted []byte, err error) {

	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	block, err := stdaes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	plantText = PKCS7Padding(plantText, block.BlockSize())
	if len(plantText)%block.BlockSize() != 0 {
		return nil, errors.New("Need a multiple of the blocksize")
	}

	encrypted = make([]byte, len(plantText))

	mode := cipher.NewCBCEncrypter(block, key[:stdaes.BlockSize])
	mode.CryptBlocks(encrypted, plantText)

	return encrypted, nil
}

func EncryptToString(plantText, key string) (base64Text string, err error) {
	keyBytes := Key(key)
	ciphertext, err := Encrypt([]byte(plantText), keyBytes)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(ciphertext, key []byte) (decrypted []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	block, err := stdaes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	decrypted = make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, key[:stdaes.BlockSize])
	mode.CryptBlocks(decrypted, ciphertext)

	decrypted = PKCS7UnPadding(decrypted, block.BlockSize())

	return decrypted, nil
}

func DecryptToString(base64Text, key string) (plantText string, err error) {
	keyBytes := Key(key)
	ciphertext, err := base64.StdEncoding.DecodeString(base64Text)
	if err != nil {
		return "", err
	}
	ciphertext, err = Decrypt(ciphertext, keyBytes)
	if err != nil {
		return "", err
	}
	return string(ciphertext), nil
}
