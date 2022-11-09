// Import the required packages for upload and admin.
package cdn

import (
	"example.com/feed_backend/src/configs"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
)

func UploadFiles(c *gin.Context, fullPath string) (*uploader.UploadResult, error) {
	// Add your Cloudinary credentials.
	cld, _ := cloudinary.NewFromParams(configs.CloudinaryBucketName(), configs.CloudinaryApiKey(), configs.CloudinarySecretKey())
	return cld.Upload.Upload(c, fullPath, uploader.UploadParams{PublicID: fullPath})

}
