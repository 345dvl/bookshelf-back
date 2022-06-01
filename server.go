package main

import (
	"net/http"

	"github.com/345dvl/bookshelf-back/controller"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Air on Docker!!")
	})
	// e.GET("/user", controller.getUser)
	user := controller.GetUser()
	e.Logger.Fatal(e.Start(":8080"))
}
