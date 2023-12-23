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

func HandleThumbnailUpload(c *gin.Context) {
	file, header, err := c.Request.FormFile("thumbnail")
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

	fileName := GfileName + ".png"
	out, err := os.Create("./userUploadDatas/thumbnails/" + fileName)
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

	c.String(200, "File uploaded and converted to PNG successfully")
}
