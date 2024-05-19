package main

import (
	"Dekkoto/cmd/myapp/handler"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nfnt/resize"
	"github.com/resend/resend-go/v2"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// login is a handler function that serves the login page.
// It parses the login.html template and executes it, sending the output to the client.
// If there is an error during parsing or execution of the template, it sends a server error response.
func (app *application) login(c *gin.Context) {
	// If logged-In
	if userInfo.UserId != 0 && userInfo.Email != "" {
		app.homePage(c)
		return
	}

	t, err := template.ParseFiles("ui/html/login.html")
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}

	app.deviceInfo(c.Request)
	err = t.Execute(c.Writer, nil)
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}
}

// register is a handler function that serves the registration page.
// It parses the register.html template and executes it, sending the output to the client.
// If there is an error during parsing or execution of the template, it sends a server error response.
func (app *application) register(c *gin.Context) {
	// If logged-In
	if userInfo.UserId != 0 && userInfo.Email != "" {
		app.homePage(c)
		return
	}

	t, err := template.ParseFiles("ui/html/register.html")
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}

	err = t.Execute(c.Writer, nil)
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}

}

// admin is a handler function that serves the admin page.
// It parses the admin.html template and executes it, sending the output to the client.
// If there is an error during parsing or execution of the template, it sends a server error response.
func (app *application) admin(c *gin.Context) {
	t, err := template.ParseFiles("ui/html/admin.html")
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}

	err = t.Execute(c.Writer, nil)
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}
}

// registerPostRequest is a handler function that handles the registration of a new user.
// It reads the user data from the request, validates it, and registers the user in the database.
// If there is an error during any of these steps, it sends an appropriate error response.
func (app *application) registerPostRequest(c *gin.Context) {
	// Response struct
	type User struct {
		Name            string `json:"name"`
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var userData User
	// Bind the JSON data from the request to the userData struct
	if err := c.BindJSON(&userData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	encryptedPass, err := encrypt(userData.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to encrypt password",
		})
		return

	}

	err = app.database.registerUser(userData.Name, userData.Email, encryptedPass)

	if err != nil {
		if err.Error() == "user already exists" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "User already exists",
			})
			return
		}
		// For other errors during registration, return a generic error response
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to register user",
		})
		return
	}

	userID, err := app.database.userId(userData.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user id",
		})
		return
	}

	err = app.database.saveUserProfile("", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save user profile",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User created successfully",
	})

}

// loginPostRequest is a handler function that handles the login of a user.
// It reads the user data from the request, validates it, and logs the user in.
// If there is an error during any of these steps, it sends an appropriate error response.
func (app *application) loginPostRequest(c *gin.Context) {

	// Response struct
	type User struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var userData User

	// Bind the JSON data from the request to the userData struct
	if err := c.BindJSON(&userData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	fmt.Println(encrypt(userData.Password))

	err := app.database.loginUser(userData.Email, userData.Password)

	if err != nil {
		if err.Error() == "user does not exist" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "User does not exist",
			})
			return
		}
		// For other errors during login, return a generic error response
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User logged in successfully",
	})

	// Get user id from the database
	userID, err := app.database.userId(userData.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user id",
		})
		return
	}

	userInfo.UserId = userID
	userInfo.Email = userData.Email
	handler.VideoDetailsInfo.UploaderId = strconv.Itoa(userID)
}

// uploadVideo is a handler function that handles the upload of a video.
// It reads the video data from the request, validates it, and uploads the video to the database.
// If there is an error during any of these steps, it sends an appropriate error response.
func (app *application) uploadVideo(c *gin.Context) {
	videoInfo := handler.VideoDetailsInfo
	// get current date
	currentDate := time.Now().Format("2006-01-02")

	// Convert arrays to comma-separated strings
	categoryString := videoInfo.Genres

	// Map genre and category strings to their respective IDs
	var categoryID int
	// remove all the whitespaces
	videoInfo.Types = strings.ReplaceAll(videoInfo.Types, " ", "")
	categoryID, err := app.database.getCategoryID(categoryString)

	// check if category is empty
	if err != nil {
		//c.String(500, "Failed to get category ID with error: "+err.Error())
		fmt.Println("Error getting category ID with error: " + err.Error())
		return
	}

	//genreID, err := app.database.getGenreID(videoInfo.Genres[0])
	// there are multiple genres do it for all
	// remove all the whitespaces
	categoryString = strings.ReplaceAll(categoryString, " ", "")
	genreID, err := app.database.getGenreID(videoInfo.Types)
	if err != nil {
		//c.String(500, "Failed to get genre ID with error: "+err.Error())
		fmt.Println("Error getting genre ID with error: " + err.Error())
		return
	}

	// print all the data
	fmt.Println(videoInfo.VideoTitle, videoInfo.VideoDescription, videoInfo.VideoStoragePath, videoInfo.ThumbnailStoragePath, videoInfo.UploaderId, currentDate, videoInfo.VideoDuration, genreID, categoryID)

	err = app.database.uploadVideo(
		videoInfo.VideoTitle,
		videoInfo.VideoDescription,
		videoInfo.VideoStoragePath,
		videoInfo.ThumbnailStoragePath,
		videoInfo.UploaderId,
		currentDate,
		videoInfo.VideoDuration,
		categoryID, // Use mapped category ID
		genreID,    // Use mapped genre ID
	)

	if err != nil {
		fmt.Println("Failed to upload video with error: " + err.Error())
		return
	}
	return
}

// terminateVideo is a handler function that handles the termination of a video upload.
// It deletes the video file and the thumbnail file from the server.
// If there is an error during any of these steps, it sends an appropriate error response.
func (app *application) terminateVideo(c *gin.Context) {
	// Delete the video file
	err := os.Remove(handler.VideoDetailsInfo.VideoStoragePath)
	if err != nil {
		c.String(500, "Failed to delete video file")
	}

	// Delete the thumbnail file
	err = os.Remove(handler.VideoDetailsInfo.ThumbnailStoragePath)
	if err != nil {
		c.String(500, "Failed to delete thumbnail file")
	}

	c.String(200, "Terminated successfully")
}

// Data is a global variable of type map with string keys and values of any type (interface{}).
// It is used to hold various data, such as video information, that needs to be accessed across different functions.
var Data map[string]interface{}

// showVideos is a handler function that serves the videos page.
// It fetches the videos data from the database and sends it to the client.
// If there is an error during fetching the videos data, it sends a server error response.
func (app *application) showVideos(c *gin.Context) {
	// If not logged-In
	if userInfo.UserId == 0 && userInfo.Email == "" {
		app.error404(c)
		return
	}

	// Check if user have access to user access page
	_, _, editDelete, _, _, _, err := app.database.userAccess(userInfo.UserId)
	if err != nil {
		app.error403(c)
		return
	}
	if editDelete == 0 {
		app.error403(c)
		return
	}

	t, err := template.ParseFiles("ui/html/admin/adminTables.html")
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}

	err = t.Execute(c.Writer, nil)
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}

	// Get the userID from the context or session
	userID := userInfo.UserId

	// Call the videoDescForTable function with the userID
	videos, err := app.database.videoDescForTable(userID)
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}

	// Create a data map to hold the videos data
	Data = map[string]interface{}{
		"Videos": videos,
	}
}

// showVideosPost is a handler function that handles the post request of the videos page.
// It sends the videos data to the client as a JSON response.
func (app *application) showVideosPost(c *gin.Context) {
	videoList, err := app.database.allVideoList()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "Failed to fetch video details",
			"success": false,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully Fetched video details",
		"success": true,
		"videos":  videoList,
	})
}

// editVideo is a handler function that serves the video editing page.
// It parses the adminEditVideo.html template and executes it, sending the output to the client.
// If there is an error during parsing or execution of the template, it sends a server error response.
func (app *application) editVideo(c *gin.Context) {
	// If not logged-In
	if userInfo.UserId == 0 && userInfo.Email == "" {
		app.error404(c)
		return
	}

	t, err := template.ParseFiles("ui/html/adminEditVideo.html")
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}

	err = t.Execute(c.Writer, nil)
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}
}

// showCategoriesName is a handler function that handles the fetching of the category name.
// It reads the category ID from the request, fetches the category name from the database, and sends it to the client.
// If there is an error during any of these steps, it sends an appropriate error response.
func (app *application) showCategoriesName(c *gin.Context) {
	type category struct {
		CategoryID int `json:"categoryID"`
	}
	var categoryData category
	if err := c.BindJSON(&categoryData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	categoryName, err := app.database.getCategoryName(categoryData.CategoryID)
	if err != nil {
		fmt.Println("Error getting category name")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Category data fetched successfully",
		"success":      true,
		"categoryName": categoryName,
	})
}

// showGenresName is a handler function that handles the fetching of the genre name.
// It reads the genre ID from the request, fetches the genre name from the database, and sends it to the client.
// If there is an error during any of these steps, it sends an appropriate error response.
func (app *application) showGenresName(c *gin.Context) {
	type genre struct {
		GenreID int `json:"genreID"`
	}
	var genreData genre
	if err := c.BindJSON(&genreData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	genreName, err := app.database.getGenreName(genreData.GenreID)
	if err != nil {
		fmt.Println("Error getting genre name")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Genre data fetched successfully",
		"success":   true,
		"genreName": genreName,
	})
}

