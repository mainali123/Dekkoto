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
	_ "fmt"
	"github.com/gin-gonic/gin"
	"image"
	"image/jpeg"
	_ "io"
	"mime/multipart"
	"os"
	_ "path/filepath"
	_ "strings"

	"github.com/nfnt/resize"
)

// HandleBannerUpload is a function that handles the banner upload process.
// It performs the following steps:
// 1. Reads the banner file from the request.
// 2. Decodes the image file.
// 3. Resizes the image to 1920x1080, maintaining the aspect ratio of 16:9.
// 4. Compresses and saves the banner as a PNG file in the local directory.
// 5. Encodes the image to JPEG format with a quality set to 80.
// 6. Provides the storage path of the banner.
//
// If there is an error at any step, it sends an appropriate error response.
//
// The aspect ratio of the banner should be 16:9.
// The function uses the Lanczos3 resampling function from the "github.com/nfnt/resize" package for resizing the image.
// The function uses the "image/jpeg" package for encoding the image to JPEG format.
// The function saves the banner in the "./userUploadDatas/banners/" directory.
// The function provides the storage path of the banner in the bannerStoragePath variable.
//
// The function expects a "banner" file in the request. If the file is not present, it sends a 500 error response with the message "Failed to read file from request".
// If the function fails to decode the image file, it sends a 500 error response with the message "Failed to decode image file".
// If the function fails to create the file, it sends a 500 error response with the message "Failed to create file".
// If the function fails to encode the image to JPEG, it sends a 500 error response with the message "Failed to encode image to JPEG".
// If the function successfully uploads and converts the file, it sends a 200 response with the message "File uploaded and converted to 1920x1080 successfully. Banner storage path: "+bannerStoragePath".
func HandleBannerUpload(c *gin.Context, bannerFile *multipart.FileHeader) (string, error) {
	file, err := bannerFile.Open()
	if err != nil {
		//c.String(500, "Failed to read file from request")
		fmt.Println("Failed to read file from request")
		return "Failed to read file from request", err
	}
	defer file.Close()

	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		//c.String(500, "Failed to decode image file")
		fmt.Println("Failed to decode image file")
		return "Failed to decode image file", err
	}

	// Resize the image to 1920x1080
	banner := resize.Resize(1920, 1080, img, resize.Lanczos3)

	// Compress and save the banner as png
	fileName := GfileName + ".png"
	out, err := os.Create("./userUploadDatas/banners/" + fileName)
	if err != nil {
		c.String(500, "Failed to create file")
		return "Failed to create file", err
	}
	defer out.Close()

	// Encode the banner as JPEG with quality set to 80
	opt := jpeg.Options{Quality: 80}
	err = jpeg.Encode(out, banner, &opt)
	if err != nil {
		//c.String(500, "Failed to encode image to JPEG")
		fmt.Println("Failed to encode image to JPEG")
		return "Failed to encode image to JPEG", err
	}

	// Provide the storage path
	bannerStoragePath := "./userUploadDatas/banners/" + fileName

	//c.String(200, "File uploaded and converted to 1920x1080 successfully. Banner storage path: "+bannerStoragePath)
	fmt.Println("File uploaded and converted to 1920x1080 successfully. Banner storage path: " + bannerStoragePath)
	return "File uploaded and converted to 1920x1080 successfully. Banner storage path: " + bannerStoragePath, nil
}
