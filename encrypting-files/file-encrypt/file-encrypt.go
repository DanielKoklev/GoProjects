package filecrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"io"
	"os"

	"golang.org/x/crypto/pbkdf2"
)

func Encrypt(source string, password []byte){
	if _, err := os.Stat(source); os.IsNotExist(err) {
		panic(err.Error())
	}
	
	srcFile, err := os.Open(source)
	if err != nil {
		panic(err.Error())
	}

	defer srcFile.Close()

	plainText, err := io.ReadAll(srcFile)
	if err != nil {
		panic(err.Error())
	}

	key := password
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	
	dk := pbkdf2.Key(key, nonce, 4096, 32, sha1.New)

	block, err := aes.NewCipher(dk)
	if err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)

	if err != nil {
		panic(err.Error())
	}

	cipherText := aesgcm.Seal(nil, nonce, plainText, nil)
	cipherText = append(cipherText, nonce...)

	destFile, err := os.Create(source)
	if err != nil {
		panic(err.Error())
	}

	defer destFile.Close().Error()

	_, err := destFile.Write(cipherText)
	
	// not finished
}

func Decrypt(source string, password []byte){

}