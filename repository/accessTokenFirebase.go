package repository

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
)

type AccessTokenFirebaseModel struct {
	Client *firestore.Client
}

func Init(key string, projectID string) AccessTokenFirebaseModel {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		panic(err)
	}

	if key == "" {
		panic("Encryption key is required")
	}

	encryptionKey = []byte(key)

	fmt.Println("Firestore client initialized")

	return AccessTokenFirebaseModel{client}
}

func (at AccessTokenFirebaseModel) StoreAccessToken(userId string, token string) {
	ctx := context.Background()
	encryptedToken, err := encrypt(token)
	if err != nil {
		log.Fatal(err)
	}

	_, _, err = at.Client.Collection("access_tokens").Add(ctx, map[string]interface{}{
		"user_id":    userId,
		"token":      encryptedToken,
		"created_at": firestore.ServerTimestamp,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (at AccessTokenFirebaseModel) GetAccessToken(userId string) (string, error) {
	ctx := context.Background()
	iter := at.Client.Collection("access_tokens").Where("user_id", "==", userId).OrderBy("created_at", firestore.Desc).Limit(1).Documents(ctx)
	doc, err := iter.Next()
	if err != nil {
		return "", err
	}

	encryptedToken, err := doc.DataAt("token")
	if err != nil {
		return "", err
	}

	token, err := decrypt(encryptedToken.(string))
	if err != nil {
		return "", err
	}

	return token, nil
}