// editVideoPost is a handler function that handles the post request of the video editing page.
// It reads the video data from the request, validates it, and sends it to the client.
// If there is an error during any of these steps, it sends an appropriate error response.
func (app *application) editVideoPost(c *gin.Context) {
	fmt.Println("edit video post")

	type Video struct {
		VideoID     int    `json:"videoID"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Genre       string `json:"genre"`
		Category    string `json:"category"`
		Thumbnail   string `json:"thumbnail"`
		Banner      string `json:"banner"`
		FileName    string `json:"fileName"`
	}

	var videoData Video

	// Bind the JSON data from the request to the userData struct
	if err := c.ShouldBindJSON(&videoData); err != nil {
		fmt.Println("Error binding JSON data: ", err) // Print the error
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}
	//fmt.Println(videoData)

	var compressImage = func(imagePath string, imageName string, width uint, height uint) error {
		// Open the image file
		imgFile, err := os.Open(imagePath)
		if err != nil {
			return fmt.Errorf("failed to open image file: %v", err)
		}
		defer imgFile.Close()

		// Decode the image
		img, _, err := image.Decode(imgFile)
		if err != nil {
			return fmt.Errorf("failed to decode image: %v", err)
		}

		// Resize the image to 1080x1920
		thumbnail := resize.Resize(width, height, img, resize.Lanczos3)
		//thumbnail := resize.Resize(1080, 1920, img, resize.Lanczos3)

		// Compress and save the thumbnail as png
		//out, err := os.Create("./userUploadDatas/thumbnails/" + imageName)
		out, err := os.Create(imagePath)
		if err != nil {
			return fmt.Errorf("failed to create file: %v", err)
		}
		defer out.Close()

		// Encode the thumbnail as PNG
		err = jpeg.Encode(out, thumbnail, nil)
		if err != nil {
			return fmt.Errorf("failed to encode image to PNG: %v", err)
		}
		return nil
	}

	if videoData.Thumbnail != "Same Image" {
		// Split the data URI to get only the base64-encoded part
		thumbnailDataParts := strings.Split(videoData.Thumbnail, ",")
		if len(thumbnailDataParts) != 2 {
			fmt.Println("Invalid data URI format")
			return
		}

		// Decode the base64 data
		thumbnailBytes, err := base64.StdEncoding.DecodeString(thumbnailDataParts[1])
		if err != nil {
			fmt.Println("Error decoding thumbnail data:", err)
			return
		}

		// Save the decoded image to a file
		err = ioutil.WriteFile("./userUploadDatas/thumbnails/"+videoData.FileName, thumbnailBytes, 0644)
		if err != nil {
			fmt.Println("Error saving image:", err)
			return
		}

		fmt.Println("Thumbnail image saved successfully!")
		err = compressImage("./userUploadDatas/thumbnails/"+videoData.FileName, videoData.FileName, 1080, 1920)
		if err != nil {
			fmt.Println("Error compressing image:", err)
			return
		}
		fmt.Println("Thumbnail image compressed successfully!")
	}

	if videoData.Banner != "Same Image" {
		// Split the data URI to get only the base64-encoded part
		bannerDataParts := strings.Split(videoData.Banner, ",")
		if len(bannerDataParts) != 2 {
			fmt.Println("Invalid data URI format")
			return
		}

		// Decode the base64 data
		bannerBytes, err := base64.StdEncoding.DecodeString(bannerDataParts[1])
		if err != nil {
			fmt.Println("Error decoding thumbnail data:", err)
			return
		}

		// Save the decoded image to a file
		err = ioutil.WriteFile("./userUploadDatas/banners/"+videoData.FileName, bannerBytes, 0644)
		if err != nil {
			fmt.Println("Error saving image:", err)
			return
		}

		fmt.Println("Banner image saved successfully!")
		err = compressImage("./userUploadDatas/banners/"+videoData.FileName, videoData.FileName, 1920, 1080)
		if err != nil {
			fmt.Println("Error compressing image:", err)
			return
		}
		fmt.Println("Banner image compressed successfully!")
	}

	err := app.database.updateVideo(videoData.Title, videoData.Description, videoData.Genre, videoData.Category, videoData.VideoID)
	if err != nil {
		fmt.Println("Error updating video details:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update video details",
			"success": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Video details updated successfully",
		"success": true,
	})
}

// updateVideoDetails is a handler function that handles the updating of video details.
// It reads the updated video data from the request, validates it, and updates the video details in the database.
// If there is an error during any of these steps, it sends an appropriate error response.
func (app *application) updateVideoDetails(c *gin.Context) {
	type Video struct {
		VideoID     string `json:"videoID"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Genre       string `json:"genre"`
		Type        string `json:"type"`
	}

	var videoData Video

	// Bind the JSON data from the request to the userData struct
	if err := c.ShouldBindJSON(&videoData); err != nil {
		fmt.Println("Error binding JSON data")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}
	fmt.Println(videoData)

	// get the genre id
	categoryID, err := app.database.getCategoryID(videoData.Genre)
	if err != nil {
		fmt.Println("Error getting genre id" + err.Error())
		return
	}

	// get the category id
	genreID, err := app.database.getGenreID(videoData.Type)
	if err != nil {
		fmt.Println("Error getting category id" + err.Error())
		return
	}

	// convert video id to int
	videoIDInt, err := strconv.Atoi(videoData.VideoID)
	if err != nil {
		fmt.Println("Error converting videoID to int" + err.Error())
	}

	// update the video details
	err = app.database.videoDescForEdit(videoIDInt, videoData.Title, videoData.Description, categoryID, genreID)
	if err != nil {
		fmt.Println("Error updating video details" + err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Video details updated successfully",
		"success": true,
	})
}

// deleteVideo is a handler function that handles the deletion of a video.
// It reads the video ID from the request, deletes the video file and thumbnail file from the server, and deletes the video from the database.
// If there is an error during any of these steps, it sends an appropriate error response.
func (app *application) deleteVideo(c *gin.Context) {
	type Video struct {
		VideoID int `json:"videoID"`
	}

	var videoData Video

	if err := c.ShouldBindJSON(&videoData); err != nil {
		fmt.Println("Error binding JSON data:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	// convert video id to int
	/*videoIDInt, err := strconv.Atoi(videoData.VideoID)
	if err != nil {
		fmt.Println("Error converting videoID to int:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Video ID",
		})
		return
	}*/

	videoName, thumbnailName, err := app.database.deleteVideoFromFile(videoData.VideoID)
	if err != nil {
		fmt.Println("Error getting video name:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get video name",
		})
		return
	}

	deleteVideo := os.Remove(videoName)
	if deleteVideo != nil {
		fmt.Println("Error deleting video:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete video",
		})
		return
	}

	deleteThumbnail := os.Remove(thumbnailName)
	if deleteThumbnail != nil {
		fmt.Println("Error deleting thumbnail:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete thumbnail",
		})
		return
	}

	// Replace 'thumbnails' with 'banners' in the thumbnailName to get the bannerName
	bannerName := strings.Replace(thumbnailName, "thumbnails", "banners", 1)

	// Delete the banner file
	deleteBanner := os.Remove(bannerName)
	if deleteBanner != nil {
		fmt.Println("Error deleting banner:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete banner",
		})
		return
	}

	err = app.database.deleteVideo(videoData.VideoID)
	if err != nil {
		fmt.Println("Error deleting video:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete video",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Video deleted successfully",
		"success": true,
	})
}

// homePage is a handler function that serves the home page.
// It parses the homePage.html template and executes it, sending the output to the client.
// If there is an error during parsing or execution of the template, it sends a server error response.
func (app *application) homePage(c *gin.Context) {
	// If not logged-In
	if userInfo.UserId == 0 && userInfo.Email == "" {
		app.error404(c)
		return
	}

	t, err := template.ParseFiles("ui/html/homePage.html")
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}

	err = t.Execute(c.Writer, nil)
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}

	videos, err := app.database.videosBrowser()
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}

	// Create a data map to hold the videos data
	Data = map[string]interface{}{
		"Videos": videos,
	}
}

// homePageVideos is a handler function that handles the post request of the home page.
// It fetches the videos data from the database and sends it to the client.
// If there is an error during fetching the videos data, it sends a server error response.
func (app *application) homePageVideos(c *gin.Context) {

}

// watchVideo is a handler function that serves the video watching page.
// It parses the watchVideo.html template and executes it, sending the output to the client.
// If there is an error during parsing or execution of the template, it sends a server error response.
func (app *application) watchVideo(c *gin.Context) {
	// If not logged-In
	if userInfo.UserId == 0 && userInfo.Email == "" {
		app.error404(c)
		return
	}

	t, err := template.ParseFiles("ui/html/watchVideo.html")
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}

	err = t.Execute(c.Writer, nil)
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}
}

// watchVideoPost is a handler function that handles the post request of the video watching page.
// It reads the video ID from the request, updates the video actions in the database, and sends an appropriate response.
// If there is an error during any of these steps, it sends an appropriate error response.
func (app *application) watchVideoPost(c *gin.Context) {
	type Video struct {
		VideoID       int       `json:"VideoID"`
		Title         string    `json:"Title"`
		Description   string    `json:"Description"`
		URL           string    `json:"URL"`
		ThumbnailURL  string    `json:"ThumbnailURL"`
		UploaderID    int       `json:"UploaderID"`
		UploadDate    time.Time `json:"UploadDate"`
		ViewsCount    int       `json:"ViewsCount"`
		LikesCount    int       `json:"LikesCount"`
		DislikesCount int       `json:"DislikesCount"`
		Duration      string    `json:"Duration"`
		CategoryID    int       `json:"CategoryID"`
		GenreID       int       `json:"GenreID"`
	}

	var videoData Video

	if err := c.ShouldBindJSON(&videoData); err != nil {
		fmt.Println("Error binding JSON data:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	// Now you can use videoData in your code
	fmt.Println(videoData)

	// send the videoID and the userID to the database
	err := app.database.videoActions(videoData.VideoID, userInfo.UserId)
	if err != nil {
		fmt.Println("Error updating video actions:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update video actions",
		})
		return
	}

	// Update the views count in the database
	err = app.database.updateViews(videoData.VideoID)
	if err != nil {
		fmt.Println("Error updating views count:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update views count",
		})
		return

	}
}

