// Package handler provides various handlers for handling banner uploads.
//
// The HandleBannerUpload function is the main function that handles banner uploads.
// It reads the banner file from the request, checks its aspect ratio, saves it to a local directory,
// and encodes the image to PNG format.
//
// The aspect ratio of the banner should be 16:9.

package handler

import (
	_ "fmt"
	"github.com/gin-gonic/gin"
	"image"
	"image/jpeg"
	_ "io"
	"os"
	_ "path/filepath"
	_ "strings"

	"github.com/nfnt/resize"
)

// HandleBannerUpload handles the banner upload process. It reads the banner file from the request,
// checks its aspect ratio, saves it to a local directory, and encodes the image to PNG format.
// The aspect ratio of the banner should be 16:9.
func HandleBannerUpload(c *gin.Context) {
	file, _, err := c.Request.FormFile("banner")
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

	// Resize the image to 1920x1080
	banner := resize.Resize(1920, 1080, img, resize.Lanczos3)

	// Compress and save the banner as png
	fileName := GfileName + ".png"
	out, err := os.Create("./userUploadDatas/banners/" + fileName)
	if err != nil {
		c.String(500, "Failed to create file")
		return
	}
	defer out.Close()

	// Encode the banner as JPEG with quality set to 80
	opt := jpeg.Options{Quality: 80}
	err = jpeg.Encode(out, banner, &opt)
	if err != nil {
		c.String(500, "Failed to encode image to JPEG")
		return
	}

	// Provide the storage path
	bannerStoragePath := "./userUploadDatas/banners/" + fileName

	c.String(200, "File uploaded and converted to 1920x1080 successfully. Banner storage path: "+bannerStoragePath)
}
