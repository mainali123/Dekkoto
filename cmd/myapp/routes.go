package main

import (
	"Dekkoto/cmd/myapp/handler"
	"github.com/gin-gonic/gin"
)

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
	//router.GET("/editVideo/:videoID/:title/:description/:categoryID/:genreID", app.editVideo)
	router.GET("/editVideo", app.editVideo)
	router.POST("/editVideo", app.editVideoPost)
	//router.POST("/editSelectedVideo", app.editVideoPost)

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
}
