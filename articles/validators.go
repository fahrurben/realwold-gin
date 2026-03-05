package articles

type ArticleValidator struct {
	Article struct {
		Title       string   `json:"title" binding:"required,min=4"`
		Description string   `json:"description" binding:"required,max=2048"`
		Body        string   `json:"body" binding:"required,max=2048"`
		AuthorID    uint     `json:"author_id"`
		Tags        []string `json:"tagList"`
	} `json:"article"`
}

func NewArticleValidator() ArticleValidator {
	return ArticleValidator{}
}

type CommentValidator struct {
	Comment struct {
		Body string `json:"body" binding:"required,min=4"`
	} `json:"comment"`
}

func NewCommentValidator() CommentValidator {
	return CommentValidator{}
}
