package model

type User struct {
	UserName string `json:"user_name"`
	Bio      string `json:"bio"`
	FullName string `json:"full_name"`
	UserType string `json:"user_type"`
}

type Request struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
