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
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GfileName is a global variable that stores the name of the uploaded file.
var GfileName string

// videoDetails is a struct that holds various details about the uploaded video.
type videoDetails struct {
	FileName             string `json:"fileName"`
	VideoTitle           string `json:"videoTitle"`
	VideoDescription     string `json:"videoDescription"`
	VideoStoragePath     string `json:"videoStoragePath"`
	ThumbnailStoragePath string `json:"thumbnailStoragePath"`
	UploaderId           string `json:"uploaderId"`
	VideoDuration        string `json:"videoDuration"`
	Genres               string `json:"genres"`
	Types                string `json:"types"`
}

// VideoDetailsInfo is a global variable of type videoDetails that stores the details of the uploaded video.
var VideoDetailsInfo videoDetails

// HandleVideoUpload handles the video upload process. It reads the video file from the request,
// saves it to a local directory, gets the video's duration and quality, encodes the video to H.265 format,
// and deletes the old file.
//
// The function expects a "video" file in the request. If the file is not present, it sends a 500 error response with the message "Failed to read file from request".
// If the function fails to create the file, it sends a 500 error response with the message "Failed to create file".
// If the function fails to save the file, it sends a 500 error response with the message "Failed to save file".
// If the function fails to get the video info, it sends a 500 error response with the message "Failed to get video info".
// If the function fails to convert the video, it sends a 500 error response with the message "Failed to convert video".
// If the function fails to stop using the old file, it sends a 500 error response with the message "Failed to stop using the old file".
// If the function fails to delete the old file, it sends a 500 error response with the message "Failed to delete old file".
// If the function successfully uploads, converts, and deletes the old file, it sends a 200 response with the message "File uploaded, converted to H.265, and old file deleted successfully".
//
// The function uses the "os", "io", "os/exec", "strconv", "math/rand", "time", "net/http", "fmt", "github.com/gin-gonic/gin", and "github.com/google/uuid" packages.
// The function saves the video in the "./userUploadDatas/videos/" directory.
// The function provides the storage path of the video in the VideoDetailsInfo.VideoStoragePath variable.
// The function provides the duration of the video in the VideoDetailsInfo.VideoDuration variable.
func HandleVideoUpload(c *gin.Context) {
	file, _, err := c.Request.FormFile("video")
	if err != nil {
		fmt.Println("ERROR 1")
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
		fmt.Println("ERROR 2")
		c.String(500, "Failed to create file")
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		fmt.Println("ERROR 3")
		c.String(500, "Failed to save file")
		return
	}

	// Before encoding the video, get its duration and quality
	videoFilePath := "./userUploadDatas/videos/" + fileName
	duration, quality, err := getVideoInfo(videoFilePath)
	if err != nil {
		fmt.Println("ERROR 4")
		c.String(500, "Failed to get video info")
		return
	}

	// convert quality from string to int
	qualityInt, err := strconv.Atoi(quality)
	// convert duration from string to float64
	durationFloat, err := strconv.ParseFloat(duration, 64)
	fmt.Println("Duration:", durationFloat, "Quality:", qualityInt)

	// Encode the video
	// if the quality of video is 1080p or higher
	var cmd *exec.Cmd

	switch {
	case qualityInt >= 1080:
		cmd = encodeVideo(fileName, "4M", "28", "fast", "1280x720")
	case qualityInt > 720 && qualityInt < 1080:
		cmd = encodeVideo(fileName, "2M", "24", "fast", "854x480")
	case qualityInt > 480 && qualityInt <= 720:
		cmd = encodeVideo(fileName, "1M", "20", "fast", "640x360")
	default:
		cmd = encodeVideo(fileName, "0.5M", "18", "fast", "480x270")
	}

	//// Set hardware acceleration flags if supported
	cmd.Env = append(os.Environ(),
		"CUDA_VISIBLE_DEVICES=0", // Utilize the first GPU device
	) // Optional: create a log for debugging ffmpeg

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
		//c.String(500, "Failed to delete old file")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to delete old file",
			"success": false,
		})
		return
	}

	VideoDetailsInfo.FileName = fileName
	VideoDetailsInfo.VideoStoragePath = "./userUploadDatas/videos/" + fileName + "_encoded" + ".mp4"
	VideoDetailsInfo.VideoDuration = duration

	//c.String(200, "File uploaded, converted to AV1, and old file deleted successfully")
	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded and old file deleted successfully",
		"success": true,
	})
}

// encodeVideo is a helper function that returns a command to encode a video file using ffmpeg.
// It takes the following parameters:
// - fileName: the name of the video file to be encoded.
// - bitrate: the bitrate to be used for the encoding process.
// - crf: the Constant Rate Factor (CRF) to be used for the encoding process.
// - preset: the preset to be used for the encoding process.
// - resolution: the resolution to be used for the encoding process.
//
// The function uses the "os/exec" package.
// The function returns a command that can be run to encode the video file.
//
// The function expects the video file to be located in the "./userUploadDatas/videos/" directory.
// The function saves the encoded video in the same directory with "_encoded" appended to the file name and ".mp4" as the file extension.
func encodeVideo(fileName, bitrate, crf, preset, resolution string) *exec.Cmd {
	return exec.Command("C:\\ffmpeg-6.1-full_build\\bin\\ffmpeg",
		"-i", "./userUploadDatas/videos/"+fileName,
		"-c:v", "libx264",
		"-b:v", bitrate,
		"-crf", crf,
		"-preset", preset,
		"-vf", "scale="+resolution,
		"./userUploadDatas/videos/"+fileName+"_encoded"+".mp4")
}

// getVideoQuality determines the quality of the video based on its height.
func getVideoQuality(height int) string {
	switch {
	case height <= 240:
		return "240"
	case height <= 480:
		return "480"
	case height <= 720:
		return "720"
	case height <= 1080:
		return "1080"
	// Add more cases as needed for higher resolutions
	default:
		return "2160"
	}
}

// getVideoInfo uses the ffprobe command to get the height and duration of the video.
func getVideoInfo(filePath string) (string, string, error) {
	// Get video height
	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=height", "-of", "default=noprint_wrappers=1:nokey=1", filePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", "", err
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) < 1 {
		return "", "", fmt.Errorf("unexpected ffprobe output: %s", output)
	}

	heightStr := strings.TrimSpace(lines[0])
	height, err := strconv.Atoi(heightStr)
	if err != nil {
		return "", "", fmt.Errorf("failed to convert height to integer: %s", err)
	}

	quality := getVideoQuality(height)

	// Get video duration
	cmd = exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", filePath)
	output, err = cmd.CombinedOutput()
	if err != nil {
		return "", "", err
	}

	lines = strings.Split(string(output), "\n")
	if len(lines) < 1 {
		return "", "", fmt.Errorf("unexpected ffprobe output: %s", output)
	}

	durationStr := strings.TrimSpace(lines[0])
	duration, err := strconv.ParseFloat(durationStr, 64)
	if err != nil {
		return "", "", fmt.Errorf("failed to convert duration to float: %s", err)
	}

	return strconv.FormatFloat(duration, 'f', 2, 64), quality, nil
}
