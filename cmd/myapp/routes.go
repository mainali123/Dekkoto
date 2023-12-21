package main

import "github.com/gin-gonic/gin"

func (app *application) routes(router *gin.Engine) {

	// Define the route for the handlers
	router.GET("/", app.login)
}
