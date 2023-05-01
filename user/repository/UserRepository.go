package repository

import (
	"context"
	db_generated "zhasa2.0/db/sqlc"
	"zhasa2.0/user/entities"
)

type UserRepository interface {
	CreateUser(request entities.CreateUserRequest) error
	GetUserByPhone(phone entities.Phone) (*entities.User, error)
	AddUserCode(userId int32, code int32) (int32, error)
	GetActualUserCode(userId int32) (*entities.UserAuth, error)
}

type PostgresUserRepository struct {
	ctx     context.Context
	querier db_generated.Querier
}

func (pur PostgresUserRepository) AddUserCode(userId int32, code int32) (int32, error) {
	params := db_generated.CreateUserCodeParams{
		UserID: userId,
		Code:   code,
	}
	return pur.querier.CreateUserCode(pur.ctx, params)
}

func (pur PostgresUserRepository) GetActualUserCode(userId int32) (*entities.UserAuth, error) {
	data, err := pur.querier.GetUserCode(pur.ctx, userId)
	if err != nil {
		return nil, err
	}
	return &entities.UserAuth{
		Code:      entities.OtpCode(data.Code),
		UserId:    data.UserID,
		CreatedAt: data.CreatedAt,
	}, err
}

func NewUserRepository(ctx context.Context, querier db_generated.Querier) UserRepository {
	return PostgresUserRepository{
		ctx:     ctx,
		querier: querier,
	}
}

func (pur PostgresUserRepository) GetUserByPhone(phone entities.Phone) (*entities.User, error) {
	res, err := pur.querier.GetUserByPhone(pur.ctx, string(phone))
	user := entities.User{
		Id:        res.ID,
		Phone:     entities.Phone(res.Phone),
		Avatar:    entities.Avatar{},
		FirstName: entities.Name(res.FirstName),
		LastName:  entities.Name(res.LastName),
	}
	return &user, err
}

func (pur PostgresUserRepository) CreateUser(request entities.CreateUserRequest) error {
	params := db_generated.CreateUserParams{
		Phone:     string(request.Phone),
		FirstName: string(request.FirstName),
		LastName:  string(request.LastName),
	}
	return pur.querier.CreateUser(pur.ctx, params)
}
