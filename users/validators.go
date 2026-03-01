package users

type LoginValidator struct {
	User struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" bidning:"required"`
	} `json:"user"`
}
