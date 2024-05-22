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
