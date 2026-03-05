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

type Comment struct {
	gorm.Model
	AuthorModelID  uint
	AuthorModel    users.UserModel
	ArticleModelID uint
	ArticleModel   ArticleModel
	Body           string `gorm:"type:text"`
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

func List(tag string, author string, favorited string, offset uint, limit uint) ([]*ArticleModel, int64) {
	db := common.GetDB()
	var articles []*ArticleModel
	var count int64
	condition := make(map[string]interface{})

	query := db.Preload("Author").Preload("Tags")

	if tag != "" {
		query = query.Joins("inner join article_tags on article_tags.article_model_id = article_models.id")
		query = query.Joins("inner join tag_models on article_tags.tag_model_id = tag_models.id")
	}

	if author != "" {
		query = query.Joins("inner join user_models on article_models.author_id = user_models.id")
	}

	if tag != "" {
		condition["tag_models.tag"] = tag
	}

	if author != "" {
		condition["user_models.username"] = author
	}

	queryData := query.Where(condition).Limit(int(limit)).Offset(int(offset))
	queryCount := query.Where(condition).Limit(int(limit)).Offset(int(offset))

	queryData.Find(&articles)
	queryCount.Count(&count)

	return articles, count
}

func (article ArticleModel) setFavoriteBy(user users.UserModel) error {
	db := common.GetDB()

	var existingModel FavoriteModel
	db.Where(FavoriteModel{
		UserModelID:    user.ID,
		ArticleModelID: article.ID,
	}).First(&existingModel)

	if existingModel.ID == 0 {
		favoriteModel := FavoriteModel{UserModel: user, ArticleModel: article}
		result := db.Save(&favoriteModel)
		return result.Error
	}

	return nil
}

func (article ArticleModel) unFavoriteBy(user users.UserModel) error {
	db := common.GetDB()

	var existingModel FavoriteModel
	db.Where(FavoriteModel{
		UserModelID:    user.ID,
		ArticleModelID: article.ID,
	}).First(&existingModel)

	if existingModel.ID != 0 {
		db.Delete(&existingModel)
	}

	return nil
}
