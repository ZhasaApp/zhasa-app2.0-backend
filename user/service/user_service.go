package service

import (
	. "zhasa2.0/user/entities"
	"zhasa2.0/user/repository"
)

type UserService interface {
	GetUserByPhone(phone Phone) (*User, error)
	CreateUser(request CreateUserRequest) error
	UploadAvatar(userId UserId, avatarUrl string) error
}

type DBUserService struct {
	repo repository.UserRepository
}

func (dus DBUserService) UploadAvatar(userId UserId, avatarUrl string) error {
	return dus.repo.UploadAvatar(userId, avatarUrl)
}

func NewUserService(repo repository.UserRepository) UserService {
	return DBUserService{
		repo: repo,
	}
}

func (dus DBUserService) CreateUser(request CreateUserRequest) error {
	return dus.repo.CreateUser(request)
}

func (dus DBUserService) GetUserByPhone(phone Phone) (*User, error) {
	return dus.repo.GetUserByPhone(phone)
}
