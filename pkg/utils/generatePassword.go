package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

var (
	lowerChars   = "abcdefghijklmnopqrstuvwxyz"
	upperChars   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digitChars   = "0123456789"
	specialChars = "!@#$%^&*()-_=+[]{}<>?/|"
	allChars     = lowerChars + upperChars + digitChars + specialChars
)

func getRandomChar(charset string) (byte, error) {
	index, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
	if err != nil {
		return 0, err
	}
	return charset[index.Int64()], nil
}

func shufflePassword(password []byte) ([]byte, error) {
	for i := len(password) - 1; i > 0; i-- {
		jBig, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			return nil, err
		}
		j := int(jBig.Int64())
		password[i], password[j] = password[j], password[i]
	}
	return password, nil
}

func GeneratePassword(length int) (string, error) {
	if length < 3 {
		return "", fmt.Errorf("la longitud mínima debe ser al menos 3")
	}

	password := make([]byte, 0, length)

	requiredSets := []string{upperChars, digitChars, specialChars}
	for _, set := range requiredSets {
		char, err := getRandomChar(set)
		if err != nil {
			return "", err
		}
		password = append(password, char)
	}

	for len(password) < length {
		char, err := getRandomChar(allChars)
		if err != nil {
			return "", err
		}
		password = append(password, char)
	}

	shuffled, err := shufflePassword(password)
	if err != nil {
		return "", err
	}

	return string(shuffled), nil
}
