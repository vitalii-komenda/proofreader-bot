package repository

var encryptionKey []byte // 32 bytes

type AccessToken interface {
	StoreAccessToken(userId string, token string)
	GetAccessToken(userId string) (string, error)
}
