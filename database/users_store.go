package database

import (
	"errors"
	"github.com/erwaen/Chirpy/types"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserNotFound = errors.New("user not found")
var ErrUserAlreadyExist = errors.New("User already exists")

// CreateUser creates a new user and saves it to disk
func (db *DB) CreateUser(email string, password string) (types.LoggedUser, error) {
	_, err := db.GetUserByEmail(email)

	// Save the user to the database
	if err == nil {
		return types.LoggedUser{}, ErrUserAlreadyExist
	} else if err != ErrUserNotFound {
		return types.LoggedUser{}, err
	}

	allData, err := db.loadDB()
	if err != nil {
		return types.LoggedUser{}, err
	}

	newID := 0
	for id := range allData.Users {
		if id > newID {
			newID = id
		}
	}
	newID++

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return types.LoggedUser{}, err
	}


	newUser := types.User{
		Id:       newID,
		Email:    email,
		Password: string(hash),
	}

	allData.Users[newID] = newUser
	err = db.writeDB(allData)
	if err != nil {
		return types.LoggedUser{}, err
	}
	loggdUser := types.LoggedUser{
		Id:    newUser.Id,
		Email: newUser.Email,
	}

	return loggdUser, nil
}

func (db *DB) GetUserByEmail(email string) (types.User, error) {
	allDB, err := db.loadDB()
	if err != nil {
		return types.User{}, err
	}
	found := false
	var retUser types.User
	for _, user := range allDB.Users {
		if user.Email == email {
			found = true
			retUser = user
			break
		}
	}
	if !found {
		return retUser, ErrUserNotFound
	}
	return retUser, nil
}
