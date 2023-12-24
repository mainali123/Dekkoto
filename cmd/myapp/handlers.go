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
	genresStr := strings.Join(videoInfo.Genres, ",")
	//typesStr := strings.Join(videoInfo.Types, ",")

	// Map genre and category strings to their respective IDs
	var categoryID int
	fmt.Println("Sending to database (types): ", videoInfo.Types)
	categoryID, err := app.database.getCategoryID(videoInfo.Types)

	// check if category is empty
	//if err != nil {
	//	c.String(500, "Failed to get category ID with error: "+err.Error())
	//}

	//genreID, err := app.database.getGenreID(videoInfo.Genres[0])
	// there are multiple genres do it for all
	fmt.Println("Sending to database (genres): ", genresStr)
	genreID, err := app.database.getGenreID(genresStr)
	//if err != nil {
	//	c.String(500, "Failed to get genre ID with error: "+err.Error())
	//}

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
			"error": "Failed to upload video",
		})
		c.String(500, "Failed to upload video with error "+err.Error())
		return
	}
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
