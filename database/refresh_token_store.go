package database

import (
	"time"

	"github.com/erwaen/Chirpy/types"
)

func (db *DB) InsertRefreshToken(userID int, refreshToken string, expiresAt time.Duration) (types.RefreshToken, error) {

	expireTime := time.Now().Add(expiresAt)
	newRefreshTokenStruct := types.RefreshToken{
		RefreshToken: refreshToken,
		UserID:       userID,
		ExpireAt:     expireTime,
	}

	dat, err := db.loadDB()
	if err != nil {
		return types.RefreshToken{}, err
	}
	dat.RefreshTokens[refreshToken] = newRefreshTokenStruct
	err = db.writeDB(dat)
	if err != nil {
		return types.RefreshToken{}, err
	}
	return newRefreshTokenStruct, nil
}
