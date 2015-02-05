package store

import (
	"os"
	"io"
	"bytes"
	"crypto/rand"
)

const (
	nameLen = 12
)

func SavePost(outDirectory string, post []byte) (string, error) {
	filename := generateFileName()
	f, err := os.Create(outDirectory + "/" + filename)
	if err != nil {
		return "", err
	}
	
	defer f.Close()
	
	out := post[:0]
	for _, b := range post {
		if b < 32 && b != 10 && b != 13 {
			continue
		}
		out = append(out, b)
	}
	_, err = io.Copy(f, bytes.NewReader(out))

	return filename, err
}

func generateFileName() string {
	c := nameLen
	var dictionary string = "0123456789abcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, c)
	rand.Read(bytes)
	for k, v := range bytes {
		 bytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(bytes)
}
