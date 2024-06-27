package types

type RefreshToken struct {
	ExpireAt     string `json:"expires_at"`
	RefreshToken string `json:"refresh_token"`
	UserID       int    `json:"user_id"`
}
