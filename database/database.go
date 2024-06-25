package database

import (
	"encoding/json"
	"io/fs"
	"log"
	"os"
	"sort"
	"sync"
    "errors"
	"github.com/erwaen/Chirpy/types"
)

var ErrNotExist = errors.New("resource does not exist")

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]types.Chirp `json:"chirps"`
    Users map[int]types.User `json:"users"`
}

func newDBStructure() DBStructure {
	return DBStructure{
		Chirps: map[int]types.Chirp{},
        Users: map[int]types.User{},
	}
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	newDB := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}

	_, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		err := newDB.ensureDB()
		if err != nil {
			log.Fatalf("Error on creating new file %s", err)
			return &DB{}, err
		}
	}
	return newDB, nil
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	nDBStruct := newDBStructure()
	data, err := json.Marshal(nDBStruct)
	if err != nil {
		log.Fatalf("Error on marshal new db structure %s", err)
		return err
	}
	error := os.WriteFile(db.path, data, fs.ModePerm)
	if error != nil {
		log.Fatalf("Error on writing new file for new db structure %s", err)
		return err
	}
	return nil
}


// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	fileData, err := os.ReadFile(db.path)

	allData := DBStructure{}
	if err != nil {
		log.Fatalf("Error on reading file when loadDB %s", err)
		return allData, err
	}
	err = json.Unmarshal(fileData, &allData)
	if err != nil {
		log.Fatalf("Error on unmarshal file when loadDB %s", err)
		return allData, err
	}
	return allData, nil

}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]types.Chirp, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()
	chirps, err := db.loadDB()
	if err != nil {
		return []types.Chirp{}, err
	}

	var chirpList []types.Chirp
	for _, chirp := range chirps.Chirps {
		chirpList = append(chirpList, chirp)
	}
	sort.Slice(chirpList, func(i, j int) bool {
		return chirpList[i].Id < chirpList[j].Id
	})

	return chirpList, nil
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (types.Chirp, error) {

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
		Id:   newID,
		Body: body,
	}

	chirps.Chirps[newID] = newChirp
    err = db.writeDB(chirps)
	if err != nil {
		return types.Chirp{}, err
	}
	return newChirp, nil

}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	// save to file again
	data, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}
	err = os.WriteFile(db.path, data, fs.ModePerm)
	if err != nil {
		return err
	}
	return nil 
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