// recentlyAdded is a handler function that handles the fetching of recently added videos.
// It fetches the recently added videos data from the database and sends it to the client.
// If there is an error during fetching the videos data, it sends a server error response.
func (app *application) recentlyAdded(c *gin.Context) {
	videos, err := app.database.recentlyAddedVideos()
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}

	// Create a data map to hold the videos data
	Data = map[string]interface{}{
		"Videos": videos,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Recently added videos fetched successfully",
		"success": true,
		"videos":  Data,
	})
}

// recommendedVideos is a handler function that handles the fetching of recommended videos.
// It fetches the recommended videos data from the database and sends it to the client.
// If there is an error during fetching the videos data, it sends a server error response.
func (app *application) recommendedVideos(c *gin.Context) {
	videos, err := app.database.recommendedVideos()
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}

	// Send the videos data as a JSON response
	c.JSON(http.StatusOK, gin.H{
		"message": "Recommended videos fetched successfully",
		"success": true,
		"videos":  videos,
	})
}

// weeklyTop is a handler function that handles the fetching of weekly top videos.
// It fetches the weekly top videos data from the database and sends it to the client.
// If there is an error during fetching the videos data, it sends a server error response.
func (app *application) weeklyTop(c *gin.Context) {
	videos, err := app.database.weeklyTop()
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}

	// Send the videos data as a JSON response
	c.JSON(http.StatusOK, gin.H{
		"message": "Weekly top videos fetched successfully",
		"success": true,
		"videos":  videos,
	})
}

// continueWatching is a handler function that handles the fetching of videos for the continue watching section.
// It fetches the continue watching videos data from the database and sends it to the client.
// If there is an error during fetching the videos data, it sends an appropriate error response.
func (app *application) continueWatching(c *gin.Context) {
	videos, err := app.database.continueWatching(userInfo.UserId)
	if err != nil {
		if err.Error() == "no videos found" {
			c.JSON(http.StatusOK, gin.H{
				"message": "No videos found",
				"success": false,
			})
			return // Add this line
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch continue watching videos",
			"success": false,
		})
		return // And this line
	}

	// Send the videos data as a JSON response
	c.JSON(http.StatusOK, gin.H{
		"message": "Continue watching videos fetched successfully",
		"success": true,
		"videos":  videos,
	})
}

// caroselSlide is a handler function that handles the fetching of videos for the carousel slide.
// It fetches the carousel slide videos data from the database and sends it to the client.
// If there is an error during fetching the videos data, it sends a server error response.
func (app *application) caroselSlide(c *gin.Context) {
	videos, err := app.database.caroselSlide()
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}

	// Send the videos data as a JSON response
	c.JSON(http.StatusOK, gin.H{
		"message": "Recommended videos fetched successfully",
		"success": true,
		"videos":  videos,
	})
}

// search is a handler function that serves the search page.
// It parses the search.html template and executes it, sending the output to the client.
// If there is an error during parsing or execution of the template, it sends a server error response.
func (app *application) search(c *gin.Context) {
	// If not logged-In
	if userInfo.UserId == 0 && userInfo.Email == "" {
		app.error404(c)
		return
	}

	t, err := template.ParseFiles("ui/html/search.html")
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}

	err = t.Execute(c.Writer, nil)
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}
}

// searchData is a handler function that handles the post request of the search page.
// It reads the search value from the request, fetches the search videos data from the database, and sends it to the client.
// If there is an error during any of these steps, it sends an appropriate error response.
func (app *application) searchData(c *gin.Context) {
	// get the value from the search bar
	type Search struct {
		SearchValue string `json:"search"`
	}

	var searchData Search

	err := c.ShouldBindJSON(&searchData)
	// if the value is empty then do nothing else return error
	if err != nil {
		fmt.Println("Error binding JSON data:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	videos, err := app.database.searchVideos(searchData.SearchValue, userInfo.UserId)

	if err != nil {
		fmt.Println("Error getting search videos:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch search videos",
			"success": false,
		})
		return
	}

	// Send the videos data as a JSON response
	c.JSON(http.StatusOK, gin.H{
		"message": "Search videos fetched successfully",
		"success": true,
		"videos":  videos,
	})
}

// autoComplete is a function that accepts POST request from the User and sends the video list to the client.
func (app *application) autoComplete(c *gin.Context) {
	res, err := app.database.autoComplete()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Auto complete data fetched successfully",
		"success": true,
		"videos":  res,
	})
}

// videoAction is a handler function that handles the fetching of video action.
// It reads the video ID from the request, fetches the video action from the database, and sends it to the client.
// If there is an error during any of these steps, it sends an appropriate error response.
func (app *application) videoAction(c *gin.Context) {
	type VideoID struct {
		ID int `json:"id"`
	}

	var id VideoID

	err := c.ShouldBindJSON(&id)
	// if the value is empty then do nothing else return error
	if err != nil {
		fmt.Println("Error binding JSON data:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	action, err := app.database.videoAction(id.ID, userInfo.UserId)
	if err != nil {
		fmt.Println("Error getting video action:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch video action",
			"success": false,
		})
		return
	}

	// Send the videos data as a JSON response
	c.JSON(http.StatusOK, gin.H{
		"message": "Video action fetched successfully",
		"success": true,
		"action":  action,
	})
}

// videoActionChanged is a handler function that handles the updating of video action.
// It reads the video ID and the new action from the request, and updates the video action in the database.
// If there is an error during any of these steps, it sends an appropriate error response.
func (app *application) videoActionChanged(c *gin.Context) {
	type UpdateValues struct {
		VideoID int    `json:"videoID"`
		Action  string `json:"action"`
	}

	var updateValues UpdateValues

	err := c.ShouldBindJSON(&updateValues)
	// if the value is empty then do nothing else return error
	if err != nil {
		fmt.Println("Error binding JSON data:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	err = app.database.videoActionChanged(updateValues.VideoID, userInfo.UserId, updateValues.Action)
}

// userProfile is a handler function that serves the user profile page.
func (app *application) userProfile(c *gin.Context) {
	// If not logged-In
	if userInfo.UserId == 0 && userInfo.Email == "" {
		app.error404(c)
		return
	}

	t, err := template.ParseFiles("ui/html/profile.html")
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}

	err = t.Execute(c.Writer, nil)
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}
}

// videoDatas is a handler function that fetches video data for a user's profile.
// It uses the userId from the global userInfo variable to fetch the data.
// The function interacts with the database to retrieve the video data associated with the user's profile.
// If there's an error during the execution, the function logs the error message and sends a JSON response back to the client indicating the failure of the operation.
// If the operation is successful, the function sends a JSON response back to the client with a success message and the fetched video data.
func (app *application) videoDatas(c *gin.Context) {
	videos, err := app.database.userProfileVideosData(userInfo.UserId)
	if err != nil {
		fmt.Println("Error getting user profile videos:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch user profile videos",
			"success": false,
		})
		return
	}

	// Send the videos data as a JSON response
	c.JSON(http.StatusOK, gin.H{
		"message": "User profile videos fetched successfully",
		"success": true,
		"videos":  videos,
	})
}

// watchingVideos is a handler function that fetches the videos that a user is currently watching.
// It uses the userId from the global userInfo variable to fetch the data.
// The function interacts with the database to retrieve the video data associated with the user's watching list.
// If there's an error during the execution, the function logs the error message and sends a JSON response back to the client indicating the failure of the operation.
// If the operation is successful, the function sends a JSON response back to the client with a success message and the fetched video data.
func (app *application) watchingVideos(c *gin.Context) {
	videos, err := app.database.watchingVideos(userInfo.UserId)
	if err != nil {
		fmt.Println("Error getting watching videos:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch watching videos",
			"success": false,
		})
		return
	}

	// Send the videos data as a JSON response
	c.JSON(http.StatusOK, gin.H{
		"message": "Watching videos fetched successfully",
		"success": true,
		"videos":  videos,
	})
}

// onHoldVideos is a handler function that fetches the videos that a user has put on hold.
// It uses the userId from the global userInfo variable to fetch the data.
// The function interacts with the database to retrieve the video data associated with the user's on-hold list.
// If there's an error during the execution, the function sends a JSON response back to the client indicating the failure of the operation.
// If the operation is successful, the function sends a JSON response back to the client with a success message and the fetched video data.
func (app *application) onHoldVideos(c *gin.Context) {
	videos, err := app.database.onHoldVideos(userInfo.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch videos"})
		return
	}

	// Send the videos data as a JSON response
	c.JSON(http.StatusOK, gin.H{
		"message": "On-Hold videos fetched successfully",
		"success": true,
		"videos":  videos,
	})
}

