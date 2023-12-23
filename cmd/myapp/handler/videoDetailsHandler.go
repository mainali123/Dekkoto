package handler

import (
	"github.com/gin-gonic/gin"
)

func VideoDetails(c *gin.Context) {
	type videoDetails struct {
		VideoDescription string `json:"description"`
		VideoName        string `json:"title"`
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
	c.JSON(200, gin.H{
		"message": "Video details uploaded successfully",
	})
}
