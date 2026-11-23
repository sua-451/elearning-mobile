package models

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"size:50;unique;not null" json:"username"`
	Password  string    `gorm:"size:255;not null" json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Material struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"size:200;not null" json:"title"`
	Slug        string    `gorm:"size:200;unique;not null" json:"slug"`
	Description string    `gorm:"type:text" json:"description"`
	OrderIndex  int       `gorm:"default:0" json:"order_index"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Lessons     []Lesson  `gorm:"foreignKey:MaterialID;constraint:OnDelete:CASCADE" json:"lessons,omitempty"`
}

type Lesson struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	MaterialID    uint      `gorm:"not null;index" json:"material_id"`
	Title         string    `gorm:"size:200;not null" json:"title"`
	Slug          string    `gorm:"size:200;not null" json:"slug"`
	OrderIndex    int       `gorm:"default:0" json:"order_index"`
	ImageURL      string    `gorm:"size:500" json:"image_url"`
	Content       string    `gorm:"type:longtext" json:"content"`
	Notes         string    `gorm:"type:text" json:"notes"`
	StarterCode   string    `gorm:"type:text" json:"starter_code"`
	EnableDartpad bool      `gorm:"default:false" json:"enable_dartpad"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
