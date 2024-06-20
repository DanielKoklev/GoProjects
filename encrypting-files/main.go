package main

import (
	"fmt"
	"os"
	"bytes"
	"golang.org/x/term"
	"github.com/DanielKoklev/Encrypting-Files/file-encrypt"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(0)
	}
	function := os.Args[1]

	switch function {
	case "help":
		printHelp()
	case "encrypt":
		encryptHandle()
	case "decrypt":
		decryptHandle()
	default:
		fmt.Println("Run encrypt to encrypt a file or decrypt to decrypt the file")
		os.Exit(0)
	}
}

func printHelp() {
	fmt.Println("Usage:")
	fmt.Println("")
	fmt.Println("\t go run . encrypt /path/to/file")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("\t encrypt\tEncrypts a file with provided password")
	fmt.Println("\t decrypt\tDecrypts a file with provided password")
	fmt.Println("\t help\tDisplays this message\n")
}

func encryptHandle(){
	if len(os.Args) > 3 {
		fmt.Println("Path to file is not specified")
		os.Exit(0)
	}

	file := os.Args[2]
	if !validateFile(file) {
		panic("File not found!")
	}

	password := getPassword()
	fmt.Println("\nEncrypting...")
	filecrypt.Encrypt(file, password)
	fmt.Println("\nFile successfully encrypted!")

}

func decryptHandle(){
	if len(os.Args) > 3 {
		fmt.Println("Path to file is not specified")
		os.Exit(0)
	}

	file := os.Args[2]
	if !validateFile(file) {
		panic("File not found!")
	}

	fmt.Println("Enter a password for decryption: ")
	password, _ := term.ReadPassword(0)
	fmt.Println("\nDecrypting...")
	filecrypt.Decrypt(file, password)
	fmt.Println("File successfully decrypted!")

}

func getPassword() []byte {
	fmt.Println("Enter password: ")
	password, _ := term.ReadPassword(0)
	fmt.Println("\nConfirm the password: ")
	confirmedPassword, _ := term.ReadPassword(0)
	if !validatePassword(password, confirmedPassword){
		fmt.Println("\nPasswords need to match.\nPlease try again.")
		return getPassword()
	}
	return password
}

func validatePassword(password []byte, confirmedPassword []byte) bool {
	if !bytes.Equal(password, confirmedPassword) {
		return false
	}
	return true
}

func validateFile(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err){
		return false
	}
	return true
}