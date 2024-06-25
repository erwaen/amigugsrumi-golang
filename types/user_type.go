package types

type User struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
    Password string `json:"password"`
}

type LoginUser struct {
	Email string `json:"email"`
    Password string `json:"password"`
}

type LoggedUser struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}
