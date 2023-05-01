package service

import (
	"zhasa2.0/user/entities"
	"zhasa2.0/user/repository"
)

type UserService interface {
	GetUserByPhone(phone entities.Phone) (*entities.User, error)
	CreateUser(request entities.CreateUserRequest) error
}

type DBUserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return DBUserService{
		repo: repo,
	}
}

func (dus DBUserService) CreateUser(request entities.CreateUserRequest) error {
	return dus.repo.CreateUser(request)
}

func (dus DBUserService) GetUserByPhone(phone entities.Phone) (*entities.User, error) {
	return dus.repo.GetUserByPhone(phone)
}
