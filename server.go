package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/labstack/echo"
	"google.golang.org/api/option"
)

// 一旦mainパッケージだけで実装。後でパッケージ分ける

type UserRegistrationParams struct {
	Name                 string `json:"name"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
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
	ctx := context.Background()
	app, err := FirebaseInit()
	if err != nil {
		log.Fatalf("error initializing firebase app: %v\n", err)
	}
	client, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error initializing firebase client: %v\n", err)
	}

	uid := "03bLtdlTyeT7V1NUQTlABbkhLRB3"
	u, err := client.GetUser(ctx, uid)
	if err != nil {
			log.Fatalf("error getting user %s: %v\n", uid, err)
	}
	log.Printf("Successfully fetched user data: %#v\n", u.UserInfo)

	return c.JSON(http.StatusOK, u)
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

	params := new(UserRegistrationParams)
	if err := c.Bind(params); err != nil {
		return c.JSON(http.StatusBadRequest, "エラーが発生しました。不正なリクエストです。")
	}

	if params.Password != params.PasswordConfirmation {
		return c.JSON(http.StatusBadRequest, "エラーが発生しました。パスワードと確認用パスワードが一致していません。")
	}

	createParams := (&auth.UserToCreate{}).
		Email(params.Email).
		EmailVerified(false).
		DisplayName(params.Name).
		Password(params.Password).
		Disabled(false)
	u, err := client.CreateUser(ctx, createParams)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "エラーが発生しました。不正なリクエストです。")
	}
	log.Printf("Successfully created user: %#v\n", u.UserInfo)

	return c.JSON(http.StatusOK, u)
}

func updateUser(c echo.Context) error {
	ctx := context.Background()
	app, err := FirebaseInit()
	if err != nil {
		log.Fatalf("error initializing firebase app: %v\n", err)
	}
	client, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error initializing firebase client: %v\n", err)
	}

	var uid string

	params := (&auth.UserToUpdate{}).
		Email("new1@a.com").
		EmailVerified(false).
		DisplayName("Alder").
		Password("newpassoword!").
		Disabled(false)
	u, err := client.UpdateUser(ctx, uid, params)
	if err != nil {
		log.Fatalf("error updating user: %v\n", err)
	}

	log.Printf("Successfully updated user: %#v\n", u.UserInfo)

	return c.JSON(http.StatusOK, u)
}

func deleteUser(c echo.Context) error {
	var err error
	return err
}

func createCustomTokenByUID(uid string) (customToken string){
	ctx := context.Background()
	app, err := FirebaseInit()
	if err != nil {
		log.Fatalf("error initializing firebase app: %v\n", err)
	}
	client, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error initializing firebase client: %v\n", err)
	}

	customToken, err = client.CustomToken(ctx, uid)
	if err != nil {
		log.Fatalf("error minting custom token: %v\n", err)
	}

	return
}

func genIDTokenForDebug(uid string) string{
	ctx := context.Background()
	app, err := FirebaseInit()
	if err != nil {
		log.Fatalf("error initializing firebase app: %v\n", err)
	}
	client, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error initializing firebase client: %v\n", err)
	}

	customToken, err := client.CustomToken(ctx, uid)
	if err != nil {
		log.Fatalf("error creating custom token: %v\n", err)
	}
	payload := map[string]interface{}{
		"token":             token,
		"returnSecureToken": true,
	}
	req, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("error occured: %v\n", err)
	}
	idToken, err := getIDTokenFromUID(uid)
	response := getDebugIDToken(idToken)
	if err != nil {
		panic("エラーが発生しました")
	}

	return c.JSON(http.StatusOK, response)
}

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Air on Docker!!")
	})
	e.GET("/users", getUser)
	e.POST("/users", createUser)
	e.PATCH("/users", updateUser)
	e.DELETE("/users/:id", deleteUser)
	e.Logger.Fatal(e.Start(":8080"))
}
