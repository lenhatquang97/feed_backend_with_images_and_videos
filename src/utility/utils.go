package utility

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/google/uuid"
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

// Download random image from source https://source.unsplash.com/random
func DownloadRandomImageIntoFolderId(id string) string {
	// Get the data
	resp, err := http.Get("https://source.unsplash.com/random")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Create the file
	folderPath := "./files/" + id
	randomIdJpg := uuid.New().String() + ".jpg"

	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		os.Mkdir(folderPath, 0755)
	}

	out, err := os.Create(folderPath + "/" + randomIdJpg)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(2 * time.Second)

	return "https://feeduiclone.win/static/files/" + id + "/" + randomIdJpg
}

func GenerateBatchImages(num int, id string) []string {
	var urls []string
	for i := 0; i < num; i++ {
		result := DownloadRandomImageIntoFolderId(id)
		urls = append(urls, result)
	}
	return urls
}
