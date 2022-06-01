package controller

import (
	"net/http"

	"github.com/labstack/echo"
)

type User struct {
	Ñame  string
	Email string
}

func getUser(c echo.Context) error {
	user := User{"hato", "k.f.kntn@gmail.com"}
	return c.JSON(http.StatusOK, user)
}