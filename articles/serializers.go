package articles

import (
	"github.com/fahrurben/realworld-gin/common"
	"github.com/fahrurben/realworld-gin/users"
	"github.com/gin-gonic/gin"
)

type ArticleUserSerializer struct {
	Model *users.UserModel
}

type ArticleUserResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Bio      string `json:"bio"`
	Image    string `json:"image"`
}

func (self *ArticleUserSerializer) Response() ArticleUserResponse {
	articleUser := ArticleUserResponse{
		Username: self.Model.Username,
		Email:    self.Model.Email,
		Bio:      self.Model.Bio,
		Image:    self.Model.Image,
	}

	return articleUser
}

type ArticleSerializer struct {
	C     *gin.Context
	Model *ArticleModel
}

type ArticleResponse struct {
	Title          string              `json:"title"`
	Slug           string              `json:"slug"`
	Description    string              `json:"description"`
	Body           string              `json:"body"`
	Author         ArticleUserResponse `json:"author"`
	Tags           []string            `json:"tagList"`
	Favorited      bool                `json:"favorited"`
	FavoritesCount uint                `json:"favoritesCount"`
	CreatedAt      string              `json:"created_at"`
	UpdatedAt      string              `json:"updated_at"`
}

func (self *ArticleSerializer) Response() ArticleResponse {
	myUserModel := self.C.MustGet("my_user_model").(users.UserModel)

	articleUserSerializer := ArticleUserSerializer{Model: &self.Model.Author}
	articleResponse := ArticleResponse{
		Title:          self.Model.Title,
		Slug:           self.Model.Slug,
		Description:    self.Model.Description,
		Body:           self.Model.Body,
		Author:         articleUserSerializer.Response(),
		Favorited:      self.Model.isFavoritedBy(myUserModel),
		FavoritesCount: self.Model.favoritesCount(),
		CreatedAt:      self.Model.CreatedAt.Format(common.DATE_FORMAT),
		UpdatedAt:      self.Model.UpdatedAt.Format(common.DATE_FORMAT),
	}

	tags := []string{}

	for _, tagModel := range self.Model.Tags {
		tags = append(tags, tagModel.Tag)
	}

	articleResponse.Tags = tags

	return articleResponse
}
