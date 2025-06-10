package main

import "time"

// Post represents a bulletin board post.
type Post struct {
    ID        uint      `gorm:"primaryKey"`
    Title     string    `json:"title"`
    Content   string    `json:"content"`
    CreatedAt time.Time
}
