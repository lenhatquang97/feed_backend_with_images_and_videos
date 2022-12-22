package utility

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"os/exec"
	"strings"
)

func ExecuteGetThumbnail(file *multipart.FileHeader, id string) {
	thumbnailName := strings.Replace(file.Filename, ".mp4", ".png", 1)
	cmd := exec.Command("ffmpeg", "-y", "-i", "./files/"+id+"/"+file.Filename, "-ss", "00:00:01.000", "-vframes", "1", "./files/"+id+"/"+thumbnailName)
	fmt.Println(cmd.String())
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

// SHA-256 file hash with path
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
