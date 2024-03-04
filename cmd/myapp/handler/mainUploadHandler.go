package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os/exec"
	"strings"
)

func MainUpload(c *gin.Context) {

	// Video
	videoFile, err := c.FormFile("video")
	if err != nil {
		fmt.Println("Error: ", err)
		c.JSON(400, gin.H{
			"message": "Video file not found",
		})
		return
	}

	// Call HandleVideoUpload function with videoFile as a parameter
	message, err := HandleVideoUpload(c, videoFile)

	if err != nil {
		fmt.Println("Error: ", err)
		fmt.Println("Message: ", message)
	}

	// Thumbnail
	thumbnailFile, err := c.FormFile("thumbnail")
	if err != nil {
		fmt.Println("Error: ", err)
		c.JSON(400, gin.H{
			"message": "Thumbnail file not found",
		})
		return
	}

	message, err = HandleThumbnailUpload(c, thumbnailFile)
	if err != nil {
		fmt.Println("Error: ", err)
		fmt.Println("Message: ", message)
	}

	// Banner
	bannerFile, err := c.FormFile("banner")
	if err != nil {
		fmt.Println("Error: ", err)
		c.JSON(400, gin.H{
			"message": "Banner file not found",
		})
		return
	}

	message, err = HandleBannerUpload(c, bannerFile)
	if err != nil {
		fmt.Println("Error: ", err)
		fmt.Println("Message: ", message)
		return
	}

	title := c.PostForm("title")
	description := c.PostForm("description")
	category := c.PostForm("type")
	genres := c.PostFormArray("genre")

	genresCSV := strings.Join(genres, ",")

	message, err = VideoDetails(c, title, description, category, genresCSV)

	if err != nil {
		fmt.Println("Error: ", err)
		fmt.Println("Message: ", message)
		return
	}

	println("Video uploaded successfully")

	notification()

	// Trigger /confirmVideo route using POST method
	c.Request.Method = "POST"
	c.Request.URL.Path = "/confirmVideo"
	c.Writer.WriteHeader(http.StatusFound)
	c.Writer.Header().Set("Location", "/aakash")

}

func notification() {
	command := "F:\\[FYP]\\Dekkoto\\internal\\toast64.exe " +
		"--app-id \"Dekkoto\" " +
		"--title \"Video_Uploaded\" " +
		"--message \"Video_upload_successful\" " +
		"--duration \"long\" "

	//exec.Command
	cmd := exec.Command("cmd", "/C", command)
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error: ", err)
	}
}
