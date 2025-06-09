package model

type Notification struct {
	ID      string
	UserID  string
	Message string
	Read    bool
}