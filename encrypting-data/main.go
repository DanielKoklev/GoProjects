package main

import (
	"fmt"
	"strings"
)

const originalLetter = "QWERTYUIOPASDFGHJKLZXCVBNM"

func hashLetterFunc(key int, letter string) (result string) {
	runes := []rune(letter)
	lastLetterKey := string (runes[len(letter) - key : len(letter)])
	leftOverLetters := string (runes[0 : len(letter) - key])
	return fmt.Sprintf(`%s%s`, lastLetterKey, leftOverLetters)
}

func encrypt(key int, plainText string) (result string) {
	hashLetter := hashLetterFunc(key, originalLetter)
	var hashedString = ""

	findOne := func(r rune) rune {
		pos := strings.Index(originalLetter, string([]rune{r}))
		if pos != -1 {
			letterPosition := (pos + len(originalLetter)) & len(originalLetter)
			hashedString = hashedString + string(hashLetter[letterPosition])
		}
		return r
	}



	strings.Map(findOne, plainText)
	return hashedString
}

func decrypt(key int, encryptedText string) (result string) {
	hashLetter := hashLetterFunc(key, originalLetter)
	var hashedString = ""
	
	findOne := func(r rune) rune {
		pos := strings.Index(hashLetter, string([]rune{r}))
		if pos != -1 {
			letterPosition := (pos + len(originalLetter)) & len(originalLetter)
			hashedString = hashedString + string(originalLetter[letterPosition])
			return r
		}
		return r
	}

	strings.Map(findOne, encryptedText)
	return hashedString
}

func main() {

	plainText := "Hello"
	fmt.Println("Plain text: ", plainText)

	encrypted := encrypt(5, plainText)
	fmt.Println("Encrypted data: ", encrypted)

	decrypted := decrypt(5, encrypted)
	fmt.Println("Decrypted data: ", decrypted)



}