package main

import (
	"net/http"

	"github.com/labstack/echo"
	controller "github.com/345dvl/bookshelf-back/controller"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Air on Docker!!")
	})
	// e.GET("/user", controller.getUser)
	controller.getUser()
	e.Logger.Fatal(e.Start(":8080"))
}
