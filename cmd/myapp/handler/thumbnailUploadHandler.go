// Package handler provides various handlers for handling thumbnail uploads.
//
// The HandleThumbnailUpload function is the main function that handles thumbnail uploads.
// It reads the thumbnail file from the request, checks its aspect ratio, saves it to a local directory,
// and encodes the image to PNG format.
//
// The aspect ratio of the thumbnail should be 16:9.

package handler

import (
	_ "fmt"
	"github.com/gin-gonic/gin"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	_ "io"
	"os"
	_ "path/filepath"
	_ "strings"

	"github.com/nfnt/resize"
)

// HandleThumbnailUpload handles the thumbnail upload process. It reads the thumbnail file from the request,
// checks its aspect ratio, saves it to a local directory, and encodes the image to PNG format.
// The aspect ratio of the thumbnail should be 16:9.
func HandleThumbnailUpload(c *gin.Context) {
	file, _, err := c.Request.FormFile("thumbnail")
	if err != nil {
		c.String(500, "Failed to read file from request")
		return
	}
	defer file.Close()

	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		c.String(500, "Failed to decode image file")
		return
	}

	// Resize the image to 1080x1920
	thumbnail := resize.Resize(1080, 1920, img, resize.Lanczos3)

	// Compress and save the thumbnail as png
	fileName := GfileName + ".png"
	out, err := os.Create("./userUploadDatas/thumbnails/" + fileName)
	if err != nil {
		c.String(500, "Failed to create file")
		return
	}
	defer out.Close()

	// Encode the thumbnail as JPEG with quality set to 80
	opt := jpeg.Options{Quality: 80}
	err = jpeg.Encode(out, thumbnail, &opt)
	if err != nil {
		c.String(500, "Failed to encode image to JPEG")
		return
	}

	// Provide the storage path
	VideoDetailsInfo.ThumbnailStoragePath = "./userUploadDatas/thumbnails/" + fileName

	c.String(200, "File uploaded and converted to 1080x1920 successfully")
}
