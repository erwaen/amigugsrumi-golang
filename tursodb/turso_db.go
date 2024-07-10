package tursodb

import (
	"database/sql"
	"fmt"

	"github.com/erwaen/Chirpy/types"
)

type TursoDB struct {
	db *sql.DB
}

func NewTursoDB(db *sql.DB) *TursoDB {
	return &TursoDB{db: db}
}

func (t *TursoDB) GetUsers() ([]types.TursoUser, error) {
	rows, err := t.db.Query("SELECT id, name FROM users")
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	var users []types.TursoUser
	for rows.Next() {
		var user types.TursoUser
		if err := rows.Scan(&user.ID, &user.Name); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %v", err)
	}

	return users, nil
}

func (t *TursoDB) GetItems() ([]types.TursoItem, error) {
	rows, err := t.db.Query("SELECT id, title, description, image_src, image_alt, price, stock FROM items")
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()


    var items []types.TursoItem
    for rows.Next() {
        var item types.TursoItem
        err := scanItem(rows, &item)
        if err != nil {
            return nil, fmt.Errorf("error scanning row: %v", err)
        }
        items = append(items, item)
    }

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %v", err)
	}

	return items, nil
}


func (t *TursoDB) GetItemsStock() ([]types.TursoItemStock, error) {
    query:= fmt.Sprintf("SELECT `id`, `stock` from items")
	rows, err := t.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()


    var items []types.TursoItemStock
    for rows.Next() {
        var item types.TursoItemStock
        err := rows.Scan(&item.ID, &item.Stock)
        if err != nil {
            return nil, fmt.Errorf("error scanning row: %v", err)
        }
        items = append(items, item)
    }

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %v", err)
	}

	return items, nil
}

func scanItem(rows *sql.Rows, item *types.TursoItem) error {
    return rows.Scan(
        &item.ID,
        &item.Title,
        &item.Description,
        &item.Image.Src,
        &item.Image.Alt,
        &item.Price,
        &item.Stock,
    )
}

func (t *TursoDB) CreateUser(name string) (int, error) {
	result, err := t.db.Exec("INSERT INTO users (name) VALUES (?)", name)
	if err != nil {
		return 0, fmt.Errorf("failed to insert user: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %v", err)
	}

	return int(id), nil
}
