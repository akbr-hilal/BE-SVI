package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Post struct {
	Id          int       `json:"id" gorm:"primary_key"`
	Title       string    `json:"title" validate:"required,min=20"`
	Content     string    `json:"content" validate:"required,min=200"`
	Category    string    `json:"category" validate:"required,min=3"`
	CreatedDate time.Time `json:"created_date"`
	UpdatedDate time.Time `json:"updated_date"`
	Status      string    `json:"status" validate:"required,oneof=publish draft trash"`
}

var (
	db       *gorm.DB
	err      error
	validate *validator.Validate
)

func main() {
	dsn := "root:@tcp(localhost:3306)/article?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open("mysql", dsn)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	db.AutoMigrate(&Post{})

	validate = validator.New()

	r := gin.Default()

	// Set CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://localhost:5173", "*"} // Allow React app origin
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}                        // Allow specific methods
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	r.POST("/article", createArticle)
	r.GET("/article", getAllArticles)
	r.GET("/article/:id", getArticle)
	r.PUT("/article/:id", updateArticle)
	r.DELETE("/article/:id", deleteArticle)

	r.Run()
}

func createArticle(c *gin.Context) {
	var post Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validate.Struct(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validation_error": err.Error()})
		return
	}

	post.CreatedDate = time.Now()
	post.UpdatedDate = time.Now()
	db.Create(&post)
	c.JSON(http.StatusOK, gin.H{"message": "Article created successfully"})
}

func getAllArticles(c *gin.Context) {
	var posts []Post
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 10
	}
	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		offset = 0
	}
	status := c.Query("status")

	var total int
	query := db.Model(&Post{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get articles with limit and offset
	query = db.Limit(limit).Offset(offset)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if err := query.Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total": total,
		"data":  posts,
	})
}

func getArticle(c *gin.Context) {
	var post Post
	id := c.Param("id")
	db.First(&post, id)
	if post.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Article not found"})
		return
	}
	c.JSON(http.StatusOK, post)
}

func updateArticle(c *gin.Context) {
	var post Post
	id := c.Param("id")
	db.First(&post, id)
	if post.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Article not found"})
		return
	}

	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validate.Struct(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validation_error": err.Error()})
		return
	}

	post.UpdatedDate = time.Now()
	db.Save(&post)
	c.JSON(http.StatusOK, gin.H{"message": "Article updated successfully"})
}

func deleteArticle(c *gin.Context) {
	var post Post
	id := c.Param("id")
	db.First(&post, id)
	if post.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Article not found"})
		return
	}
	db.Delete(&post)
	c.JSON(http.StatusOK, gin.H{"message": "Article deleted successfully"})
}
