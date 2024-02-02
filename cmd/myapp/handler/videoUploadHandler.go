// Package handler provides various handlers for handling video uploads.
//
// The HandleVideoUpload function is the main function that handles video uploads.
// It reads the video file from the request, saves it to a local directory,
// gets the video's duration and quality, encodes the video to H.265 format,
// and deletes the old file.
//
// The getVideoQuality function determines the quality of the video based on its height.
//
// The getVideoInfo function uses the ffprobe command to get the height and duration of the video.

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
	FileName             string   `json:"fileName"`
	VideoTitle           string   `json:"videoTitle"`
	VideoDescription     string   `json:"videoDescription"`
	VideoStoragePath     string   `json:"videoStoragePath"`
	ThumbnailStoragePath string   `json:"thumbnailStoragePath"`
	UploaderId           string   `json:"uploaderId"`
	VideoDuration        string   `json:"videoDuration"`
	Genres               []string `json:"genres"`
	Types                string   `json:"types"`
}

// VideoDetailsInfo is a global variable of type videoDetails that stores the details of the uploaded video.
var VideoDetailsInfo videoDetails

// HandleVideoUpload handles the video upload process. It reads the video file from the request,
// saves it to a local directory, gets the video's duration and quality, encodes the video to H.265 format,
// and deletes the old file.
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

	// Before encoding the video, get its duration and quality
	videoFilePath := "./userUploadDatas/videos/" + fileName
	duration, quality, err := getVideoInfo(videoFilePath)
	if err != nil {
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
		"message": "File uploaded, converted to H.265, and old file deleted successfully",
		"success": true,
	})
}

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
