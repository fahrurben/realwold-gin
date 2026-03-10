package users

import (
	"net/http"

	"github.com/fahrurben/realworld-gin/common"

	"github.com/gin-gonic/gin"
)

func UsersRegister(router *gin.RouterGroup) {
	router.POST("/login", UsersLogin)
	router.POST("", UserRegister)
}

func UsersLogin(c *gin.Context) {
	loginValidator := LoginValidator{}
	if err := c.ShouldBindJSON(&loginValidator); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userModel, err := FindOneUser(&UserModel{Email: loginValidator.User.Email})

	if err != nil {
		c.JSON(http.StatusUnauthorized, common.NewError("Wrong email or password", err))
		return
	}

	if userModel.checkPassword(loginValidator.User.Password) != nil {
		c.JSON(http.StatusUnauthorized, common.NewError("Wrong email or password", err))
		return
	}

	UpdateContextUserModel(c, userModel.ID)
	serializer := UserSerializer{&userModel}
	c.JSON(http.StatusOK, gin.H{"user": serializer.Response()})
}

func UserRegister(c *gin.Context) {
	registerValidator := RegisterValidator{}
	if err := c.ShouldBindJSON(&registerValidator); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userModel, err := Register(registerValidator)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	UpdateContextUserModel(c, userModel.ID)
	serializer := UserSerializer{userModel}
	c.JSON(http.StatusOK, gin.H{"user": serializer.Response()})
}
