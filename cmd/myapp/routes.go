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

	router.GET("/adminPanel", app.admin)
	router.POST("/uploadVideo", handler.HandleVideoUpload)
	router.POST("/uploadThumbnail", handler.HandleThumbnailUpload)
	router.POST("/videoDetails", handler.VideoDetails)
	router.POST("/confirmVideo", app.uploadVideo)
	router.POST("/terminateVideo", app.terminateVideo)
}
