package database

import (
	"time"

	"github.com/erwaen/Chirpy/types"
)


func (db *DB) InsertRefreshToken(userID, refreshToken string, expiresAt time.Duration()) (types.RefreshToken, error) {
