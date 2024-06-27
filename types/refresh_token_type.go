package types

import "time"

type RefreshToken struct {
	ExpireAt     time.Time `json:"expires_at"`
	RefreshToken string `json:"refresh_token"`
	UserID       int    `json:"user_id"`
}
