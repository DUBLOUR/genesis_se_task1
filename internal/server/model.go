package server

type User struct {
	Email        string
	PasswordHash string
	Token        string //Random string generated at registration
}