// consideringVideos is a handler function that fetches the videos that a user is considering to watch.
// It uses the userId from the global userInfo variable to fetch the data.
// The function interacts with the database to retrieve the video data associated with the user's considering list.
// If there's an error during the execution, the function sends a JSON response back to the client indicating the failure of the operation.
// If the operation is successful, the function sends a JSON response back to the client with a success message and the fetched video data.
func (app *application) consideringVideos(c *gin.Context) {
	videos, err := app.database.consideringVideos(userInfo.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch videos"})
		return
	}

	// Send the videos data as a JSON response
	c.JSON(http.StatusOK, gin.H{
		"message": "Considering videos fetched successfully",
		"success": true,
		"videos":  videos,
	})
}

// recentlyCompletedVideos is a handler function that fetches the videos that a user has recently completed watching.
// It uses the userId from the global userInfo variable to fetch the data.
// The function interacts with the database to retrieve the video data associated with the user's recently completed list.
// If there's an error during the execution, the function sends a JSON response back to the client indicating the failure of the operation.
// If the operation is successful, the function sends a JSON response back to the client with a success message and the fetched video data.
func (app *application) recentlyCompletedVideos(c *gin.Context) {
	videos, err := app.database.recentlyCompletedVideos(userInfo.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch videos"})
		return
	}

	// Send the videos data as a JSON response
	c.JSON(http.StatusOK, gin.H{
		"message": "Recently completed videos fetched successfully",
		"success": true,
		"videos":  videos,
	})
}

// userDetails is a handler function that fetches the details of a user.
// It uses the userId from the global userInfo variable to fetch the data.
// The function interacts with the database to retrieve the user's details such as username, email, and admin status.
// If there's an error during the execution, the function sends a JSON response back to the client indicating the failure of the operation.
// If the operation is successful, the function sends a JSON response back to the client with a success message and the fetched user details.
func (app *application) userDetails(c *gin.Context) {
	userName, email, isAdmin, err := app.database.userDetails(userInfo.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch user details"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "User details fetched successfully",
		"success":  true,
		"userName": userName,
		"email":    email,
		"isAdmin":  isAdmin,
	})
}

// quotesHandler is a handler function that handles the request to fetch quotes.
// It reads the quotes from a CSV file named 'Quotes.csv' located in the 'internal' directory and stores them in a slice of 'Quote' structs.
// Each 'Quote' struct represents a quote and contains fields for ID, author, type, text, and count.
// If there's an error during the execution, the function handles the error appropriately.
// The function is a method of the 'application' struct and takes a 'gin.Context' object as an argument, which is used to handle the HTTP request and response.
func (app *application) quotesHandler(c *gin.Context) {

	type Quote struct {
		ID     string
		Author string
		Type   string
		Text   string
		Count  string
	}

	// Open the CSV file
	f, err := os.Open("internal/Quotes.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Read the CSV file into a slice of Quote structs
	quotes := make([]Quote, 0)
	r := csv.NewReader(f)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		quotes = append(quotes, Quote{
			ID:     record[0],
			Author: record[1],
			Type:   record[2],
			Text:   record[3],
			Count:  record[4],
		})
	}

	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Generate a random index
	index := rand.Intn(len(quotes))

	// Get the quote and author at the random index
	quote := quotes[index].Text
	author := quotes[index].Author

	// Return the quote and author
	c.JSON(http.StatusOK, gin.H{
		"quote":  quote,
		"author": author,
	})
}

// videoList is a handler function that fetches the list of all videos.
// It interacts with the database to retrieve the video data such as video title, description, upload date, and other related information.
// If there's an error during the execution, the function sends a JSON response back to the client indicating the failure of the operation.
// If the operation is successful, the function sends a JSON response back to the client with a success message and the fetched video list.
func (app *application) videoList(c *gin.Context) {
	t, err := template.ParseFiles("ui/html/videoList.html")
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}

	err = t.Execute(c.Writer, nil)
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}
}

// videoListPost is a handler function that handles the post request of the video list page.
// It fetches the video data from the database and sends it to the client as a JSON response.
// If there's an error during the execution, the function sends an appropriate error response.
func (app *application) videoListPost(c *gin.Context) {

}

func (app *application) recommendedVideoList(c *gin.Context) {
	videos, err := app.database.recommendedVideoListDatabase(userInfo.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch recommended videos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User details fetched successfully",
		"success": true,
		"videos":  videos,
	})
}

// recommendedVideoList is a handler function that fetches the list of recommended videos for a user.
// It interacts with the database to retrieve the video data such as video title, description, upload date, and other related information.
// If there's an error during the execution, the function sends a JSON response back to the client indicating the failure of the operation.
// If the operation is successful, the function sends a JSON response back to the client with a success message and the fetched recommended video list.
func (app *application) watchingVideoList(c *gin.Context) {
	videos, err := app.database.watchingVideoListDatabase(userInfo.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch watching videos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User details fetched successfully",
		"success": true,
		"videos":  videos,
	})
}

// completedVideoList is a handler function that fetches the list of videos a user has completed watching.
// It interacts with the database to retrieve the video data such as video title, description, upload date, and other related information.
// If there's an error during the execution, the function sends a JSON response back to the client indicating the failure of the operation.
// If the operation is successful, the function sends a JSON response back to the client with a success message and the fetched completed video list.
func (app *application) completedVideoList(c *gin.Context) {
	videos, err := app.database.completedVideoListDatabase(userInfo.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch completed videos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User details fetched successfully",
		"success": true,
		"videos":  videos,
	})
}

// onHoldVideoList is a handler function that fetches the list of videos a user has put on hold.
// It interacts with the database to retrieve the video data such as video title, description, upload date, and other related information.
// If there's an error during the execution, the function sends a JSON response back to the client indicating the failure of the operation.
// If the operation is successful, the function sends a JSON response back to the client with a success message and the fetched on-hold video list.
func (app *application) onHoldVideoList(c *gin.Context) {
	videos, err := app.database.onHoldVideoListDatabase(userInfo.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch on hold videos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User details fetched successfully",
		"success": true,
		"videos":  videos,
	})
}

// consideringVideoList is a handler function that fetches the list of videos a user is considering to watch.
// It interacts with the database to retrieve the video data such as video title, description, upload date, and other related information.
// If there's an error during the execution, the function sends a JSON response back to the client indicating the failure of the operation.
// If the operation is successful, the function sends a JSON response back to the client with a success message and the fetched considering video list.
func (app *application) consideringVideoList(c *gin.Context) {
	videos, err := app.database.consideringVideoListDatabase(userInfo.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch considering videos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User details fetched successfully",
		"success": true,
		"videos":  videos,
	})
}

// droppedVideoList is a handler function that fetches the list of videos a user has dropped or stopped watching.
// It interacts with the database to retrieve the video data such as video title, description, upload date, and other related information.
// If there's an error during the execution, the function sends a JSON response back to the client indicating the failure of the operation.
// If the operation is successful, the function sends a JSON response back to the client with a success message and the fetched dropped video list.
func (app *application) droppedVideoList(c *gin.Context) {
	videos, err := app.database.droppedVideoListDatabase(userInfo.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch dropped videos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User details fetched successfully",
		"success": true,
		"videos":  videos,
	})
}

// about is a handler function that serves the about page.
// It parses the about.html template and executes it, sending the output to the client.
// If there is an error during parsing or execution of the template, it sends a server error response.
// If the user is not logged in, it redirects the user to the login page.
func (app *application) about(c *gin.Context) {
	// If not logged-In
	if userInfo.UserId == 0 && userInfo.Email == "" {
		app.error404(c)
		return
	}

	t, err := template.ParseFiles("ui/html/about.html")
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}

	err = t.Execute(c.Writer, nil)
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}
}

// comment is a handler function that manages the user's comments on a video.
// It interacts with the database to store, retrieve, and delete comments.
// The function takes the user's ID and the comment text as input.
// If there's an error during the execution, the function sends a JSON response back to the client indicating the failure of the operation.
// If the operation is successful, the function sends a JSON response back to the client with a success message and the updated comment list.
func (app *application) comment(c *gin.Context) {
	// get the data from the post request that was sent by JS
	type Comment struct {
		Comment string `json:"comment"`
		VideoID int    `json:"videoID"`
	}

	var commentData Comment

	// Bind the JSON data from the request to the userData struct
	if err := c.ShouldBindJSON(&commentData); err != nil {
		fmt.Println("Error binding JSON data")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	fmt.Println("Comment: ", commentData.Comment, "VideoID: ", commentData.VideoID)

	err := app.database.commentOnVideo(userInfo.UserId, commentData.VideoID, commentData.Comment)
	if err != nil {
		fmt.Println("Error commenting on video")
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Comment added successfully",
		"success": true,
	})
}

