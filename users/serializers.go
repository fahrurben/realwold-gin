package users

import (
	"github.com/fahrurben/realworld-gin/common"
)

type UserSerializer struct {
	model *UserModel
}

type UserResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Bio      string `json:"bio"`
	Image    string `json:"image"`
	Token    string `json:"token"`
}

func (self *UserSerializer) Response() UserResponse {
	myUserModel := self.model
	image := ""
	if myUserModel.Image != "" {
		image = myUserModel.Image
	}
	user := UserResponse{
		Username: myUserModel.Username,
		Email:    myUserModel.Email,
		Bio:      myUserModel.Bio,
		Image:    image,
		Token:    common.GenToken(myUserModel.ID),
	}
	return user
}

type ProfileSerialier struct {
	model        *UserModel
	loggedInUser *UserModel
}

type ProfileResponse struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	Following bool   `json:"following"`
}

func (self *ProfileSerialier) ProfileResponse() ProfileResponse {
	db := common.GetDB()

	myUserModel := self.model
	image := ""
	if myUserModel.Image != "" {
		image = myUserModel.Image
	}

	following := false
	if self.loggedInUser != self.model {
		count := db.Model(&self.loggedInUser).Where(&UserModel{Username: myUserModel.Username}).Association("Follows").Count()
		if count > 0 {
			following = true
		}
	}
	profile := ProfileResponse{
		Username:  myUserModel.Username,
		Email:     myUserModel.Email,
		Bio:       myUserModel.Bio,
		Image:     image,
		Following: following,
	}
	return profile
}
