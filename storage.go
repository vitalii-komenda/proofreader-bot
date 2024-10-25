package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var encryptionKey []byte // 32 bytes

type AccessToken struct {
	Token string
}

func initDB(key string) *sql.DB {
	db, err := sql.Open("sqlite3", "./tokens.db")
	if err != nil {
		panic(err)
	}

	if key == "" {
		panic("Encryption key is required")
	}

	encryptionKey = []byte(key)

	createTableSQL := `CREATE TABLE IF NOT EXISTS access_tokens (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,	
		"user_id" TEXT,
		"created_at" DATETIME DEFAULT CURRENT_TIMESTAMP,
		"token" TEXT
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		panic(err)
	}
	fmt.Println("DB inited")

	return db
}

func encrypt(text string) (string, error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(text), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decrypt(encryptedText string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func storeAccessToken(userId string, token string) {
	encryptedToken, err := encrypt(token)
	if err != nil {
		log.Fatal(err)
	}

	insertTokenSQL := `INSERT INTO access_tokens(token, user_id) VALUES (?, ?)`
	statement, err := db.Prepare(insertTokenSQL)
	if err != nil {
		log.Fatal(err)
	}
	defer statement.Close()

	_, err = statement.Exec(encryptedToken, userId)
	if err != nil {
		log.Fatal(err)
	}
}

func getAccessToken(userId string) (string, error) {
	row := db.QueryRow(`SELECT token FROM access_tokens WHERE user_id = ? ORDER BY id DESC LIMIT 1`, userId)
	var encryptedToken string
	err := row.Scan(&encryptedToken)
	if err != nil {
		return "", err
	}

	token, err := decrypt(encryptedToken)
	if err != nil {
		return "", err
	}

	return token, nil
}