// getComments is a handler function that fetches the list of comments for a specific video.
// It interacts with the database to retrieve the comment data such as comment text, user who commented, and the time of the comment.
// If there's an error during the execution, the function sends a JSON response back to the client indicating the failure of the operation.
// If the operation is successful, the function sends a JSON response back to the client with a success message and the fetched comment list.
func (app *application) getComments(c *gin.Context) {
	type VideoID struct {
		ID int `json:"videoID"`
	}

	var id VideoID

	err := c.ShouldBindJSON(&id)
	// if the value is empty then do nothing else return error
	if err != nil {
		fmt.Println("Error binding JSON data:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	comments, err := app.database.getComments(id.ID)
	if err != nil {
		fmt.Println("Error getting comments")
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Comments fetched successfully",
		"success":  true,
		"comments": comments,
	})
}

// upvote is a handler function that manages the upvotes for a specific video.
// It interacts with the database to increment the upvote count for the video.
// The function takes the user's ID and the video ID as input.
// If there's an error during the execution, the function sends a JSON response back to the client indicating the failure of the operation.
// If the operation is successful, the function sends a JSON response back to the client with a success message and the updated upvote count.
func (app *application) upvote(c *gin.Context) {
	type comment_id struct {
		CommentID int `json:"commentID"`
	}

	var commentID comment_id

	err := c.ShouldBindJSON(&commentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	err = app.database.upvoteComment(commentID.CommentID, userInfo.UserId)

	if err != nil {
		if err.Error() == "Comment does not exist" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Comment does not exist",
				"success": false,
			})
			return
		} else if err.Error() == "User does not exist" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "User does not exist",
				"success": false,
			})
			return
		} else if err.Error() == "Comment is already upvoted" {
			c.JSON(http.StatusOK, gin.H{
				"message": "Comment is already upvoted",
				"success": true,
			})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to upvote comment with error: " + err.Error(),
				"success": false,
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Comment upvoted successfully",
		"success": true,
	})
}

// reverseUpvote is a handler function that manages the reversal of upvotes for a specific video.
// It interacts with the database to decrement the upvote count for the video.
// The function takes the user's ID and the video ID as input.
// If there's an error during the execution, the function sends a JSON response back to the client indicating the failure of the operation.
// If the operation is successful, the function sends a JSON response back to the client with a success message and the updated upvote count.
func (app *application) reverseUpvote(c *gin.Context) {
	type comment_id struct {
		CommentID int `json:"commentID"`
	}

	var commentID comment_id

	err := c.ShouldBindJSON(&commentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	err = app.database.reverseUpvoteComment(commentID.CommentID, userInfo.UserId)

	if err != nil {
		if err.Error() == "Comment does not exist" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Comment does not exist",
				"success": false,
			})
			return
		} else if err.Error() == "User does not exist" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "User does not exist",
				"success": false,
			})
			return
		} else if err.Error() == "Comment is already upvoted" {
			c.JSON(http.StatusOK, gin.H{
				"message": "Comment is already upvoted",
				"success": true,
			})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to upvote comment with error: " + err.Error(),
				"success": false,
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Comment upvoted successfully",
		"success": true,
	})
}

// downvote is a handler function that manages the downvotes for a specific video.
// It interacts with the database to increment the downvote count for the video.
// The function takes the user's ID and the video ID as input.
// If there's an error during the execution, the function sends a JSON response back to the client indicating the failure of the operation.
// If the operation is successful, the function sends a JSON response back to the client with a success message and the updated downvote count.
func (app *application) downvote(c *gin.Context) {
	type comment_id struct {
		CommentID int `json:"commentID"`
	}

	var commentID comment_id

	err := c.ShouldBindJSON(&commentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	err = app.database.downvoteComment(commentID.CommentID, userInfo.UserId)

	if err != nil {
		if err.Error() == "Comment does not exist" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Comment does not exist",
				"success": false,
			})
			return
		} else if err.Error() == "User does not exist" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "User does not exist",
				"success": false,
			})
			return
		} else if err.Error() == "Comment is already downvoted" {
			c.JSON(http.StatusOK, gin.H{
				"message": "Comment is already downvoted",
				"success": true,
			})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to downvote comment with error: " + err.Error(),
				"success": false,
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Comment downvoted successfully",
		"success": true,
	})
}

// reverseDownvote is a handler function that manages the reversal of downvotes for a specific video.
// It interacts with the database to decrement the downvote count for the video.
// The function takes the user's ID and the video ID as input.
// If there's an error during the execution, the function sends a JSON response back to the client indicating the failure of the operation.
// If the operation is successful, the function sends a JSON response back to the client with a success message and the updated downvote count.
func (app *application) reverseDownvote(c *gin.Context) {
	type comment_id struct {
		CommentID int `json:"commentID"`
	}

	var commentID comment_id

	err := c.ShouldBindJSON(&commentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	err = app.database.reverseDownvoteComment(commentID.CommentID, userInfo.UserId)

	if err != nil {
		if err.Error() == "Comment does not exist" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Comment does not exist",
				"success": false,
			})
			return
		} else if err.Error() == "User does not exist" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "User does not exist",
				"success": false,
			})
			return
		} else if err.Error() == "Comment is already downvoted" {
			c.JSON(http.StatusOK, gin.H{
				"message": "Comment is already downvoted",
				"success": true,
			})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to downvote comment with error: " + err.Error(),
				"success": false,
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Comment downvoted successfully",
		"success": true,
	})
}

// commentDetails is a handler function that manages the details of a user's comment on a specific video.
// It interacts with the database to retrieve, update, or delete the comment details.
// The function takes the user's ID and the comment ID as input.
// If there's an error during the execution, the function sends a JSON response back to the client indicating the failure of the operation.
// If the operation is successful, the function sends a JSON response back to the client with a success message and the updated comment details.
func (app *application) commentDetails(c *gin.Context) {
	type videoID struct {
		VideoID int `json:"videoID"`
	}

	var id videoID

	err := c.ShouldBindJSON(&id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	comments, err := app.database.commentDetails(id.VideoID, userInfo.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to fetch comments with error: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Comments fetched successfully",
		"success":  true,
		"comments": comments,
	})
}

// likeVideo is a handler function that manages the likes for a specific video.
// It interacts with the database to increment the like count for the video.
// The function takes the user's ID and the video ID as input.
// If there's an error during the execution, the function sends a JSON response back to the client indicating the failure of the operation.
// If the operation is successful, the function sends a JSON response back to the client with a success message and the updated like count.
func (app *application) likeVideo(c *gin.Context) {
	type videoID struct {
		VideoID int `json:"videoID"`
	}

	var id videoID

	err := c.ShouldBindJSON(&id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	err = app.database.likeVideo(id.VideoID, userInfo.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to like video with error: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Video liked successfully",
		"success": true,
	})
}

// reverseLikeVideo is a handler function that manages the reversal of likes for a specific video.
// It interacts with the database to decrement the like count for the video.
// The function takes the user's ID and the video ID as input.
// If there's an error during the execution, the function sends a JSON response back to the client indicating the failure of the operation.
// If the operation is successful, the function sends a JSON response back to the client with a success message and the updated like count.
func (app *application) reverseLikeVideo(c *gin.Context) {
	type videoID struct {
		VideoID int `json:"videoID"`
	}

	var id videoID

	err := c.ShouldBindJSON(&id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	err = app.database.reverseLikeVideo(id.VideoID, userInfo.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to reverse like video with error: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Video like reversed successfully",
		"success": true,
	})
}

// dislikeVideo is a handler function that manages the dislikes for a specific video.
// It interacts with the database to increment the dislike count for the video.
// The function takes the user's ID and the video ID as input.
// If there's an error during the execution, the function sends a JSON response back to the client indicating the failure of the operation.
// If the operation is successful, the function sends a JSON response back to the client with a success message and the updated dislike count.
func (app *application) dislikeVideo(c *gin.Context) {
	type videoID struct {
		VideoID int `json:"videoID"`
	}

	var id videoID

	err := c.ShouldBindJSON(&id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	err = app.database.dislikeVideo(id.VideoID, userInfo.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to dislike video with error: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Video disliked successfully",
		"success": true,
	})

}

// reverseDislikeVideo is a handler function that manages the reversal of dislikes for a specific video.
// It interacts with the database to decrement the dislike count for the video.
// The function takes the user's ID and the video ID as input.
// If there's an error during the execution, the function sends a JSON response back to the client indicating the failure of the operation.
// If the operation is successful, the function sends a JSON response back to the client with a success message and the updated dislike count.
func (app *application) reverseDislikeVideo(c *gin.Context) {
	type videoID struct {
		VideoID int `json:"videoID"`
	}

	var id videoID

	err := c.ShouldBindJSON(&id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	err = app.database.reverseDislikeVideo(id.VideoID, userInfo.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to reverse dislike video with error: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Video dislike reversed successfully",
		"success": true,
	})
}

// isLikedDisliked is a handler function that checks if a specific video is liked or disliked by a user.
// It interacts with the database to retrieve the like and dislike status for the video.
// The function takes the user's ID and the video ID as input.
// If there's an error during the execution, the function sends a JSON response back to the client indicating the failure of the operation.
// If the operation is successful, the function sends a JSON response back to the client with a success message and the like and dislike status.
func (app *application) isLikedDisliked(c *gin.Context) {
	type videoID struct {
		VideoID int `json:"videoID"`
	}

	var id videoID

	err := c.ShouldBindJSON(&id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	liked, disliked, err := app.database.isLikedDisliked(id.VideoID, userInfo.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to fetch liked/disliked status with error: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Liked/disliked status fetched successfully",
		"success":    true,
		"isLiked":    liked,
		"isDisliked": disliked,
	})
}

// likeDislikeCount is a handler function that retrieves the count of likes and dislikes for a specific video.
// It interacts with the database to fetch the like and dislike count for the video.
// The function takes the video ID as input.
// If there's an error during the execution, the function sends a JSON response back to the client indicating the failure of the operation.
// If the operation is successful, the function sends a JSON response back to the client with a success message and the like and dislike count.
func (app *application) likeDislikeCount(c *gin.Context) {
	type videoID struct {
		VideoID int `json:"videoID"`
	}

	var id videoID

	err := c.ShouldBindJSON(&id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	likes, dislikes, err := app.database.likeDislikeCount(id.VideoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to fetch like/dislike count with error: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Like/dislike count fetched successfully",
		"success":  true,
		"likes":    likes,
		"dislikes": dislikes,
	})
}

