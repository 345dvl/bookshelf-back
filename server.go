package main

import (
	"context"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/labstack/echo"
	"google.golang.org/api/option"
)

// 一旦mainパッケージだけで実装。後でパッケージ分ける

type User struct {
	Ñame  string
	Email string
}

func FirebaseInit() (app *firebase.App, err error) {
	ctx := context.Background()
	opt := option.WithCredentialsFile("bookshelf-back-firebase-adminsdk.json")

	app, err = firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("error initializing firebase app: %v\n", err)
	}
	return
}

func getUser(c echo.Context) error {
	user := User{"hato", "k.f.kntn@gmail.com"}
	return c.JSON(http.StatusOK, user)
}

func createUser(c echo.Context) error {
	ctx := context.Background()
	app, err := FirebaseInit()
	if err != nil {
		log.Fatalf("error initializing firebase app: %v\n", err)
	}
	client, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error initializing firebase client: %v\n", err)
	}

	params := (&auth.UserToCreate{}).
		Email("a2@a.com").
		EmailVerified(false).
		DisplayName("John Doe").
		Password("passoword!").
		Disabled(false)
	u, err := client.CreateUser(ctx, params)
	if err != nil {
		log.Fatalf("error creating user: %v\n", err)
	}
	log.Printf("Successfully created user: %#v\n", u.UserInfo)

	return c.JSON(http.StatusOK, u)
}

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Air on Docker!!")
	})
	e.GET("/users", getUser)
	e.POST("/users", createUser)
	e.Logger.Fatal(e.Start(":8080"))
}
