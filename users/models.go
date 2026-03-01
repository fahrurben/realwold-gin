package users

import (
	"github.com/fahrurben/realworld-gin/common"

	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"column:username"`
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
