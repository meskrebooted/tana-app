package models

import "time"

type Space struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:160;not null" json:"name"`
	Type      string    `gorm:"size:30;not null" json:"type"` // desk | room | booth
	Capacity  int       `gorm:"not null;default:1" json:"capacity"`
	Location  string    `gorm:"size:160" json:"location"`
	Vibe      string    `gorm:"size:200" json:"vibe"` // breve descrizione giocosa dello spazio
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