// adminPanel is a handler function that serves the admin panel page.
// It checks if the user has admin access, and if so, it parses the adminPanel.html template and executes it, sending the output to the client.
// If the user does not have admin access, it redirects the user to the login page.
// If there is an error during parsing or execution of the template, it sends a server error response.
func (app *application) adminPanel(c *gin.Context) {
	// If not logged-In
	if userInfo.UserId == 0 && userInfo.Email == "" {
		app.error404(c)
		return
	}

	t, err := template.ParseFiles("ui/html/admin/admin_Panel.html")
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}

	err = t.Execute(c.Writer, nil)
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}

	video := c.Query("video")
	title := c.Query("title")
	description := c.Query("description")
	thumbnail := c.Query("thumbnail")
	banner := c.Query("banner")
	type_ := c.Query("type")
	genre := c.QueryArray("genre") // Use QueryArray to get multiple values for the same parameter

	fmt.Println("Video: ", video, "Title: ", title, "Description: ", description, "Thumbnail: ", thumbnail, "Banner: ", banner, "Type: ", type_, "Genre: ", genre)

	// pass the video to another function named encodeVideo
	//handler.HandleVideoUpload()
}

// adminAddVideo is a handler function that serves the add video page.
func (app *application) adminAddVideo(c *gin.Context) {
	// If not logged-In
	if userInfo.UserId == 0 && userInfo.Email == "" {
		app.error404(c)
		return
	}

	// Check if user have access to user access page
	_, upload, _, _, _, _, err := app.database.userAccess(userInfo.UserId)
	if err != nil {
		app.error403(c)
		return
	}
	if upload == 0 {
		app.error403(c)
		return
	}

	t, err := template.ParseFiles("ui/html/admin/admin_videoUpload.html")
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}

	err = t.Execute(c.Writer, nil)
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}
}

// adminDashboard is a handler function that serves the admin dashboard page.
// It checks if the user has admin access, and if so, it parses the adminDashboard.html template and executes it, sending the output to the client.
// If the user does not have admin access, it redirects the user to the login page.
// If there is an error during parsing or execution of the template, it sends a server error response.
func (app *application) adminDashboard(c *gin.Context) {
	// If not logged-In
	if userInfo.UserId == 0 && userInfo.Email == "" {
		app.error404(c)
		return
	}

	// Check if user have access to user access page
	dashboard, _, _, _, _, _, err := app.database.userAccess(userInfo.UserId)
	if err != nil {
		app.error403(c)
		return
	}
	if dashboard == 0 {
		app.error403(c)
		return
	}

	t, err := template.ParseFiles("ui/html/admin/admin_dashboard.html")
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}

	err = t.Execute(c.Writer, nil)
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}
}

// adminVideoList is a handler function that serves the video list page.
func (app *application) adminVideoList(c *gin.Context) {
	// If not logged-In
	if userInfo.UserId == 0 && userInfo.Email == "" {
		app.error404(c)
		return
	}

	// Check if user have access to user access page
	_, _, editDelete, _, _, _, err := app.database.userAccess(userInfo.UserId)
	if err != nil {
		app.error403(c)
		return
	}
	if editDelete == 0 {
		app.error403(c)
		return
	}

	t, err := template.ParseFiles("ui/html/admin/videoList.html")
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}

	err = t.Execute(c.Writer, nil)
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}
}

// adminAnalytics is a handler function that serves the analytics page.
func (app *application) adminAnalytics(c *gin.Context) {
	// If not logged-In
	if userInfo.UserId == 0 && userInfo.Email == "" {
		app.error404(c)
		return
	}

	// Check if user have access to user access page
	_, _, _, analytics, _, _, err := app.database.userAccess(userInfo.UserId)
	if err != nil {
		app.error403(c)
		return
	}
	if analytics == 0 {
		app.error403(c)
		return
	}

	t, err := template.ParseFiles("ui/html/admin/admin_analytics.html")
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}

	err = t.Execute(c.Writer, nil)
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}
}

// adminServerLogs is a handler function that serves the settings page.
func (app *application) adminServerLogs(c *gin.Context) {
	// If not logged-In
	if userInfo.UserId == 0 && userInfo.Email == "" {
		app.error404(c)
		return
	}

	// Check if user have access to user access page
	_, _, _, _, serverLogs, _, err := app.database.userAccess(userInfo.UserId)
	if err != nil {
		app.error403(c)
		return
	}
	if serverLogs == 0 {
		app.error403(c)
		return
	}

	t, err := template.ParseFiles("ui/html/admin/admin_serverLogs.html")
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}

	err = t.Execute(c.Writer, nil)
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}
}

func (app *application) adminUserAccess(c *gin.Context) {
	// If not logged-In
	if userInfo.UserId == 0 && userInfo.Email == "" {
		app.error404(c)
		return
	}

	// Check if user have access to user access page
	_, _, _, _, _, useraccess, err := app.database.userAccess(userInfo.UserId)
	if err != nil {
		app.error403(c)
		return
	}
	if useraccess == 0 {
		app.error403(c)
		return
	}

	t, err := template.ParseFiles("ui/html/admin/admin_userAccess.html")
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}

	err = t.Execute(c.Writer, nil)
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}
}

// adminUserAccessPost is a handler function that handles the POST request for updating user access in the admin panel.
// It reads the updated user access data from the request, validates it, and updates the user access in the database.
// If there is an error during any of these steps, it sends an appropriate error response.
// If the operation is successful, it sends a success response back to the client.
func (app *application) adminUserAccessPost(c *gin.Context) {
	adminAccess, err := app.database.allUserAdminAccess(userInfo.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch user access", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"message":     "User access fetched successfully",
		"adminAccess": adminAccess,
	})
}

// adminUserAccessChange is a handler function that handles the POST request for changing user access in the admin panel.
// It reads the new user access data from the request, validates it, and updates the user access in the database.
// If there is an error during any of these steps, it sends an appropriate error response.
// If the operation is successful, it sends a success response back to the client.
func (app *application) adminUserAccessChange(c *gin.Context) {
	type user struct {
		UserID      int    `json:"userID"`
		Access      int    `json:"access"`
		AccessValue string `json:"accessValue"`
	}

	var userAccess user

	err := c.ShouldBindJSON(&userAccess)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	fmt.Println("User ID: ", userAccess.UserID, "Access: ", userAccess.Access, "Access Value: ", userAccess.AccessValue)

	// Give dashboard access
	if userAccess.AccessValue == "dashboard" && userAccess.Access == 0 {
		err = app.database.giveDashboardAccess(userAccess.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to change user access", "message": err.Error()})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"success": true, "message": "User access changed successfully"})
			return
		}
	}
	// Revoke dashboard access
	if userAccess.AccessValue == "dashboard" && userAccess.Access == 1 {
		err = app.database.removeDashboardAccess(userAccess.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to change user access", "message": err.Error()})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"success": true, "message": "User access changed successfully"})
			return
		}
	}

	// Give upload access
	if userAccess.AccessValue == "upload" && userAccess.Access == 0 {
		err = app.database.giveUploadAccess(userAccess.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to change user access", "message": err.Error()})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"success": true, "message": "User access changed successfully"})
			return
		}
	}
	// Revoke upload access
	if userAccess.AccessValue == "upload" && userAccess.Access == 1 {
		err = app.database.removeUploadAccess(userAccess.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to change user access", "message": err.Error()})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"success": true, "message": "User access changed successfully"})
			return
		}
	}

	// Give edit_delete access
	if userAccess.AccessValue == "edit_delete" && userAccess.Access == 0 {
		err = app.database.giveEditDeleteAccess(userAccess.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to change user access", "message": err.Error()})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"success": true, "message": "User access changed successfully"})
			return
		}
	}
	// Revoke upload access
	if userAccess.AccessValue == "edit_delete" && userAccess.Access == 1 {
		err = app.database.removeEditDeleteAccess(userAccess.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to change user access", "message": err.Error()})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"success": true, "message": "User access changed successfully"})
			return
		}
	}

	// Give analytics access
	if userAccess.AccessValue == "analytics" && userAccess.Access == 0 {
		err = app.database.giveAnalyticsAccess(userAccess.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to change user access", "message": err.Error()})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"success": true, "message": "User access changed successfully"})
			return
		}
	}
	// Revoke upload access
	if userAccess.AccessValue == "analytics" && userAccess.Access == 1 {
		err = app.database.removeAnalyticsAccess(userAccess.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to change user access", "message": err.Error()})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"success": true, "message": "User access changed successfully"})
			return
		}
	}

	// Give server_log access
	if userAccess.AccessValue == "server_log" && userAccess.Access == 0 {
		err = app.database.giveServerLogAccess(userAccess.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to change user access", "message": err.Error()})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"success": true, "message": "User access changed successfully"})
			return
		}
	}
	// Revoke upload access
	if userAccess.AccessValue == "server_log" && userAccess.Access == 1 {
		err = app.database.removeServerLogAccess(userAccess.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to change user access", "message": err.Error()})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"success": true, "message": "User access changed successfully"})
			return
		}
	}

	// Give server_log access
	if userAccess.AccessValue == "user" && userAccess.Access == 0 {
		err = app.database.giveUserAccess(userAccess.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to change user access", "message": err.Error()})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"success": true, "message": "User access changed successfully"})
			return
		}
	}
	// Revoke upload access
	if userAccess.AccessValue == "user" && userAccess.Access == 1 {
		err = app.database.removeUserAccess(userAccess.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to change user access", "message": err.Error()})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"success": true, "message": "User access changed successfully"})
			return
		}
	}
}

