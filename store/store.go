package store

import (
	"bytes"
	"io"
	"os"
)

const (
	// FriendlyIdLen is the length of any saved posts's filename.
	FriendlyIdLen = 13
)

// SavePost writes the given bytes to a file with a randomly generated name in
// the given directory.
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
