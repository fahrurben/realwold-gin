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
	ID             uint                `json:"id"`
	Title          string              `json:"title"`
	Slug           string              `json:"slug"`
	Description    string              `json:"description"`
	Body           string              `json:"body"`
	Author         ArticleUserResponse `json:"author"`
	Tags           []string            `json:"tagList"`
	Favorited      bool                `json:"favorited"`
	FavoritesCount uint                `json:"favoritesCount"`
	CreatedAt      string              `json:"createdAt"`
	UpdatedAt      string              `json:"updatedAt"`
}

type ArticlesSerializer struct {
	C        *gin.Context
	Articles []*ArticleModel
}

func (self *ArticleSerializer) Response() ArticleResponse {
	myUserModel := self.C.MustGet("my_user_model").(users.UserModel)

	articleUserSerializer := ArticleUserSerializer{Model: &self.Model.Author}
	articleResponse := ArticleResponse{
		ID:             self.Model.ID,
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

func (self *ArticlesSerializer) Response() []ArticleResponse {
	myUserModel := self.C.MustGet("my_user_model").(users.UserModel)
	results := []ArticleResponse{}

	for _, articleModel := range self.Articles {
		articleUserSerializer := ArticleUserSerializer{Model: &articleModel.Author}
		articleResponse := ArticleResponse{
			Title:          articleModel.Title,
			Slug:           articleModel.Slug,
			Description:    articleModel.Description,
			Body:           articleModel.Body,
			Author:         articleUserSerializer.Response(),
			Favorited:      articleModel.isFavoritedBy(myUserModel),
			FavoritesCount: articleModel.favoritesCount(),
			CreatedAt:      articleModel.CreatedAt.Format(common.DATE_FORMAT),
			UpdatedAt:      articleModel.UpdatedAt.Format(common.DATE_FORMAT),
		}

		tags := []string{}

		for _, tagModel := range articleModel.Tags {
			tags = append(tags, tagModel.Tag)
		}

		articleResponse.Tags = tags

		results = append(results, articleResponse)
	}

	return results
}

type CommentSerializer struct {
	C     *gin.Context
	Model *Comment
}

type CommentsSerializer struct {
	C        *gin.Context
	Comments []*Comment
}

type CommentResponse struct {
	ID        uint                `json:"id"`
	Body      string              `json:"body"`
	Author    ArticleUserResponse `json:"author"`
	CreatedAt string              `json:"createdAt"`
	UpdatedAt string              `json:"updatedAt"`
}

func (self *CommentSerializer) Response() CommentResponse {

	authorSerializer := ArticleUserSerializer{Model: &self.Model.AuthorModel}
	commentResponse := CommentResponse{
		ID:        self.Model.ID,
		Body:      self.Model.Body,
		Author:    authorSerializer.Response(),
		CreatedAt: self.Model.CreatedAt.Format(common.DATE_FORMAT),
		UpdatedAt: self.Model.UpdatedAt.Format(common.DATE_FORMAT),
	}

	return commentResponse
}

func (self *CommentsSerializer) Response() []CommentResponse {
	commentResponses := []CommentResponse{}

	for _, comment := range self.Comments {
		authorSerializer := ArticleUserSerializer{Model: &comment.AuthorModel}
		commentResponse := CommentResponse{
			ID:        comment.ID,
			Body:      comment.Body,
			Author:    authorSerializer.Response(),
			CreatedAt: comment.CreatedAt.Format(common.DATE_FORMAT),
			UpdatedAt: comment.UpdatedAt.Format(common.DATE_FORMAT),
		}

		commentResponses = append(commentResponses, commentResponse)
	}

	return commentResponses
}
