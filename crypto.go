package main

import (
	"io"
	"crypto/rand"
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"
	"fmt"
	"crypto/aes"
	"crypto/cipher"
)


func base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}



func generateHash(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		// TODO: Properly handle error
		fmt.Println(err.Error())
	}
	fmt.Println("The new hash: " + string(hash))
	fmt.Println(password)
	encodedHash := base64.StdEncoding.EncodeToString([]byte(hash))
	return string(encodedHash)
}


func compareHash(password, hash string) bool {


	originalHash, _ := base64.StdEncoding.DecodeString(hash)
	fmt.Println("original hash: " + string(originalHash))
	fmt.Println(password)

	if err := bcrypt.CompareHashAndPassword([]byte(originalHash), []byte(password)); err != nil {
		// TODO: Properly handle error
		fmt.Println(err.Error())
		return false
	}
	return true
}


func decrypt(ciphertextString string) string {
	key, err := base64.StdEncoding.DecodeString(masterKey)
	if err != nil {
		fmt.Println(err.Error())
	}

	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextString)
	if err != nil {
		panic(err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)
	return string(ciphertext)
	// Output: some plaintext
}

func encrypt(text string) string {
	key, err := base64.StdEncoding.DecodeString(masterKey)
	if err != nil {
		fmt.Println(err.Error())
	}

	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	encodedCiphertext := base64.StdEncoding.EncodeToString(ciphertext)

	return encodedCiphertext

	// It's important to remember that ciphertexts must be authenticated
	// (i.e. by using crypto/hmac) as well as being encrypted in order to
	// be secure.
}


// The function called to collect data to do the initial render of the secrets.html page
// It will decrypt all of the user's applications and display the passwords as ******
func decryptUserApplications(userSecrets map[string]string) map[string]string {
	plaintext := make(map[string]string, 0)

	// Decryption happens in memory, never in persistent storage
	for application, _ := range userSecrets {
		plaintextApplication := decrypt(application)

		plaintext[plaintextApplication] = "******"

	}

	return plaintext
}

// Used when a user clicks the magnifying glass
// Assume's we are passing in a map with decrypted applications
// Returns the decrypted secret
func decryptUserSecret(userSecrets map[string]string, requestedApp string) map[string]string {
	plaintext := make(map[string]string, 0)

	// Decryption happens in memory, never in persistent storage
	for application, password := range userSecrets {
		plaintextApplication := decrypt(application)
		if(plaintextApplication == requestedApp) {
			plaintext[plaintextApplication] = decrypt(password)
		} else {
			plaintext[plaintextApplication] = "******"
		}

	}

	return plaintext
}