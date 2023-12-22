package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"time"
)

// Declare fileName as a global variable
var GfileName string

func HandleVideoUpload(c *gin.Context) {
	file, _, err := c.Request.FormFile("video")
	if err != nil {
		c.String(500, "Failed to read file from request")
		return
	}
	defer file.Close()

	currentTime := time.Now().Format("20060102150405.000000")
	uniqueID := uuid.New()
	randomNumber := strconv.Itoa(rand.Intn(999999999999))

	GfileName = currentTime + "_" + uniqueID.String() + "_" + randomNumber
	fileName := GfileName

	out, err := os.Create("./userUploadDatas/videos/" + fileName)
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

	// encoding the video
	//cmd := exec.Command("C:\\ffmpeg-6.1-full_build\\bin\\ffmpeg",
	//	"-i",
	//	"./userUploadDatas/videos/"+fileName,
	//	"-c:v", "libvpx-vp9",
	//	//"-preset", "ultrafast", // Use a faster preset
	//	//"-hwaccel", "auto", // Use hardware acceleration if available
	//	"-b:v", "45M", // Increase bitrate for better quality
	//	"-crf", "0", // Lower CRF for better quality
	//	"./userUploadDatas/videos/"+fileName+"_encoded"+".webm")

	cmd := exec.Command("C:\\ffmpeg-6.1-full_build\\bin\\ffmpeg",
		"-i",
		"./userUploadDatas/videos/"+fileName,
		"-c:v", "hevc_nvenc", // NVENC codec for H.265 encoding
		"-b:v", "6M",
		"-crf", "26", // Lower CRF for better quality
		"-preset", "fast", // Adjust preset according to speed/quality trade-off
		"./userUploadDatas/videos/"+fileName+"_encoded"+".mp4")

	//// Set hardware acceleration flags if supported
	cmd.Env = append(os.Environ(),
		"CUDA_VISIBLE_DEVICES=0",                  // Utilize the first GPU device
		"FFREPORT=file=./ffmpeg_log.log:level=48") // Optional: create a log for debugging ffmpeg

	err = cmd.Run()
	if err != nil {
		fmt.Println(err) // Print out the error
		c.String(500, "Failed to convert video")
		return
	}

	// stop using the old file
	err = out.Close()
	if err != nil {
		fmt.Println(err) // Print out the error
		c.String(500, "Failed to stop using the old file")
		return
	}

	// Delete the old file
	err = os.Remove("./userUploadDatas/videos/" + fileName)
	if err != nil {
		fmt.Println("Failed to delete old file:", err) // Print out the error
		c.String(500, "Failed to delete old file")
		return
	}

	c.String(200, "File uploaded, converted to AV1, and old file deleted successfully")
}

func HandleThumbnailUpload(c *gin.Context) {
	file, _, err := c.Request.FormFile("thumbnail")
	if err != nil {
		c.String(500, "Failed to read file from request")
		return
	}
	defer file.Close()

	fileName := GfileName

	out, err := os.Create("./userUploadDatas/thumbnails/" + fileName)
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
