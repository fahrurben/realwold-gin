package articles

import (
	"github.com/fahrurben/realworld-gin/common"
	"github.com/fahrurben/realworld-gin/users"
	"gorm.io/gorm"
)

type ArticleModel struct {
	gorm.Model
	Title          string
	Slug           string `gorm:"uniqueIndex`
	Description    string `gorm:"size:2048"`
	Body           string `gorm:"size:2048"`
	AuthorID       uint
	Author         users.UserModel
	Tags           []TagModel      `gorm:"many2many:article_tags"`
	FavoriteModels []FavoriteModel `gorm:"ForeignKey:ArticleModelID"`
}

type TagModel struct {
	gorm.Model
	Tag           string         `gorm:"uniqueIndex"`
	ArticleModels []ArticleModel `gorm:"many2many:article_tags"`
}

type FavoriteModel struct {
	gorm.Model
	UserModelID    uint
	UserModel      users.UserModel
	ArticleModelID uint
	ArticleModel   ArticleModel
}

func SaveOne(data interface{}) error {
	db := common.GetDB()
	error := db.Save(data).Error
	return error
}

func FindOne(condition interface{}) (*ArticleModel, error) {
	db := common.GetDB()
	var model ArticleModel
	result := db.Preload("Author").Preload("Tags").Where(condition).First(&model)

	return &model, result.Error
}

func Delete(condition interface{}) error {
	db := common.GetDB()
	result := db.Where(condition).Delete(&ArticleModel{})

	return result.Error
}

func (article ArticleModel) favoritesCount() uint {
	db := common.GetDB()
	var count int64
	db.Model(&FavoriteModel{}).Where(FavoriteModel{
		ArticleModelID: article.ID,
	}).Count(&count)
	return uint(count)
}

func (article ArticleModel) isFavoritedBy(user users.UserModel) bool {
	db := common.GetDB()
	var favoriteModel FavoriteModel
	db.Where(FavoriteModel{
		UserModelID:    user.ID,
		ArticleModelID: article.ID,
	}).First(&favoriteModel)
	return favoriteModel.ID != 0
}
