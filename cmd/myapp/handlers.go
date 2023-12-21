package main

import (
	"github.com/gin-gonic/gin"
	"html/template"
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
