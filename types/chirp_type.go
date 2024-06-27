package types

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
    AuthorID int `json:"author_id"`
}
