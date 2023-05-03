package repository

import (
	"context"
	db_generated "zhasa2.0/db/sqlc"
	"zhasa2.0/user/entities"
)

type UserRepository interface {
	CreateUser(request entities.CreateUserRequest) error
	GetUserByPhone(phone entities.Phone) (*entities.User, error)
	GetUserById(id entities.UserId) (*entities.User, error)
	AddUserCode(userId entities.UserId, code entities.OtpCode) (entities.OtpId, error)
	GetAuthCodeById(id entities.OtpId) (*entities.UserAuth, error)
}

type PostgresUserRepository struct {
	ctx     context.Context
	querier db_generated.Querier
}

func (pur PostgresUserRepository) AddUserCode(userId entities.UserId, code entities.OtpCode) (entities.OtpId, error) {
	params := db_generated.CreateUserCodeParams{
		UserID: int32(userId),
		Code:   int32(code),
	}
	res, err := pur.querier.CreateUserCode(pur.ctx, params)
	if err != nil {
		return 0, err
	}
	return entities.OtpId(res), err
}

func (pur PostgresUserRepository) GetAuthCodeById(otpId entities.OtpId) (*entities.UserAuth, error) {
	data, err := pur.querier.GetAuthCodeById(pur.ctx, int32(otpId))
	if err != nil {
		return nil, err
	}
	return &entities.UserAuth{
		Code:      entities.OtpCode(data.Code),
		UserId:    entities.UserId(data.UserID),
		CreatedAt: data.CreatedAt,
	}, err
}

func (pur PostgresUserRepository) GetUserById(userId entities.UserId) (*entities.User, error) {
	res, err := pur.querier.GetUserById(pur.ctx, int32(userId))
	if err != nil {
		return nil, err
	}
	user := entities.User{
		Id:        res.ID,
		Phone:     entities.Phone(res.Phone),
		Avatar:    entities.Avatar{},
		FirstName: entities.Name(res.FirstName),
		LastName:  entities.Name(res.LastName),
	}
	return &user, err
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
