package repository

import (
	"context"
	generated "zhasa2.0/db/sqlc"
	. "zhasa2.0/user/entities"
)

type GetAuthCodeByIdFunc func(otpId OtpId) (*UserAuth, error)

func NewGetAuthCodeByIdFunc(ctx context.Context, store generated.UserStore) GetAuthCodeByIdFunc {
	return func(otpId OtpId) (*UserAuth, error) {
		data, err := store.GetAuthCodeById(ctx, int32(otpId))
		if err != nil {
			return nil, err
		}
		return &UserAuth{
			Code:      OtpCode(data.Code),
			UserId:    UserId(data.UserID),
			CreatedAt: data.CreatedAt,
		}, nil
	}
}
