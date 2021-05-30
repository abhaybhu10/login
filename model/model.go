package model

type User struct {
	Email    string
	ID       string
	Password string
	Name     string
}

type Session struct {
	UserId    string
	SessionID string
}
