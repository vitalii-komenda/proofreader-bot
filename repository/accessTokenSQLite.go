package repository

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type AccessTokenSQLiteModel struct {
	DB *sql.DB
}

func InitSQLite(key string) AccessTokenSQLiteModel {
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

	return AccessTokenSQLiteModel{db}
}

func (at AccessTokenSQLiteModel) StoreAccessToken(userId string, token string) {
	encryptedToken, err := encrypt(token)
	if err != nil {
		log.Fatal(err)
	}

	insertTokenSQL := `INSERT INTO access_tokens(token, user_id) VALUES (?, ?)`
	statement, err := at.DB.Prepare(insertTokenSQL)
	if err != nil {
		log.Fatal(err)
	}
	defer statement.Close()

	_, err = statement.Exec(encryptedToken, userId)
	if err != nil {
		log.Fatal(err)
	}
}

func (at AccessTokenSQLiteModel) GetAccessToken(userId string) (string, error) {
	row := at.DB.QueryRow(`SELECT token FROM access_tokens WHERE user_id = ? ORDER BY id DESC LIMIT 1`, userId)
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
