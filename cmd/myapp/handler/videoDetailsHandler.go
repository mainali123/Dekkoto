package handler

import "github.com/gin-gonic/gin"

func VideoDetails(c *gin.Context) {
	type videoDetails struct {
		VideoName        string `json:"videoName"`
		VideoDescription string `json:"videoDescription"`
	}

	var videoDetailsStruct videoDetails

	err := c.BindJSON(&videoDetailsStruct)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

}
