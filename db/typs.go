package db

import "gorm.io/gorm"

type User struct {
	gorm.Model
	PhoneNumber string
	Uuid        string
	Token       int64
	Images      []Image
	Messages    []Message
}

type Image struct {
	gorm.Model
	URL      string
	UserID   uint
	User     User
	Messages []Message
}

type Message struct {
	gorm.Model
	Text    string
	UserID  uint
	User    User
	ImageID uint
	Image   Image
}
