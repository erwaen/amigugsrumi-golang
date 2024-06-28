package types

type User struct {
	Id          int    `json:"id"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

type UpdateUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUser struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

type LoggedUser struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}
