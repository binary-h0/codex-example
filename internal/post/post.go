package post

import "time"

// Post represents a forum post
type Post struct {
   ID        uint      `gorm:"primaryKey" json:"id"`
   Title     string    `gorm:"type:varchar(255)" json:"title"`
   Content   string    `gorm:"type:text" json:"content"`
   CreatedAt time.Time `json:"created_at"`
}
