package entities

import (
	"time"
	entities2 "zhasa2.0/manager/entities"
	entities3 "zhasa2.0/sale/entities"
	"zhasa2.0/user/entities"
)

type BranchDirector struct {
	entities.User
	BranchDirectorId BranchDirectorId
}

type BranchDirectorId int32

type SalesManagerGoal struct {
	SalesManagerId entities2.SalesManagerId
	FromDate       time.Time
	ToDate         time.Time
	Amount         entities3.SaleAmount
}