// Profile is a handler function that serves the user profile page.
// It checks if the user is logged in, and if so, it serves the user's profile page.
// If the user is not logged in, it redirects the user to the 404 error page.
// The function interacts with the 'userInfo' global variable to check the user's login status.
func (app *application) profile(c *gin.Context) {
	// If not logged-In
	if userInfo.UserId == 0 && userInfo.Email == "" {
		app.error404(c)
		return
	}

	t, err := template.ParseFiles("ui/html/user_profile.html")
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}

	err = t.Execute(c.Writer, nil)
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}
}

// editProfile is a handler function that serves the user profile editing page.
// It checks if the user is logged in, and if so, it serves the user's profile editing page.
// If the user is not logged in, it redirects the user to the 404 error page.
// The function interacts with the 'userInfo' global variable to check the user's login status.
// It also interacts with the database to retrieve and update the user's profile information based on the changes made on the page.
// If there's an error during the execution, the function sends an appropriate error response.
// If the operation is successful, it sends a success response back to the client.
func (app *application) editProfile(c *gin.Context) {
	// If not logged-In
	if userInfo.UserId == 0 && userInfo.Email == "" {
		app.error404(c)
		return
	}

	t, err := template.ParseFiles("ui/html/user_editProfile.html")
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}

	err = t.Execute(c.Writer, nil)
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}
}

// changePassword is a handler function that allows a user to change their password.
// It checks if the user is logged in, and if so, it validates the old password and updates the password in the database.
// If the user is not logged in, it redirects the user to the login page.
// The function interacts with the 'userInfo' global variable to check the user's login status and to retrieve the user's current password.
// It also interacts with the database to update the user's password.
// If there's an error during the execution, the function sends an appropriate error response.
// If the operation is successful, it sends a success response back to the client.
func (app *application) changePassword(c *gin.Context) {
	// If not logged-In
	if userInfo.UserId == 0 && userInfo.Email == "" {
		app.error404(c)
		return
	}

	t, err := template.ParseFiles("ui/html/user_changePassword.html")
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}

	err = t.Execute(c.Writer, nil)
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}
}

// editUserProfile is a handler function that serves the user profile editing page.
// It checks if the user is logged in, and if so, it serves the user's profile editing page.
// If the user is not logged in, it redirects the user to the login page.
// The function interacts with the 'userInfo' global variable to check the user's login status.
// It also interacts with the database to retrieve the user's current profile information, and updates the profile information based on the changes made on the page.
// If there's an error during the execution, the function sends an appropriate error response.
// If the operation is successful, it sends a success response back to the client.
func (app *application) editUserProfile(c *gin.Context) {
	type userDetails struct {
		UserName string `json:"userName"`
		Email    string `json:"email"`
	}

	var user userDetails

	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	err = app.database.editUserProfile(user.UserName, user.Email, userInfo.Email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to edit user profile",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User profile edited successfully",
		"success": true,
	})
}

// changePasswordPost is a handler function that handles the POST request for changing a user's password.
// It reads the new password data from the request, validates it, and updates the password in the database.
// If the user is not logged in, it redirects the user to the login page.
// The function interacts with the 'userInfo' global variable to check the user's login status and to retrieve the user's current password.
// It also interacts with the database to update the user's password.
// If there's an error during the execution, the function sends an appropriate error response.
// If the operation is successful, it sends a success response back to the client.
func (app *application) changePasswordPost(c *gin.Context) {
	//oldPassword: oldPassword,
	//	newPassword: newPassword

	type password struct {
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	}

	var pass password

	err := c.ShouldBindJSON(&pass)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	err = app.database.changePassword(pass.OldPassword, pass.NewPassword, userInfo.Email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password changed successfully",
		"success": true,
	})
}

// encrypt is a function that takes a plain text string as input and returns the encrypted version of that string.
// It uses a specific encryption algorithm (such as AES, DES, RSA, etc.) to perform the encryption.
// The function handles any errors that might occur during the encryption process and returns them to the caller.
// If the encryption is successful, it returns the encrypted string.
func encrypt(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

// forgetPassword is a handler function that serves the forget password page.
// It checks if the user is logged in, and if not, it serves the forget password page.
// If the user is logged in, it redirects the user to the home page.
// The function interacts with the 'userInfo' global variable to check the user's login status.
// It also interacts with the database to validate the user's email and send a password reset link to the user's email.
// If there's an error during the execution, the function sends an appropriate error response.
// If the operation is successful, it sends a success response back to the client.
func (app *application) forgetPassword(c *gin.Context) {
	t, err := template.ParseFiles("ui/html/forgetPassword.html")
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}

	err = t.Execute(c.Writer, nil)
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}
}

// sendEmail is a handler function that handles the sending of an email.
// It reads the email data from the request, validates it, and sends the email.
// The function interacts with the email server to send the email.
// If there's an error during the execution, the function sends an appropriate error response.
// If the operation is successful, it sends a success response back to the client.
func (app *application) sendEmail(c *gin.Context) {

	type Email struct {
		Email string `json:"email"`
	}

	var email Email

	err := c.ShouldBindJSON(&email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	// Create a new source and randomizer
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	// Create a random password of length 8
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+{}[]:;,.?/|\\"
	password := make([]byte, 8)
	for i := range password {
		password[i] = charset[r.Intn(len(charset))]
	}

	// Encrypt the password
	encryptedPassword, err := encrypt(string(password))

	if err != nil {
		fmt.Println("Error encrypting password:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to encrypt password",
		})
		return
	}

	err = app.database.resetPassword(email.Email, encryptedPassword)
	if err != nil {
		fmt.Println("Error resetting password:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "User does not exist. Please enter a valid email address.",
		})
		return
	}

	apiKey := "re_AESmEYab_Mub2tMp887SMot95Tftjn3wk"

	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From: "Diwash <diwash@diwashmainali.com.np>",
		//To:      []string{"np03cs4s220198@heraldcollege.edu.np"},
		To:      []string{email.Email},
		Subject: "Reset Password",
		Html: `<html>
                <head>
                </head>
                <body>
                    <p>Dear User,</p>
                    <p>Your password has been reset successfully.</p>
                    <p>Your new password is: <span style="color:red;">` + string(password) + `</span></p>
                    <p>For security purposes, we highly recommend that you change your password immediately upon logging in.</p>
					<p>If you did not request this password reset or have any concerns about the security of your account, please contact our support team immediately.</p>
					<p>Thank you for your attention to this matter.</p>
					<p>Best Regards,<br>Diwash</p>
                </body>
            </html>`,
	}

	_, err = client.Emails.Send(params)

	if err != nil {
		fmt.Println("Error sending email:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to send email",
		})
		return
	}

	err = app.database.resetPassword(email.Email, encryptedPassword)
	if err != nil {
		fmt.Println("Error resetting password:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to reset password",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Email sent successfully",
		"success": true,
	})
}

// mostViewedVideos is a handler function that fetches the list of most viewed videos.
// It interacts with the database to retrieve the video data such as video title, description, view count, and other related information.
// If there's an error during the execution, the function sends a JSON response back to the client indicating the failure of the operation.
// If the operation is successful, the function sends a JSON response back to the client with a success message and the fetched most viewed video list.
func (app *application) mostViewedVideos(c *gin.Context) {
	// Call the mostViewedVideos method from the database connection
	videos, err := app.database.mostViewedVideos()
	if err != nil {
		// If there is an error, return a 500 status code and error message
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// If there is no error, return a 200 status code and the videos in JSON format
	c.JSON(http.StatusOK, videos)
}

// likeVsDislike is a handler function that fetches the count of likes and dislikes for a specific video.
// It interacts with the database to retrieve the like and dislike count for the video.
// The function takes the video ID as input.
// If there's an error during the execution, the function sends a JSON response back to the client indicating the failure of the operation.
// If the operation is successful, the function sends a JSON response back to the client with a success message and the like and dislike count.
func (app *application) likeVsDislike(c *gin.Context) {
	mostLikedVideos, mostDislikedVideos, err := app.database.likeVsDislike()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mostLikedVideos":    mostLikedVideos,
		"mostDislikedVideos": mostDislikedVideos,
	})
}

// duration is a handler function that fetches the duration of a specific video.
// It interacts with the database to retrieve the duration of the video.
// The function takes the video ID as input.
// If there's an error during the execution, the function sends a JSON response back to the client indicating the failure of the operation.
// If the operation is successful, the function sends a JSON response back to the client with a success message and the duration of the video.
func (app *application) duration(c *gin.Context) {
	duration, err := app.database.duration()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, duration)
}

