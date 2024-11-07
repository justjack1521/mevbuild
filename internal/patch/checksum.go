package patch

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

func GetChecksum(path string) (string, error) {

	handle, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer handle.Close()

	var hash = sha256.New()
	if _, err := io.Copy(hash, handle); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil

}
