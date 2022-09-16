// Import the required packages for upload and admin.
package main

import (
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
)

func uploadFiles(c *gin.Context, fullPath string) (*uploader.UploadResult, error) {
	// Add your Cloudinary credentials.
	cld, _ := cloudinary.NewFromParams(CloudinaryBucketName(), CloudinaryApiKey(), CloudinarySecretKey())
	return cld.Upload.Upload(c, fullPath, uploader.UploadParams{PublicID: fullPath})

}
