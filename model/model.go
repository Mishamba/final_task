package model

import "time"

type User struct {
	ID          int
	Name        string
	Email       string
	Password    string
	Subscribers string
}

type Tweet struct {
	Tweet      string
	AuthorID   string
	UploadTime time.Time
}