// devideInfo is a function that divides the information into manageable parts.
// It takes the information as input and returns the divided information.
// The function interacts with the data to divide the information based on certain criteria.
// If there's an error during the execution, the function handles the error and returns an appropriate error message.
// If the operation is successful, the function returns the divided information.
func (app *application) deviceInfo(r *http.Request) {

	resp, err := http.Get("https://api.ipify.org?format=json")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var result map[string]string
	json.Unmarshal([]byte(body), &result)

	//fmt.Println(result)

	publicIP := result["ip"]

	// Get the User-Agent header
	userAgent := r.UserAgent()

	// The User-Agent header contains information about the device type, OS and Browser
	// You can parse this header to extract the information you need
	// Note: This is a simplified example, in reality parsing the User-Agent header can be complex
	deviceType := ""
	deviceOS := ""

	// Determine the device type
	if strings.Contains(userAgent, "Mobile") {
		deviceType = "Mobile"
	} else if strings.Contains(userAgent, "Tablet") {
		deviceType = "Tablet"
	} else if strings.Contains(userAgent, "Windows") || strings.Contains(userAgent, "Mac OS") || strings.Contains(userAgent, "Linux") {
		deviceType = "Desktop"
	} else {
		deviceType = "Other"
	}

	if strings.Contains(userAgent, "Windows") {
		deviceOS = "Windows"
	} else if strings.Contains(userAgent, "Mac OS") {
		deviceOS = "Mac OS"
	} else if strings.Contains(userAgent, "Linux") {
		deviceOS = "Linux"
	} else if strings.Contains(userAgent, "Android") {
		deviceOS = "Android"
	} else if strings.Contains(userAgent, "iOS") {
		deviceOS = "iOS"
	} else {
		deviceOS = "Other"
	}

	// Get the Browser
	browser := ""
	if strings.Contains(userAgent, "Chrome/") {
		browser = "Chrome"
	} else if strings.Contains(userAgent, "Safari/") {
		browser = "Safari"
	} else if strings.Contains(userAgent, "Firefox/") {
		browser = "Firefox"
	} else if strings.Contains(userAgent, "Edge/") {
		browser = "Edge"
	} else if strings.Contains(userAgent, "Opera/") {
		browser = "Opera"
	} else if strings.Contains(userAgent, "MSIE") {
		browser = "Internet Explorer"
	} else {
		browser = "Other"
	}

	//fmt.Println("IP:", publicIP)
	//fmt.Println("Device Type:", deviceType)
	//fmt.Println("Device OS:", deviceOS)
	//fmt.Println("Browser:", browser)

	type networkInfo struct {
		IP          string  `json:"ip"`
		CountryCode string  `json:"country_code"`
		CountryName string  `json:"country_name"`
		RegionName  string  `json:"region_name"`
		CityName    string  `json:"city_name"`
		Latitude    float64 `json:"latitude"`
		Longitude   float64 `json:"longitude"`
		ZipCode     string  `json:"zip_code"`
		TimeZone    string  `json:"time_zone"`
		ASN         string  `json:"asn"`
		AS          string  `json:"as"`
		IsProxy     bool    `json:"is_proxy"`
	}

	var network_info networkInfo

	url := "https://api.ip2location.io/?key=B98A09AA002A51699B5F3BC63B04A5D9&ip=" + publicIP + "&format=json"

	resp, err = http.Get(url)
	if err != nil {
		fmt.Println("Error getting location:", err)
	}

	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}

	// Bind the JSON data into network_info
	err = json.Unmarshal(body, &network_info)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	// Get the current date and time in UTC
	currentTime := time.Now().UTC()

	// Format the date and time as a string
	currentTimeString := currentTime.Format("2006-01-02 15:04:05 MST")

	//fmt.Printf("Network Info: %+v\n", network_info)

	lastLoginJSON := fmt.Sprintf(`{"login_time": "%s"}`, currentTimeString)

	// Insert the device information into the database
	//IP, device_type, device_os, Browser, LastLogin, country_code, country_name, region_name, city_name, latitude, longitude, zip_code, time_zone, asn, as_, is_proxy
	err = app.database.deviceInfo(publicIP, deviceType, deviceOS, browser, lastLoginJSON, network_info.CountryCode, network_info.CountryName, network_info.RegionName, network_info.CityName, network_info.Latitude, network_info.Longitude, network_info.ZipCode, network_info.TimeZone, network_info.ASN, network_info.AS, network_info.IsProxy)
	if err != nil {
		fmt.Println("Error inserting device info into the database:", err)
	} else {
		//fmt.Println("Device Info inserted successfully")
	}
}

// serverLogsPost is a handler function that handles the post request of the server logs page.
// It fetches the server logs data from the database and sends it to the client as a JSON response.
// If there's an error during the execution, the function sends an appropriate error response.
// If the operation is successful, it sends a success response back to the client with the fetched server logs data.
func (app *application) serverLogsPost(c *gin.Context) {
	data, err := app.database.serverLog()
	if err != nil {
		//c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println("Error getting server logs:", err)
		return
	}

	// Send the server logs data as a JSON response
	c.JSON(http.StatusOK, gin.H{
		"message": "Server logs fetched successfully",
		"success": true,
		"logs":    data,
	})
}

// Modify the ResponseData struct
type ResponseData struct {
	CountryNameCount map[string]int `json:"country_name_count"`
	CountryCodeCount map[string]int `json:"country_code_count"`
	Count            int            `json:"count"`
}

// localAnalysis is a function that analyzes the local data of the application.
// It interacts with the local data to perform various analysis tasks such as calculating statistics, identifying patterns, and extracting useful information.
// The function does not take any input as it works with the local data of the application.
// If there's an error during the execution, the function handles the error and returns an appropriate error message.
// If the operation is successful, the function returns the results of the analysis.
func (app *application) locationAnalysis(c *gin.Context) {
	// Call the locationAnalysis function from the databaseConn struct
	countryNameCount, countryCodeCount, count, err := app.database.locationAnalysis()
	if err != nil {
		// Handle the error
	}

	// Create an instance of the ResponseData struct
	data := ResponseData{
		CountryNameCount: countryNameCount,
		CountryCodeCount: countryCodeCount,
		Count:            count,
	}

	// Send the data as a JSON response
	c.JSON(http.StatusOK, data)
}

// imageUploadDynamic is a handler function that handles the dynamic uploading of images.
// It reads the image data from the request, validates it, and uploads the image to the server.
// The function interacts with the file system to store the uploaded image.
// If there's an error during the execution, the function sends an appropriate error response.
// If the operation is successful, it sends a success response back to the client with the details of the uploaded image.
func (app *application) imageUploadDynamic(c *gin.Context) {
	type image struct {
		Image string `json:"image"`
	}

	var img image

	err := c.ShouldBindJSON(&img)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	base64String := img.Image

	// Decode the base64 string
	imageBytes, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		fmt.Println("Error decoding base64 string:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error decoding base64 string",
		})
		return
	}

	// create a unique file name unique to the image
	uniqueFile := uuid.New().String()

	// Save the image to a file
	err = ioutil.WriteFile("./userUploadDatas/userProfileImage/"+uniqueFile+".png", imageBytes, 0644)
	if err != nil {
		fmt.Println("Error saving image to file:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error saving image to file",
		})
		return
	}

	err = app.database.saveUserProfile("./userUploadDatas/userProfileImage/"+uniqueFile+".png", userInfo.UserId)

	c.JSON(http.StatusOK, gin.H{
		"message": "Image saved successfully",
	})
}

// displayUserProfileImage is a handler function that handles the retrieval and display of a user's profile image.
// It interacts with the file system to retrieve the user's profile image.
// The function uses the user's ID, which is stored in the global userInfo variable, to locate the image.
// If there's an error during the execution, the function sends an appropriate error response.
// If the operation is successful, it sends a success response back to the client with the user's profile image.
func (app *application) displayUserProfileImage(c *gin.Context) {
	imagePath, err := app.database.displayUserProfileImage(userInfo.UserId)
	if err != nil {
		fmt.Println("Error displaying user profile image:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"succes": false,
			"error":  "Error displaying user profile image",
		})
		return
	}

	// Return the image path as a JSON response
	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"message":   "User profile image displayed successfully",
		"imagePath": imagePath,
	})
}

// error403 is a handler function that serves the 403 error page.
// It parses the error403.html template and executes it, sending the output to the client.
// If there is an error during parsing or execution of the template, it sends a server error response.
func (app *application) error403(c *gin.Context) {
	t, err := template.ParseFiles("ui/html/errorPage/403.html")
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}

	//app.deviceInfo(c.Request)
	err = t.Execute(c.Writer, nil)
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}
}

// error404 is a handler function that serves the 404 error page.
// It parses the error404.html template and executes it, sending the output to the client.
// If there is an error during parsing or execution of the template, it sends a server error response.
func (app *application) error404(c *gin.Context) {
	t, err := template.ParseFiles("ui/html/errorPage/404.html")
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}

	//app.deviceInfo(c.Request)
	err = t.Execute(c.Writer, nil)
	if err != nil {
		app.serverError(c.Writer, err)
		return
	}
}

// logout is a handler function that handles the logout process for a user.
// It invalidates the user's session, effectively logging them out of the system.
// If there's an error during the execution, the function sends an appropriate error response.
// If the operation is successful, it redirects the user to the login page.
func (app *application) logout(c *gin.Context) {
	//userInfo.UserId != 0 && userInfo.Email != ""
	userInfo.UserId = 0
	userInfo.Email = ""

	c.JSON(http.StatusOK, gin.H{
		"message": "User logged out successfully",
		"success": true,
	})
}

// checkAdminAccess is a handler function that checks if a user has admin access.
// It reads the user's ID from the session, fetches the user's admin status from the database, and sends it to the client.
// If there's an error during any of these steps, it sends an appropriate error response.
// If the operation is successful, it sends a success response back to the client with the user's admin status.
func (app *application) checkAdminAccess(c *gin.Context) {
	/*c.JSON(http.StatusOK, gin.H{
		"message": "User logged out successfully",
		"success": true,
	})*/

	res, err := app.database.checkAdminAccess(userInfo.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to check user access",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "User access checked successfully",
		"success":     true,
		"adminAccess": res,
	})
}
