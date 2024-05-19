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
