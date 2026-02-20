package models

type AccessTokenClaims struct {
	UserID    uint64 `json:"user_id"`
	UserEmail string `json:"user_email"`
}
