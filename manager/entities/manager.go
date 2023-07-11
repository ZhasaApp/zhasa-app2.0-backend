package entities

import (
	"time"
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

type SalesManagers []SalesManager

type SalesManagerResponse struct {
	Id SalesManagerId `json:"id"`
}

type EditSaleBody struct {
	ID     int32     `json:"id"`
	Date   time.Time `json:"date"`
	TypeID int32     `json:"type_id"`
	Value  int64     `json:"value"`
	Title  string    `json:"title"`
}
