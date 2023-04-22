package typs

import (
	"gorm.io/gorm"
	"time"
)

type Model struct {
	ID        string `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type User struct {
	Model
	PhoneNumber string `json:"-"`
	Name        string
	Password    string `json:"-"`
	Token       int64  `json:"-"`
	Images      []Image
	Messages    []Message
}

type Image struct {
	Model
	URL         string
	ImageDetail string
	UserID      string `json:"-"`
	User        User
	Messages    []Message
}

type Message struct {
	Model
	Text    string
	User    User `gorm:"foreignKey:UserID"`
	UserID  string
	Dst     User `gorm:"foreignKey:DstID"`
	DstID   string
	ImageID string
	Image   Image
}
