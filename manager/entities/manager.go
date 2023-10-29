package entities

import (
	"zhasa2.0/base"
	. "zhasa2.0/branch/entities"
	. "zhasa2.0/user/entities"
)

type SalesManagerId int

type RatingPlace int32

type SalesManager struct {
	Id          SalesManagerId
	FirstName   string
	LastName    string
	AvatarUrl   string
	Branch      Branch
	Ratio       base.Percent
	RatingPlace RatingPlace
	UserId      UserId
}

func (sm SalesManager) GetAvatarPointer() *string {
	if len(sm.AvatarUrl) == 0 {
		return nil
	}

	return &sm.AvatarUrl
}
