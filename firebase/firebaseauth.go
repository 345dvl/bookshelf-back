package firebaseauth

import (
	"context"
	"log"
	"firebase.google.com/go/v4"
)

func Init() (app *firebase.App, err error) {
	app, err = firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	return
}