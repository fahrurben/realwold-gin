package articles

import (
	"net/http"
	"strconv"

	"github.com/fahrurben/realworld-gin/common"
	"github.com/fahrurben/realworld-gin/users"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
)

func ArticleRegister(router *gin.RouterGroup) {
	router.POST("", ArticleSave)
	router.PUT("/:slug", ArticleUpdate)
	router.DELETE("/:slug", ArticleDelete)
	router.POST("/:slug/favorite", FavoriteArticle)
	router.DELETE("/:slug/favorite", UnfavoriteArticle)
	router.POST("/:slug/comments", AddComment)
}

func PublicRegister(router *gin.RouterGroup) {
	router.GET("/:slug", ArticleGet)
	router.GET("", ArticleList)
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

	serializer := ArticleSerializer{Model: &model, C: c}
	c.JSON(http.StatusCreated, gin.H{"article": serializer.Response()})
}

func ArticleGet(c *gin.Context) {
	slug := c.Param("slug")

	articleModel, err := FindOne(&ArticleModel{Slug: slug})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Article not found"})
	}

	serializer := ArticleSerializer{Model: articleModel, C: c}
	c.JSON(http.StatusOK, gin.H{"article": serializer.Response()})
}

func ArticleUpdate(c *gin.Context) {
	articleSlug := c.Param("slug")
	articleValidator := NewArticleValidator()

	if err := c.ShouldBindJSON(&articleValidator); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var model *ArticleModel

	model, err := FindOne(&ArticleModel{Slug: articleSlug})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Article not found"})
		return
	}

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

	serializer := ArticleSerializer{Model: model, C: c}
	c.JSON(http.StatusOK, gin.H{"article": serializer.Response()})
}

func ArticleDelete(c *gin.Context) {
	slug := c.Param("slug")

	_, err := FindOne(&ArticleModel{Slug: slug})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Article not found"})
		return
	}

	if err := Delete(&ArticleModel{Slug: slug}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func ArticleList(c *gin.Context) {
	tag := c.Query("tag")
	author := c.Query("author")
	favorited := c.Query("favorited")
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))

	articleModels, count := List(tag, author, favorited, uint(offset), uint(limit))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Article not found"})
	}

	serializer := ArticlesSerializer{Articles: articleModels, C: c}
	c.JSON(http.StatusOK, gin.H{"articles": serializer.Response(), "articlesCount": count})
}

func FavoriteArticle(c *gin.Context) {
	slug := c.Param("slug")

	articleModel, err := FindOne(&ArticleModel{Slug: slug})
	user := c.MustGet("my_user_model").(users.UserModel)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Article not found"})
	}

	articleModel.setFavoriteBy(user)

	serializer := ArticleSerializer{Model: articleModel, C: c}
	c.JSON(http.StatusOK, gin.H{"article": serializer.Response()})
}

func UnfavoriteArticle(c *gin.Context) {
	slug := c.Param("slug")

	articleModel, err := FindOne(&ArticleModel{Slug: slug})
	user := c.MustGet("my_user_model").(users.UserModel)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Article not found"})
	}

	articleModel.unFavoriteBy(user)

	serializer := ArticleSerializer{Model: articleModel, C: c}
	c.JSON(http.StatusOK, gin.H{"article": serializer.Response()})
}

func AddComment(c *gin.Context) {
	slug := c.Param("slug")

	user := c.MustGet("my_user_model").(users.UserModel)
	articleModel, err := FindOne(&ArticleModel{Slug: slug})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Article not found"})
	}

	commentValidator := NewCommentValidator()

	if err := c.ShouldBindJSON(&commentValidator); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	model := Comment{}
	model.Body = commentValidator.Comment.Body
	model.ArticleModel = *articleModel
	model.AuthorModel = user

	err = SaveOne(&model)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot save comment"})
	}

	serializer := CommentSerializer{Model: &model, C: c}
	c.JSON(http.StatusOK, gin.H{"article": serializer.Response()})
}
