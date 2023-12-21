package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
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

func (app *application) registerPostRequest(c *gin.Context) {
	fmt.Println("registerPostRequest")
	// Response struct
	type User struct {
		Name            string `json:"name"`
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var userData User
	// Bind  the JSON data from the request to the userData struct
	if err := c.BindJSON(&userData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	// TODO: Handle user data, validate it, and perform action as needed
	// HACK: Temporary solution until validation is done, echoing back the received data as JSON
	c.JSON(http.StatusOK, userData)
	fmt.Println(userData)
}
