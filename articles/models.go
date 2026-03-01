package articles

import (
	"github.com/fahrurben/realworld-gin/common"
	"github.com/fahrurben/realworld-gin/users"
	"gorm.io/gorm"
)

type ArticleModel struct {
	gorm.Model
	Title       string
	Slug        string `gorm:"uniqueIndex`
	Description string `gorm:"size:2048"`
	Body        string `gorm:"size:2048"`
	AuthorID    uint
	Author      users.UserModel
	Tags        []TagModel `gorm:"many2many:article_tags"`
}

type TagModel struct {
	gorm.Model
	Tag           string         `gorm:"uniqueIndex"`
	ArticleModels []ArticleModel `gorm:"many2many:article_tags"`
}

func SaveOne(data interface{}) error {
	db := common.GetDB()
	error := db.Save(data).Error
	return error
}

func FindOne(condition interface{}) (*ArticleModel, error) {
	db := common.GetDB()
	var model ArticleModel
	result := db.Where(condition).First(&model)

	return &model, result.Error
}
