package database

import (
	"errors"
	"github.com/erwaen/Chirpy/types"
)

var ErrUserAlreadyExist = errors.New("User already exists")

// CreateUser creates a new user and saves it to disk
func (db *DB) CreateUser(email string, password string) (types.User, error) {
	if _, err := db.GetUserByEmail(email); !errors.Is(err, ErrNotExist) {
		return types.User{}, ErrUserAlreadyExist
	}

	allData, err := db.loadDB()
	if err != nil {
		return types.User{}, err
	}

	newID := 0
	for id := range allData.Users {
		if id > newID {
			newID = id
		}
	}
	newID++

	newUser := types.User{
		Id:       newID,
		Email:    email,
		Password: password,
	}
	allData.Users[newID] = newUser

	err = db.writeDB(allData)
	if err != nil {
		return types.User{}, err
	}

	return newUser, nil
}

func (db *DB) GetUserByEmail(email string) (types.User, error) {
	allDB, err := db.loadDB()
	if err != nil {
		return types.User{}, err
	}
	for _, user := range allDB.Users {
		if user.Email == email {
			return user, nil
		}
	}
	return types.User{}, ErrNotExist
}

func (db *DB) GetUserByID(id int) (types.User, error) {
	allDB, err := db.loadDB()
	if err != nil {
		return types.User{}, err
	}
	user, ok := allDB.Users[id]
	if !ok {
		return types.User{}, ErrNotExist
	}
	return user, nil
}

func (db *DB) UpdateUser(id int, email, hashedPassword string) (types.User, error) {
	allDB, err := db.loadDB()
	if err != nil {
		return types.User{}, err
	}
	user, ok := allDB.Users[id]
	if !ok {
		return types.User{}, ErrNotExist
	}

	user.Email = email
	user.Password = hashedPassword
	allDB.Users[user.Id] = user
	err = db.writeDB(allDB)
	if err != nil {
		return types.User{}, err
	}
	return user, nil
}
