// Package main provides various handlers for handling user requests.
//
// The functions in this file handle user requests for login, registration, video upload, video termination,
// video display, video editing, video deletion, and other related operations.
// The handlers interact with the database and the user interface to provide the required functionality.
package main

import (
	"Dekkoto/cmd/myapp/handler"
	"encoding/csv"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"io"
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

// register is a handler function that serves the registration page.
// It parses the register.html template and executes it, sending the output to the client.
// If there is an error during parsing or execution of the template, it sends a server error response.
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

	fmt.Println("Email: ", userData.Email+" Password: "+userData.Password)

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

// uploadVideo is a handler function that handles the upload of a video.
// It reads the video data from the request, validates it, and uploads the video to the database.
// If there is an error during any of these steps, it sends an appropriate error response.
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

// showVideosPost is a handler function that handles the post request of the videos page.
// It sends the videos data to the client as a JSON response.
func (app *application) showVideosPost(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Video uploaded in the database successfully",
		"success": true,
		"videos":  Data,
	})
}

// editVideo is a handler function that serves the video editing page.
// It parses the adminEditVideo.html template and executes it, sending the output to the client.
// If there is an error during parsing or execution of the template, it sends a server error response.
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

func (app *application) userProfile(c *gin.Context) {
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

func (app *application) about(c *gin.Context) {
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

func (app *application) adminPanelTemp(c *gin.Context) {
	t, err := template.ParseFiles("ui/html/admin_PanelTest.html")
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

// adminAddVideoTemp is a handler function that serves the add video page.
func (app *application) adminAddVideoTemp(c *gin.Context) {
	t, err := template.ParseFiles("ui/html/admin_videoUpload.html")
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

func (app *application) adminDashboardTemp(c *gin.Context) {
	t, err := template.ParseFiles("ui/html/admin_dashboard.html")
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

// adminVideoListTemp is a handler function that serves the video list page.
func (app *application) adminVideoListTemp(c *gin.Context) {
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

// adminAnalyticsTemp is a handler function that serves the analytics page.
func (app *application) adminAnalyticsTemp(c *gin.Context) {
	t, err := template.ParseFiles("ui/html/admin_analytics.html")
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

// adminSettingsTemp is a handler function that serves the settings page.
func (app *application) adminSettingsTemp(c *gin.Context) {
	t, err := template.ParseFiles("ui/html/admin_settings.html")
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

func (app *application) profile(c *gin.Context) {
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

func (app *application) editProfile(c *gin.Context) {
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

func (app *application) changePassword(c *gin.Context) {
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
