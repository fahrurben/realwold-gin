package users

import (
	"github.com/fahrurben/realworld-gin/common"

	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"column:username;uniqueIndex"`
	Email        string `gorm:"column:email;uniqueIndex"`
	Bio          string `gorm:"column:bio;size:1024"`
	Image        string `gorm:"column:image"`
	PasswordHash string `gorm:"column:password;not null"`
}

func FindOneUser(condition any) (UserModel, error) {
	db := common.GetDB()
	var model UserModel
	err := db.Where(condition).First(&model).Error
	return model, err
}

func (u *UserModel) checkPassword(password string) error {
	bytePassword := []byte(password)
	byteHashedPassword := []byte(u.PasswordHash)
	return bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)
}

func Register(data RegisterValidator) (*UserModel, error) {
	db := common.GetDB()

	var model UserModel = UserModel{}
	model.Username = data.User.Username
	model.Email = data.User.Email
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.User.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}
	model.PasswordHash = string(hashedPassword)

	error := db.Save(&model).Error
	return &model, error
}
