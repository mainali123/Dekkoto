package handler

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Declare fileName as a global variable
var GfileName string

type videoDetails struct {
	FileName             string   `json:"fileName"`
	VideoTitle           string   `json:"videoTitle"`
	VideoDescription     string   `json:"videoDescription"`
	VideoStoragePath     string   `json:"videoStoragePath"`
	ThumbnailStoragePath string   `json:"thumbnailStoragePath"`
	UploaderId           string   `json:"uploaderId"`
	VideoDuration        string   `json:"videoDuration"`
	Genres               []string `json:"genres"`
	Types                []string `json:"types"`
}

var VideoDetailsInfo videoDetails

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
	if qualityInt >= 1080 {
		cmd = exec.Command("C:\\ffmpeg-6.1-full_build\\bin\\ffmpeg",
			"-i",
			"./userUploadDatas/videos/"+fileName,
			"-c:v", "hevc_nvenc", // NVENC codec for H.265 encoding
			"-b:v", "6M",
			"-crf", "26", // Lower CRF for better quality
			"-preset", "fast", // Adjust preset according to speed/quality trade-off
			"./userUploadDatas/videos/"+fileName+"_encoded"+".mp4")
	} else if qualityInt > 720 && qualityInt < 1080 {
		cmd = exec.Command("C:\\ffmpeg-6.1-full_build\\bin\\ffmpeg",
			"-i",
			"./userUploadDatas/videos/"+fileName,
			"-c:v", "hevc_nvenc", // NVENC codec for H.265 encoding
			"-b:v", "4M",
			"-crf", "12", // Lower CRF for better quality
			"-preset", "fast", // Adjust preset according to speed/quality trade-off
			"./userUploadDatas/videos/"+fileName+"_encoded"+".mp4")
	} else if qualityInt > 480 && qualityInt <= 720 {
		cmd = exec.Command("C:\\ffmpeg-6.1-full_build\\bin\\ffmpeg",
			"-i",
			"./userUploadDatas/videos/"+fileName,
			"-c:v", "hevc_nvenc", // NVENC codec for H.265 encoding
			"-b:v", "2M",
			"-crf", "8", // Lower CRF for better quality
			"-preset", "fast", // Adjust preset according to speed/quality trade-off
			"./userUploadDatas/videos/"+fileName+"_encoded"+".mp4")
	} else {
		cmd = exec.Command("C:\\ffmpeg-6.1-full_build\\bin\\ffmpeg",
			"-i",
			"./userUploadDatas/videos/"+fileName,
			"-c:v", "hevc_nvenc", // NVENC codec for H.265 encoding
			"-b:v", "0.3M",
			"-crf", "2", // Lower CRF for better quality
			"-preset", "fast", // Adjust preset according to speed/quality trade-off
			"./userUploadDatas/videos/"+fileName+"_encoded"+".mp4")
	}
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

	VideoDetailsInfo.FileName = fileName
	VideoDetailsInfo.VideoStoragePath = "./userUploadDatas/videos/" + fileName + "_encoded" + ".mp4"
	VideoDetailsInfo.VideoDuration = duration

	c.String(200, "File uploaded, converted to AV1, and old file deleted successfully")

}

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
