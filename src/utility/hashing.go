package utility

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"os"
)

// GetMd5 MD5 Hash with file
func GetMd5(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		log.Fatal(err)
	}

	hashInBytes := hash.Sum(nil)
	return hex.EncodeToString(hashInBytes)
}
