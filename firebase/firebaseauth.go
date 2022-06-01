package firebaseauth

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

func Init() (app *firebase.App, err error) {
	ctx := context.Background()
	opt := option.WithCredentialsFile("bookshelf-back-firebase-adminsdk.json")

	app, err = firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("error initializing firebase app: %v\n", err)
	}
	return
}
