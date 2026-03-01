package main

import (
	"log"

	"github.com/fahrurben/realworld-gin/articles"
	"github.com/fahrurben/realworld-gin/common"
	"github.com/fahrurben/realworld-gin/users"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&users.UserModel{})
	db.AutoMigrate(&articles.TagModel{})
	db.AutoMigrate(&articles.ArticleModel{})
}

func main() {
	db := common.Init()
	Migrate(db)
	sqlDB, err := db.DB()
	if err != nil {
		log.Println("failed to get sql.DB:", err)
	} else {
		defer sqlDB.Close()
	}

	router := gin.Default()

	router.RedirectTrailingSlash = false

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	v1 := router.Group("/api")

	v1.Use(users.AuthMiddleware(true))

	users.UserRegister(v1.Group("/users"))
	articles.ArticleRegister(v1.Group("/articles"))

	router.Run(":8000")
}
