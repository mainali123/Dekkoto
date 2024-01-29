// Package main provides various handlers for handling user requests.
//
// The functions in this file handle user requests for login, registration, video upload, video termination,
// video display, video editing, video deletion, and other related operations.
// The handlers interact with the database and the user interface to provide the required functionality.
//
// Global Variables:
// Data is a global variable of type map with string keys and values of any type (interface{}).
// It is used to hold various data, such as video information, that needs to be accessed across different functions.
package main

import (
	"Dekkoto/cmd/myapp/handler"
	"github.com/gin-gonic/gin"
)

// routes function is responsible for setting up the routes for the application.
// It takes a gin.Engine as an argument and sets up the routes for the application.
// The routes are defined for various operations like login, registration, video upload, video termination,
// video display, video editing, video deletion, and other related operations.
// Each route is associated with a handler function that handles the request for that route.
func (app *application) routes(router *gin.Engine) {

	// Define the route for the handlers
	router.GET("/", app.login)
	router.GET("/login", app.login)
	router.GET("/register", app.register)
	router.POST("/login", app.loginPostRequest)
	router.POST("/register", app.registerPostRequest)

	// For admin panel
	router.GET("/adminPanel", app.admin)
	router.POST("/uploadVideo", handler.HandleVideoUpload)
	router.POST("/uploadThumbnail", handler.HandleThumbnailUpload)
	router.POST("/uploadBanner", handler.HandleBannerUpload)
	router.POST("/videoDetails", handler.VideoDetails)
	router.POST("/confirmVideo", app.uploadVideo)
	router.POST("/terminateVideo", app.terminateVideo)

	// For admin video RUD operations
	router.GET("/showVideos", app.showVideos)
	router.POST("/showVideosPost", app.showVideosPost)
	router.GET("/editVideo", app.editVideo)
	router.POST("/editVideo", app.editVideoPost)

	router.POST("/showCategoriesName", app.showCategoriesName)
	router.POST("/showGenresName", app.showGenresName)

	router.POST("/editVideoPost", app.updateVideoDetails)
	router.POST("/deleteVideo", app.deleteVideo)

	// For homepage
	router.GET("/home", app.homePage)
	router.POST("/home", app.homePageVideos)

	// For watching videos
	router.GET("/watchVideo", app.watchVideo)
	router.POST("/watchVideoPost", app.watchVideoPost)
	router.POST("/recentlyAdded", app.recentlyAdded)
	router.POST("/recommendedVideos", app.recommendedVideos)
	router.POST("/weeklyTop", app.weeklyTop)
	router.POST("/continueWatching", app.continueWatching)
	router.POST("/caroselSlide", app.caroselSlide)

	// For video player
	router.POST("/videoAction", app.videoAction)
	router.POST("/videoActionChanged", app.videoActionChanged)

	// For search page
	router.GET("/search", app.search)
	router.POST("/searchData", app.searchData)

	// For user profile
	router.GET("/userProfile", app.userProfile)
}
