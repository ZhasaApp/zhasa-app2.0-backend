package repository

import (
	"context"
	"errors"
	"fmt"
	generated "zhasa2.0/db/sqlc"
	. "zhasa2.0/user/entities"
)

type AddUserCodeFunc func(userId UserId, code OtpCode) (OtpId, error)

func NewAddUserCodeFunc(ctx context.Context, store generated.UserStore) AddUserCodeFunc {
	return func(userId UserId, code OtpCode) (OtpId, error) {
		params := generated.CreateUserCodeParams{
			UserID: int32(userId),
			Code:   int32(code),
		}
		id, err := store.CreateUserCode(ctx, params)

		if err != nil {
			fmt.Println(err)
			return 0, errors.New("error creating user code")
		}
		return OtpId(id), nil
	}
}
