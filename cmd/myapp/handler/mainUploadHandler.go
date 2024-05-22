// Package handler provides various handlers for handling user requests in the application.
// These handlers are responsible for processing incoming HTTP requests related to video management,
// such as uploading videos, thumbnails, and banners, as well as handling video details.
// The package includes functions for reading files from requests, handling file uploads,
// encoding videos, resizing images, and interacting with the database.
// It also includes functions for sending notifications to the user and redirecting the user to different pages.
// This package is essential for the functionality of the video management system in the application.
package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os/exec"
	"strings"
)

// MainUpload is the main function that handles the video upload process.
// It performs the following steps:
// 1. Reads the video, thumbnail, and banner files from the request.
// 2. Calls the HandleVideoUpload, HandleThumbnailUpload, and HandleBannerUpload functions to handle the upload of the video, thumbnail, and banner respectively.
// 3. Reads the video details from the request.
// 4. Calls the VideoDetails function to handle the video details.
// 5. Calls the uploadOnDatabase function to upload the video to the database.
// 6. Redirects the user to the admin panel.
//
// If there is an error at any step, it sends an appropriate error response.
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

	uploadOnDatabase()
	c.Redirect(http.StatusFound, "/adminPanel")
}

// notification is a function that sends a notification to the user indicating that the video upload was successful.
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

// uploadOnDatabase is a function that uploads the video to the database.
// It performs the following steps:
// 1. Calls a curl command to upload the video to the database.
// 2. If there is an error, it sends a notification to the user indicating that the video upload failed.
// 3. If the video is uploaded successfully, it calls the notification function to send a notification to the user indicating that the video upload was successful.
func uploadOnDatabase() {
	command := "curl -X POST http://localhost:8080/uploadVideoInDatabase"
	cmd := exec.Command("cmd", "/C", command)
	err := cmd.Run()
	if err != nil {
		command := "F:\\[FYP]\\Dekkoto\\internal\\toast64.exe " +
			"--app-id \"Dekkoto\" " +
			"--title \"Video_Uploaded\" " +
			"--message \"Unable_to_upload_video\" " +
			"--duration \"long\" "

		//exec.Command
		cmd := exec.Command("cmd", "/C", command)
		err := cmd.Run()
		if err != nil {
			fmt.Println("Error: ", err)
		}
		return
	}
	notification()
}
