package db

import "gorm.io/gorm"

type User struct {
	gorm.Model
	PhoneNumber string
	Name        string
	Uuid        string
	Password    string `json:"password"`
	Token       int64
	Images      []Image
	Messages    []Message
}

type Image struct {
	gorm.Model
	URL         string
	Uuid        string
	ImageDetail string `json:"imageDetail"`
	UserID      uint
	User        User
	Messages    []Message
}

type Message struct {
	gorm.Model
	Text    string
	UserID  uint
	User    User
	ImageID uint
	Image   Image
}
