package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func setupDatabase() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, errors.New("DATABASE_URL is not set")
	}
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

// setupRouter initializes Gin engine with routes using the given DB instance
func setupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.GET("/posts", func(c *gin.Context) {
		var posts []Post
		if err := db.Order("created_at desc").Find(&posts).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, posts)
	})

	r.GET("/posts/:id", func(c *gin.Context) {
		var post Post
		if err := db.First(&post, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}
		c.JSON(http.StatusOK, post)
	})

	r.POST("/posts", func(c *gin.Context) {
		var input Post
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Create(&input).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, input)
	})

	r.DELETE("/posts/:id", func(c *gin.Context) {
		var post Post
		if err := db.Delete(&post, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "post deleted"})
	})

	r.PUT("/posts/:id", func(c *gin.Context) {
		var post Post
		if err := db.First(&post, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}
		if err := c.ShouldBindJSON(&post); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Save(&post).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, post)
	})

	return r
}

// main starts the application, connecting to the database and serving HTTP routes
func main() {
	db, err := setupDatabase()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := db.AutoMigrate(&Post{}); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	r := setupRouter(db)

	log.Println("Listening on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
