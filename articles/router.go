package articles

import (
	"net/http"

	"github.com/fahrurben/realworld-gin/common"
	"github.com/fahrurben/realworld-gin/users"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
)

func ArticleRegister(router *gin.RouterGroup) {
	router.POST("", ArticleSave)
}

func ArticleSave(c *gin.Context) {
	articleValidator := NewArticleValidator()

	if err := c.ShouldBindJSON(&articleValidator); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	model := ArticleModel{}
	model.Title = articleValidator.Article.Title
	model.Slug = slug.Make(model.Title)
	model.Description = articleValidator.Article.Description
	model.Body = articleValidator.Article.Body
	model.Author = c.MustGet("my_user_model").(users.UserModel)

	tagModels := []TagModel{}

	for _, tagStr := range articleValidator.Article.Tags {
		var tagModel TagModel
		db := common.GetDB()
		result := db.Where("tag = ?", tagStr).FirstOrInit(&tagModel, TagModel{Tag: tagStr})

		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		}

		tagModels = append(tagModels, tagModel)
	}

	model.Tags = tagModels

	if err := SaveOne(&model); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	serializer := ArticleSerializer{model: &model}
	c.JSON(http.StatusCreated, gin.H{"article": serializer.Response()})
}
