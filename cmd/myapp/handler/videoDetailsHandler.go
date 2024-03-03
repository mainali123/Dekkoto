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

// Modify VideoDetails to accept video details as parameters
func VideoDetails(c *gin.Context, title string, description string, category string, genresCSV string) (string, error) {
	// if title and description is not empty
	if title == "" || description == "" || category == "" || genresCSV == "" {
		//c.JSON(400, gin.H{
		//	"error": "Video title or description cannot be empty",
		//})
		return "Video title or description cannot be empty", fmt.Errorf("Video title or description cannot be empty")
	}

	VideoDetailsInfo.VideoTitle = title
	VideoDetailsInfo.VideoDescription = description
	VideoDetailsInfo.Genres = genresCSV
	VideoDetailsInfo.Types = category

	fmt.Println(VideoDetailsInfo.VideoTitle, VideoDetailsInfo.VideoDescription, VideoDetailsInfo.Genres, VideoDetailsInfo.Types)
	/*c.JSON(200, gin.H{
		"message": "Video details uploaded successfully",
	})*/
	return "Video details uploaded successfully", nil
}
