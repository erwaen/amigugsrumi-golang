package types

type TursoItem struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       struct {
		Src string `json:"src"`
		Alt string `json:"alt"`
	} `json:"image"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
}

type TursoItemStock struct {
	ID          int    `json:"id"`
	Stock int     `json:"stock"`
}
