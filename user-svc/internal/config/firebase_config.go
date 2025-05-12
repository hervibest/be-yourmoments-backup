package config

import (
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

func NewFirebaseConfig() *firebase.App {
	ctx := context.Background()
	firebaseApp, err := firebase.NewApp(ctx, nil, option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")))
	if err != nil {
		log.Fatalf("error initializing firebase app: %v", err)
	}

	return firebaseApp
}
