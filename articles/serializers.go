package articles

import "github.com/fahrurben/realworld-gin/common"

type ArticleSerializer struct {
	model *ArticleModel
}

type ArticleResponse struct {
	Title       string   `json:"title"`
	Slug        string   `json:"slug"`
	Description string   `json:"description"`
	Body        string   `json:"body"`
	AuthorID    uint     `json:"author_id"`
	Tags        []string `json:"tagList"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

func (self *ArticleSerializer) Response() ArticleResponse {
	articleResponse := ArticleResponse{
		Title:       self.model.Title,
		Slug:        self.model.Slug,
		Description: self.model.Description,
		Body:        self.model.Body,
		AuthorID:    self.model.AuthorID,
		CreatedAt:   self.model.CreatedAt.Format(common.DATE_FORMAT),
		UpdatedAt:   self.model.UpdatedAt.Format(common.DATE_FORMAT),
	}

	tags := []string{}

	for _, tagModel := range self.model.Tags {
		tags = append(tags, tagModel.Tag)
	}

	articleResponse.Tags = tags

	return articleResponse
}
