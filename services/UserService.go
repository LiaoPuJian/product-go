package services

import (
	"errors"
	"product-go/models"
	"product-go/repositories"

	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	IsPwdSuccess(userName string, pwd string) (user *models.User, err error)
	AddUser(user *models.User) (userId int64, err error)
}

type UserService struct {
	UserRepository repositories.UserRepository
}

func NewUserService(repository repositories.UserRepository) IUserService {
	return &UserService{UserRepository: repository}
}

//判断当前密码是否有效
func (u *UserService) IsPwdSuccess(userName string, pwd string) (user *models.User, err error) {
	user, err = u.UserRepository.Select(userName)
	if err != nil {
		return
	}
	isOk, err := ValidatePassword(pwd, user.Password)
	if !isOk {
		return
	}
	return
}

func (u *UserService) AddUser(user *models.User) (userId int64, err error) {
	pwd, err := GeneratePassword(user.Password)
	if err != nil {
		return
	}
	user.Password = pwd
	return u.UserRepository.Insert(user)
}

//密码加密
func GeneratePassword(pwd string) (string, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(password), nil
}

//验证密码是否一致
func ValidatePassword(pwd, password string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(pwd)); err != nil {
		return false, errors.New("密码对比失败")
	}
	return true, nil
}
