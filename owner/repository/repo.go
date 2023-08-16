package repository

import (
	"context"
	generated "zhasa2.0/db/sqlc"
	. "zhasa2.0/user/entities"
)

type OwnerRepository interface {
	GetOwnerByUserId(userId int32) (*User, error)
}

type DbOwnerRepositroy struct {
	ctx     context.Context
	querier generated.Querier
}

func (d DbOwnerRepositroy) GetOwnerByUserId(userId int32) (*User, error) {
	data, err := d.querier.GetOwnerByUserId(d.ctx, userId)
	if err != nil {
		return nil, err
	}

	user := User{
		Id:        data.UserID,
		Phone:     Phone(data.Phone),
		Avatar:    data.AvatarUrl,
		FirstName: Name(data.FirstName),
		LastName:  Name(data.LastName),
	}

	return &user, err
}

func NewOwnerRepo(ctx context.Context, querier generated.Querier) OwnerRepository {
	return DbOwnerRepositroy{
		ctx:     ctx,
		querier: querier,
	}
}
