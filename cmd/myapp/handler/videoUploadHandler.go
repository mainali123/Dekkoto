package handler

import (
	"github.com/gin-gonic/gin"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func HandleVideoUpload(c *gin.Context) {
	file, header, err := c.Request.FormFile("video")
	if err != nil {
		c.String(500, "Failed to read file from request")
		return
	}
	defer file.Close()

	// filename = currentTime + Random 5 digit from 0-9 + originalFilename + Random 5 digit from 0-9 + extension
	fileName := time.Now().Format("20060102150405") + strconv.Itoa(rand.Intn(99999)) + header.Filename + strconv.Itoa(rand.Intn(99999)) + filepath.Ext(header.Filename)

	out, err := os.Create("./videos/" + fileName)
	if err != nil {
		c.String(500, "Failed to create file")
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		c.String(500, "Failed to save file")
		return
	}

	c.String(200, "File uploaded successfully")
}
