package server

import (
   "net/http"

   "github.com/gin-gonic/gin"
   "gorm.io/gorm"
   "codex-pipeline/internal/post"
)

// SetupRouter initializes routes for posts
func SetupRouter(db *gorm.DB) *gin.Engine {
   r := gin.Default()

   r.GET("/posts", func(c *gin.Context) {
       var posts []post.Post
       if err := db.Order("created_at desc").Find(&posts).Error; err != nil {
           c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
           return
       }
       c.JSON(http.StatusOK, posts)
   })

   r.GET("/posts/:id", func(c *gin.Context) {
       var p post.Post
       if err := db.First(&p, c.Param("id")).Error; err != nil {
           c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
           return
       }
       c.JSON(http.StatusOK, p)
   })

   r.POST("/posts", func(c *gin.Context) {
       var input post.Post
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

   r.PUT("/posts/:id", func(c *gin.Context) {
       var p post.Post
       if err := db.First(&p, c.Param("id")).Error; err != nil {
           c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
           return
       }
       if err := c.ShouldBindJSON(&p); err != nil {
           c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
           return
       }
       if err := db.Save(&p).Error; err != nil {
           c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
           return
       }
       c.JSON(http.StatusOK, p)
   })

   r.DELETE("/posts/:id", func(c *gin.Context) {
       if err := db.Delete(&post.Post{}, c.Param("id")).Error; err != nil {
           c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
           return
       }
       c.JSON(http.StatusOK, gin.H{"message": "post deleted"})
   })

   return r
}
