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
	query := `
		SELECT 
			i.id, i.title, i.description, i.image_src, i.image_alt, i.price, i.stock,
			t.id, t.url_img, t.color_background, t.tagname
		FROM 
			items i
			LEFT JOIN item_tags it ON i.id = it.item_id
			LEFT JOIN tags t ON it.tag_id = t.id
	`
	rows, err := t.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	itemsMap := make(map[int]*types.TursoItem)

	for rows.Next() {
		var item types.TursoItem
		var tagID sql.NullInt64
		var urlImg, colorBackground, tagname sql.NullString

		err := rows.Scan(
			&item.ID,
			&item.Title,
			&item.Description,
			&item.Image.Src,
			&item.Image.Alt,
			&item.Price,
			&item.Stock,
			&tagID,
			&urlImg,
			&colorBackground,
			&tagname,
		)

		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		if existingItem, exists := itemsMap[item.ID]; exists {
			if tagID.Valid {
				existingItem.Tags = append(existingItem.Tags, types.TursoTag{
					ID:              int(tagID.Int64),
					URLImg:          urlImg.String,
					ColorBackground: colorBackground.String,
					Tagname:         tagname.String,
				})
			}
		} else {
			item.Tags = []types.TursoTag{}
			if tagID.Valid {
				item.Tags = append(item.Tags, types.TursoTag{
					ID:              int(tagID.Int64),
					URLImg:          urlImg.String,
					ColorBackground: colorBackground.String,
					Tagname:         tagname.String,
				})

			}
			itemsMap[item.ID] = &item
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %v", err)
	}

	var items []types.TursoItem
	for _, item := range itemsMap {
		items = append(items, *item)
	}

	return items, nil
}

func (t *TursoDB) GetItemsStock() ([]types.TursoItemStock, error) {
	query := fmt.Sprintf("SELECT `id`, `stock` from items")
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
