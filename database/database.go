package database

import (
	"encoding/json"
	"errors"
	"github.com/erwaen/Chirpy/types"
	"os"
	"sort"
	"sync"
)

var ErrNotExist = errors.New("resource does not exist")

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps        map[int]types.Chirp           `json:"chirps"`
	Users         map[int]types.User            `json:"users"`
	RefreshTokens map[string]types.RefreshToken `json:"refresh_tokens"`
}

func (db *DB) createDB() error {
	dbStructure := DBStructure{
		Chirps:        map[int]types.Chirp{},
		Users:         map[int]types.User{},
		RefreshTokens: map[string]types.RefreshToken{},
	}
	return db.writeDB(dbStructure)
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}
	err := db.ensureDB()
	return db, err
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return db.createDB()
	}
	return err
}

func (db *DB) ResetDB() error {
	err := os.Remove(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return db.ensureDB()
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbStructure := DBStructure{}
	dat, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return dbStructure, err
	}
	err = json.Unmarshal(dat, &dbStructure)
	if err != nil {
		return dbStructure, err
	}
	return dbStructure, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	dat, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}
	err = os.WriteFile(db.path, dat, 0600)
	if err != nil {
		return err
	}
	return nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps(authorID int, sortBy string) ([]types.Chirp, error) {
	chirps, err := db.loadDB()
	if err != nil {
		return []types.Chirp{}, err
	}

	var chirpList []types.Chirp

	for _, chirp := range chirps.Chirps {
		if authorID == 0 || chirp.AuthorID == authorID {
			chirpList = append(chirpList, chirp)
		}
	}
	// Set default sort order if invalid
	if sortBy != "asc" && sortBy != "desc" {
		sortBy = "asc"
	}
	sort.Slice(chirpList, func(i, j int) bool {
		if sortBy == "asc" {
			return chirpList[i].Id < chirpList[j].Id
		}
		return chirpList[i].Id > chirpList[j].Id
	})

	return chirpList, nil
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string, authorID int) (types.Chirp, error) {

	chirps, err := db.loadDB()
	if err != nil {
		return types.Chirp{}, err
	}

	newID := 0
	for id := range chirps.Chirps {
		if id > newID {
			newID = id
		}
	}
	newID++
	newChirp := types.Chirp{
		Id:       newID,
		Body:     body,
		AuthorID: authorID,
	}

	chirps.Chirps[newID] = newChirp
	err = db.writeDB(chirps)
	if err != nil {
		return types.Chirp{}, err
	}
	return newChirp, nil

}

func (db *DB) GetChirp(id int) (types.Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return types.Chirp{}, err
	}

	chirp, ok := dbStructure.Chirps[id]
	if !ok {
		return types.Chirp{}, ErrNotExist
	}

	return chirp, nil
}

func (db *DB) DeleteChirp(id int) (types.Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return types.Chirp{}, err
	}

	deletedChirp, ok := dbStructure.Chirps[id]
	if !ok {
		return types.Chirp{}, ErrNotExist
	}

	delete(dbStructure.Chirps, id)
	err = db.writeDB(dbStructure)
	if err != nil {
		return types.Chirp{}, err
	}

	return deletedChirp, nil
}
