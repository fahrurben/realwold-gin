package users

type LoginValidator struct {
	User struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" bidning:"required"`
	} `json:"user"`
}

type RegisterValidator struct {
	User struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" bidning:"required"`
	} `json:"user"`
}

type UpdateValidator struct {
	User struct {
		Email    string `json:"email" binding:"omitempty,email"`
		Username string `json:"username"`
		Password string `json:"password"`
		Bio      string `json:"bio"`
		Image    string `json:"image"`
	}
}
