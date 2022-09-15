package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const (
	AWS_S3_REGION         = "ap-southeast-1"
	AWS_S3_BUCKET         = "customfeedbucket"
	AWS_ACCESS_KEY_ID     = "AKIAV3XRZF7WLEZQMNGD"
	AWS_SECRET_ACCESS_KEY = "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"
)

func uploadFile(basePath string, filePath string) error {
	sessionObj, err := session.NewSession(&aws.Config{Region: aws.String(AWS_S3_REGION),

		Credentials: credentials.NewStaticCredentials(
			AWS_ACCESS_KEY_ID,
			AWS_SECRET_ACCESS_KEY,
			"",
		),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Get the fileName from Path
	fileName := filepath.Base(filePath)

	// Open the file from the file path
	upFile, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("could not open local filepath [%v]: %+v", filePath, err)
	}
	defer upFile.Close()

	uploader := s3manager.NewUploader(sessionObj)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket:             aws.String(AWS_S3_BUCKET),
		ACL:                aws.String("public-read"),
		Key:                aws.String(basePath + fileName),
		Body:               upFile,
		ContentDisposition: aws.String("attachment"),
	})
	return err
}
