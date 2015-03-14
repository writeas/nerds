package store

import (
	"bytes"
	"io"
	"os"
)

const (
	FriendlyIdLen = 13
)

func SavePost(outDir string, post []byte) (string, error) {
	filename := generateFileName()
	f, err := os.Create(outDir + "/" + filename)
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
	return GenerateFriendlyRandomString(FriendlyIdLen)
}
