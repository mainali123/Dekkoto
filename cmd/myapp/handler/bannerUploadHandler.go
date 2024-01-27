// Package handler provides various handlers for handling banner uploads.
//
// The HandleBannerUpload function is the main function that handles banner uploads.
// It reads the banner file from the request, checks its aspect ratio, saves it to a local directory,
// and encodes the image to PNG format.
//
// The aspect ratio of the banner should be 16:9.

package handler

import (
	"github.com/gin-gonic/gin"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

// HandleBannerUpload handles the banner upload process. It reads the banner file from the request,
// checks its aspect ratio, saves it to a local directory, and encodes the image to PNG format.
// The aspect ratio of the banner should be 16:9.
func HandleBannerUpload(c *gin.Context) {
	file, header, err := c.Request.FormFile("banner")
	if err != nil {
		c.String(500, "Failed to read file from request")
		return
	}
	defer file.Close()

	var img image.Image
	switch strings.ToLower(filepath.Ext(header.Filename)) {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(file)
		if err != nil {
			c.String(500, "Failed to decode JPEG file")
			return
		}
	case ".png":
		img, err = png.Decode(file)
		if err != nil {
			c.String(500, "Failed to decode PNG file")
			return
		}
	default:
		c.String(500, "Unsupported file format")
		return
	}

	// Check if the image has the correct aspect ratio for a thumbnail (16:9)
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	if width*9 != height*16 {
		c.String(400, "Image does not have the correct aspect ratio for a thumbnail (16:9)")
		return
	}

	fileName := GfileName + ".png"
	out, err := os.Create("./userUploadDatas/banners/" + fileName)
	if err != nil {
		c.String(500, "Failed to create file")
		return
	}
	defer out.Close()

	err = png.Encode(out, img)
	if err != nil {
		c.String(500, "Failed to encode image to PNG")
		return
	}
	VideoDetailsInfo.ThumbnailStoragePath = "./userUploadDatas/thumbnails/" + fileName

	c.String(200, "File uploaded and converted to PNG successfully")
}
