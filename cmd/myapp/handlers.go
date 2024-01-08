package main

import (
	"Dekkoto/cmd/myapp/handler"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func (app *application) login(c *gin.Context) {
	t, err := template.ParseFiles("ui/html/login.html")
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

// register handles the GET and POST requests for the /register route. It also validates the JSON data and inserts it into the database.
func (app *application) register(c *gin.Context) {
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

	err := app.database.registerUser(userData.Name, userData.Email, userData.Password)

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

	c.JSON(http.StatusOK, gin.H{
		"message": "User created successfully",
	})

}

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
			"error": "Failed to login user",
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

func (app *application) uploadVideo(c *gin.Context) {
	videoInfo := handler.VideoDetailsInfo
	// get current date
	currentDate := time.Now().Format("2006-01-02")

	// Convert arrays to comma-separated strings
	categoryString := strings.Join(videoInfo.Genres, ",")

	// Map genre and category strings to their respective IDs
	var categoryID int
	// remove all the whitespaces
	videoInfo.Types = strings.ReplaceAll(videoInfo.Types, " ", "")
	categoryID, err := app.database.getCategoryID(categoryString)

	// check if category is empty
	if err != nil {
		c.String(500, "Failed to get category ID with error: "+err.Error())
	}

	//genreID, err := app.database.getGenreID(videoInfo.Genres[0])
	// there are multiple genres do it for all
	// remove all the whitespaces
	categoryString = strings.ReplaceAll(categoryString, " ", "")
	genreID, err := app.database.getGenreID(videoInfo.Types)
	if err != nil {
		c.String(500, "Failed to get genre ID with error: "+err.Error())
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to upload video",
			"success": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Video uploaded in the database successfully",
		"success": true,
	})
}

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

// Data := map[string]interface{}
var Data map[string]interface{}

func (app *application) showVideos(c *gin.Context) {
	t, err := template.ParseFiles("ui/html/adminTables.html")
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

func (app *application) showVideosPost(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Video uploaded in the database successfully",
		"success": true,
		"videos":  Data,
	})
}

func (app *application) editVideo(c *gin.Context) {
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

func (app *application) editVideoPost(c *gin.Context) {
	fmt.Println("edit video post")

	//app.editVideo(c)

	// get the data from the post request that was sent by JS
	type Video struct {
		VideoID     string `json:"videoID"`
		Title       string `json:"title"`
		Description string `json:"description"`
		CategoryID  string `json:"categoryID"`
		GenreID     string `json:"genreID"`
	}

	var videoData Video

	rawData, _ := c.GetRawData()
	fmt.Println(string(rawData))

	// Bind the JSON data from the request to the userData struct
	if err := c.ShouldBindJSON(&videoData); err != nil {
		fmt.Println("Error binding JSON data")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	fmt.Println(videoData)

	// convert id to int
	videoIDInt, err := strconv.Atoi(videoData.VideoID)
	if err != nil {
		fmt.Println("Error converting videoID to int")
	}

	genreIDInt, err := strconv.Atoi(videoData.GenreID)
	if err != nil {
		fmt.Println("Error converting genreID to int")
	}

	genreName, err := app.database.getGenreName(genreIDInt)
	if err != nil {
		fmt.Println("Error getting genre name")
		return
	}

	// send the post request to the another js file to show the data in the form
	c.JSON(http.StatusOK, gin.H{
		"message":     "Video data fetched successfully",
		"success":     true,
		"videoID":     videoIDInt,
		"title":       videoData.Title,
		"description": videoData.Description,
		"genreID":     genreName,
	})

}

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

	err := app.database.deleteVideo(videoData.VideoID)
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
