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
	//router.GET("/adminPanel", app.admin)
	router.POST("/mainUpload", handler.MainUpload)
	//router.POST("/uploadVideo", handler.HandleVideoUpload)
	//router.POST("/uploadThumbnail", handler.HandleThumbnailUpload)
	//router.POST("/uploadBanner", handler.HandleBannerUpload)
	//router.POST("/videoDetails", handler.VideoDetails)
	router.POST("/uploadVideoInDatabase", app.uploadVideo)
	router.POST("/terminateVideo", app.terminateVideo)

	// For admin video RUD operations
	router.GET("/adminPanel/showVideos", app.showVideos)
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
	router.POST("/autoComplete", app.autoComplete)

	// For user profile
	router.GET("/userProfile", app.userProfile)
	router.POST("/videoDatas", app.videoDatas)
	router.POST("/watchingVideos", app.watchingVideos)
	router.POST("/onHoldVideos", app.onHoldVideos)
	router.POST("/consideringVideos", app.consideringVideos)
	router.POST("/recentlyCompletedVideos", app.recentlyCompletedVideos)
	router.POST("/userDetails", app.userDetails)
	router.GET("/quotes", app.quotesHandler)

	// Videos List
	router.GET("/videoslist", app.videoList)
	router.POST("/videoslist", app.videoListPost)
	router.POST("/recommendedVideoList", app.recommendedVideoList)
	router.POST("/watchingVideoList", app.watchingVideoList)
	router.POST("/completedVideoList", app.completedVideoList)
	router.POST("/onHoldVideoList", app.onHoldVideoList)
	router.POST("/consideringVideoList", app.consideringVideoList)
	router.POST("/droppedVideoList", app.droppedVideoList)

	// About
	router.GET("/about", app.about)

	// Comments
	router.POST("/comment", app.comment)
	router.POST("/displayComments", app.getComments)
	router.POST("/upvote", app.upvote)
	router.POST("/downvote", app.downvote)
	router.POST("/commentDetails", app.commentDetails)

	// like & dislike video
	router.POST("/likeVideo", app.likeVideo)
	router.POST("/reverseLikeVideo", app.reverseLikeVideo)
	router.POST("/dislikeVideo", app.dislikeVideo)
	router.POST("/reverseDislikeVideo", app.reverseDislikeVideo)
	router.POST("/isLikedDisliked", app.isLikedDisliked)
	router.POST("/likeDislikeCount", app.likeDislikeCount)

	// For user settings
	router.GET("/profile", app.profile)
	router.GET("/editProfile", app.editProfile)
	router.GET("/changePassword", app.changePassword)
	router.POST("/editUserProfile", app.editUserProfile)
	router.POST("/changePassword", app.changePasswordPost)

	// TemporaryDevelopment of AdminPanel
	/*router.GET("/aakash", app.adminPanelTemp)
	router.GET("/aakash/dashboard", app.adminDashboard)
	router.GET("/aakash/addVideo", app.adminAddVideo)
	router.GET("/aakash/videoList", app.adminVideoList)
	router.GET("/aakash/analytics", app.adminAnalytics)
	router.GET("/aakash/serverLogs", app.adminServerLogs)*/
	router.GET("/adminPanel", app.adminPanel)
	router.GET("/adminPanel/dashboard", app.adminDashboard)
	router.GET("/adminPanel/addVideo", app.adminAddVideo)
	router.GET("/adminPanel/videoList", app.adminVideoList)
	router.GET("/adminPanel/analytics", app.adminAnalytics)
	router.GET("/adminPanel/serverLogs", app.adminServerLogs)

	// Forget Password
	router.GET("/forgetPassword", app.forgetPassword)
	router.POST("/sendEmail", app.sendEmail)

	// Analytics
	router.POST("/mostViewedVideos", app.mostViewedVideos)
	router.POST("/likeVsDislike", app.likeVsDislike)
	router.POST("/duration", app.duration)

	// ServerLogs
	router.POST("/serverLogsPost", app.serverLogsPost)

	// Location analysis
	router.POST("/locationAnalysis", app.locationAnalysis)

}
