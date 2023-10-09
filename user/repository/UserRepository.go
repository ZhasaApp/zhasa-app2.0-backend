package repository

import (
	"context"
	"fmt"
	db_generated "zhasa2.0/db/sqlc"
	. "zhasa2.0/user/entities"
)

type UserRepository interface {
	CreateUser(request CreateUserRequest) (int32, error)
	GetUserByPhone(phone Phone) (*User, error)
	GetUserById(id int32) (*User, error)
	AddUserCode(userId UserId, code OtpCode) (OtpId, error)
	GetAuthCodeById(id OtpId) (*UserAuth, error)
	UploadAvatar(userId UserId, avatarUrl string) error
	DeleteAvatar(userId UserId) error
}

type PostgresUserRepository struct {
	ctx     context.Context
	querier db_generated.Querier
}

func (pur PostgresUserRepository) DeleteAvatar(userId UserId) error {
	err := pur.querier.DeleteUserAvatar(pur.ctx, int32(userId))
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}

func (pur PostgresUserRepository) UploadAvatar(userId UserId, avatarUrl string) error {
	err := pur.querier.UploadUserAvatar(pur.ctx, db_generated.UploadUserAvatarParams{
		UserID:    int32(userId),
		AvatarUrl: avatarUrl,
	})
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (pur PostgresUserRepository) AddUserCode(userId UserId, code OtpCode) (OtpId, error) {
	params := db_generated.CreateUserCodeParams{
		UserID: int32(userId),
		Code:   int32(code),
	}
	res, err := pur.querier.CreateUserCode(pur.ctx, params)
	if err != nil {
		return 0, err
	}
	return OtpId(res), err
}

func (pur PostgresUserRepository) GetAuthCodeById(otpId OtpId) (*UserAuth, error) {
	data, err := pur.querier.GetAuthCodeById(pur.ctx, int32(otpId))
	if err != nil {
		return nil, err
	}
	return &UserAuth{
		Code:      OtpCode(data.Code),
		UserId:    UserId(data.UserID),
		CreatedAt: data.CreatedAt,
	}, err
}

func (pur PostgresUserRepository) GetUserById(userId int32) (*User, error) {
	res, err := pur.querier.GetUserById(pur.ctx, userId)
	if err != nil {
		return nil, err
	}

	user := User{
		Id:        res.ID,
		Phone:     Phone(res.Phone),
		Avatar:    res.AvatarUrl,
		FirstName: Name(res.FirstName),
		LastName:  Name(res.LastName),
		UserRole: UserRole{
			Id:  res.RoleID,
			Key: res.Key,
		},
	}
	return &user, err
}

func NewUserRepository(ctx context.Context, querier db_generated.Querier) UserRepository {
	return PostgresUserRepository{
		ctx:     ctx,
		querier: querier,
	}
}

func (pur PostgresUserRepository) GetUserByPhone(phone Phone) (*User, error) {
	res, err := pur.querier.GetUserByPhone(pur.ctx, string(phone))

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	user := User{
		Id:        res.ID,
		Phone:     Phone(res.Phone),
		Avatar:    res.AvatarUrl,
		FirstName: Name(res.FirstName),
		LastName:  Name(res.LastName),
	}
	return &user, err
}

func (pur PostgresUserRepository) CreateUser(request CreateUserRequest) (int32, error) {
	params := db_generated.CreateUserParams{
		Phone:     string(request.Phone),
		FirstName: string(request.FirstName),
		LastName:  string(request.LastName),
	}
	return pur.querier.CreateUser(pur.ctx, params)
}
