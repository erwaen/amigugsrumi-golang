package types

import "database/sql"

type TursoItem struct {
	ID          int            `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Image       TursoItemImage `json:"image"`
	Price       float64        `json:"price"`
	Stock       int            `json:"stock"`
	Tags        []TursoTag     `json:"tags"`
	Size        TursoSize      `json:"size"`
}

type TursoItemImage struct {
	Src string `json:"src"`
	Alt string `json:"alt"`
}

type TursoItemStock struct {
	ID    int `json:"id"`
	Stock int `json:"stock"`
}

type TursoTag struct {
	ID              int    `json:"id"`
	URLImg          string `json:"url_img"`
	ColorBackground string `json:"color_background"`
	Tagname         string `json:"tagname"`
}

type TursoSize struct {
	Length sql.NullFloat64 `json:"length"`
	Width  sql.NullFloat64 `json:"width"`
	Height sql.NullFloat64 `json:"height"`
}
