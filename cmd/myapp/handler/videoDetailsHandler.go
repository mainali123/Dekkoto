// Package handler provides various handlers for handling user requests related to video management in the application.
//
// The package includes handlers for the following operations:
// - Video Upload: Handles the video upload process, including reading the video file from the request, saving it to a local directory, getting the video's duration and quality, encoding the video to H.265 format, and deleting the old file.
// - Video Details: Handles the video details process, including reading the video details from the request, checking if the title and description are not empty, and saving the details to a global variable.
// - Thumbnail Upload: Handles the thumbnail upload process, including reading the thumbnail file from the request, checking its aspect ratio, saving it to a local directory, and encoding the image to PNG format.
// - Banner Upload: Handles the banner upload process, including reading the banner file from the request, checking its aspect ratio, saving it to a local directory, and encoding the image to PNG format.
//
// The package uses the "os", "io", "os/exec", "strconv", "math/rand", "time", "net/http", "fmt", "github.com/gin-gonic/gin", and "github.com/google/uuid" packages.
// The package saves the video, thumbnail, and banner in the "./userUploadDatas/videos/", "./userUploadDatas/thumbnails/", and "./userUploadDatas/banners/" directories respectively.
// The package provides the storage path of the video, thumbnail, and banner in the VideoDetailsInfo.VideoStoragePath, VideoDetailsInfo.ThumbnailStoragePath, and bannerStoragePath variables respectively.
// The package provides the duration of the video in the VideoDetailsInfo.VideoDuration variable.
package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// VideoDetails handles the video details process. It reads the video details from the request,
// checks if the title and description are not empty, and saves the details to a global variable.
func VideoDetails(c *gin.Context) {
	type videoDetails struct {
		VideoDescription string   `json:"description"`
		VideoName        string   `json:"title"`
		Genres           []string `json:"genres"`
		Types            string   `json:"types"`
	}

	var videoDetailsStruct videoDetails

	err := c.BindJSON(&videoDetailsStruct)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	// if title and description is not empty
	if videoDetailsStruct.VideoName == "" || videoDetailsStruct.VideoDescription == "" {
		c.JSON(400, gin.H{
			"error": "Video title or description cannot be empty",
		})
		return
	}

	VideoDetailsInfo.VideoTitle = videoDetailsStruct.VideoName
	VideoDetailsInfo.VideoDescription = videoDetailsStruct.VideoDescription
	VideoDetailsInfo.Genres = videoDetailsStruct.Genres
	VideoDetailsInfo.Types = videoDetailsStruct.Types

	fmt.Println(VideoDetailsInfo.VideoTitle, VideoDetailsInfo.VideoDescription, VideoDetailsInfo.Genres, VideoDetailsInfo.Types)
	c.JSON(200, gin.H{
		"message": "Video details uploaded successfully",
	})
}
