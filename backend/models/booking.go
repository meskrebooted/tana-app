package models

import "time"

type Booking struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	SpaceID   uint      `gorm:"not null;index" json:"space_id"`
	StartTime time.Time `gorm:"not null" json:"start_time"`
	EndTime   time.Time `gorm:"not null" json:"end_time"`
	Status    string    `gorm:"size:20;not null;default:confirmed" json:"status"` // confirmed | cancelled
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User  User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Space Space `gorm:"foreignKey:SpaceID" json:"space,omitempty"`
}
