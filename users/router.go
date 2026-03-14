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

func UserEndpoint(router *gin.RouterGroup) {
	router.GET("", GetCurrentUser)
	router.PUT("", UpdateUser)
}

func PublicProfilesRegister(router *gin.RouterGroup) {
	router.GET("/:username", GetProfile)
}

func PrivateProfilesRegister(router *gin.RouterGroup) {
	router.POST("/:username/follow", FollowUser)
	router.DELETE("/:username/follow", UnfollowUser)
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

	if err := userModel.checkPassword(loginValidator.User.Password); err != nil {
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

func GetCurrentUser(c *gin.Context) {
	userModel := c.MustGet("my_user_model").(UserModel)

	serializer := UserSerializer{model: &userModel}
	c.JSON(http.StatusOK, gin.H{"user": serializer.Response()})
}

func UpdateUser(c *gin.Context) {
	updateValidator := UpdateValidator{}
	userModel := c.MustGet("my_user_model").(UserModel)

	if err := c.ShouldBindJSON(&updateValidator); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedModel, err := Update(userModel, updateValidator)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	serializer := UserSerializer{model: updatedModel}
	c.JSON(http.StatusOK, gin.H{"user": serializer.Response()})
}

func GetProfile(c *gin.Context) {
	username := c.Param("username")
	loggedInUser := c.MustGet("my_user_model").(UserModel)

	userModel, err := FindOneUser(&UserModel{Username: username})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	serializer := ProfileSerialier{model: &userModel, loggedInUser: &loggedInUser}
	c.JSON(http.StatusOK, gin.H{"user": serializer.ProfileResponse()})
}

func FollowUser(c *gin.Context) {
	username := c.Param("username")
	loggedInUser := c.MustGet("my_user_model").(UserModel)

	follow, err := FindOneUser(&UserModel{Username: username})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = loggedInUser.followUser(follow)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	serializer := ProfileSerialier{model: &loggedInUser, loggedInUser: &loggedInUser}
	c.JSON(http.StatusOK, gin.H{"user": serializer.ProfileResponse()})
}

func UnfollowUser(c *gin.Context) {
	username := c.Param("username")
	loggedInUser := c.MustGet("my_user_model").(UserModel)

	follow, err := FindOneUser(&UserModel{Username: username})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = loggedInUser.unfollowUser(follow)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	serializer := ProfileSerialier{model: &loggedInUser, loggedInUser: &loggedInUser}
	c.JSON(http.StatusOK, gin.H{"user": serializer.ProfileResponse()})
}
