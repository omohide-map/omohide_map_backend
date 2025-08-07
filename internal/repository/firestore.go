package repository

import (
	"context"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

type FirestoreClient struct {
	*firestore.Client
}

func NewFirestoreClient(ctx context.Context, credentialsPath string) (*FirestoreClient, error) {
	opt := option.WithCredentialsFile(credentialsPath)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}

	return &FirestoreClient{client}, nil
}

func (c *FirestoreClient) Close() error {
	return c.Client.Close()
}
