package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/joho/godotenv"
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

type UserLoginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type IDToken struct {
	IDToken string `json:"idToken"`
}

const envFile string = "./application.env"

func GetEnv(key string) (value string) {
	if env := os.Getenv(key); env != "" {
		value = env
	} else {
		loadFile := godotenv.Load(envFile)
		env, err := godotenv.Read(envFile)

		if (loadFile != nil) || (err != nil) {
			var message = "Error loading .env file"
			log.Print(message)
			panic(message)
		}

		value = env[key]
	}

	return
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

// func createCustomTokenByUID(uid string) (customToken string){
// 	ctx := context.Background()
// 	app, err := FirebaseInit()
// 	if err != nil {
// 		log.Fatalf("error initializing firebase app: %v\n", err)
// 	}
// 	client, err := app.Auth(ctx)
// 	if err != nil {
// 		log.Fatalf("error initializing firebase client: %v\n", err)
// 	}

// 	customToken, err = client.CustomToken(ctx, uid)
// 	if err != nil {
// 		log.Fatalf("error minting custom token: %v\n", err)
// 	}

// 	return
// }

func genIDToken(c echo.Context) error{
	params := new(UserLoginParams)
	if err := c.Bind(params); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	apiKey := GetEnv("FIREBASE_WEB_API_KEY")
	uri := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=%s", apiKey)

	requestBody := &UserLoginParams{Email: params.Email, Password: params.Password}
	jsonStr, err := json.Marshal(requestBody)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	req, err := http.NewRequest(
		"POST",
		uri,
		bytes.NewBuffer(jsonStr),
	)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

  req.Header.Set("Content-Type", "application/json")
	httpClient := &http.Client{}
	response, err := httpClient.Do(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	defer response.Body.Close()

	var idToken IDToken
	body, _ := io.ReadAll(response.Body)
	if err := json.Unmarshal(body, &idToken); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, idToken)
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
	e.POST("/users/id_token", genIDToken)
	e.Logger.Fatal(e.Start(":8080"))
}
