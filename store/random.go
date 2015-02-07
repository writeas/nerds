package store

import (
	"crypto/rand"
)

func Generate62RandomString(l int) string {
	return GenerateRandomString("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz", l)
}

func GenerateFriendlyRandomString(l int) string {
	return GenerateRandomString("0123456789abcdefghijklmnopqrstuvwxyz", l)
}

func GenerateRandomString(dictionary string, l int) string {
	var bytes = make([]byte, l)
	rand.Read(bytes)
	for k, v := range bytes {
		 bytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(bytes)
}
