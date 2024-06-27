package database

import (
	"errors"
	"time"

	"github.com/erwaen/Chirpy/types"
)

var ErrTokenExpired = errors.New("Token Expired")

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

func (db *DB) GetRefreshTokenStruct(refreshToken string) (types.RefreshToken, error) {
	dat, err := db.loadDB()
	if err != nil {
		return types.RefreshToken{}, err
	}

	rf, exists := dat.RefreshTokens[refreshToken]
	if !exists {
		return types.RefreshToken{}, ErrNotExist
	}
	return rf, nil
}

func (db *DB) RevokeRefreshToken(refreshToken string) (types.RefreshToken, error) {
	dat, err := db.loadDB()
	if err != nil {
		return types.RefreshToken{}, err
	}

	deleteElement, exists := dat.RefreshTokens[refreshToken]
	if !exists {
		return types.RefreshToken{}, ErrNotExist
	}

	delete(dat.RefreshTokens, refreshToken)

	err = db.writeDB(dat)
	if err != nil {
		return types.RefreshToken{}, err
	}

	return deleteElement, err
}
