package entities

import (
	"time"
	. "zhasa2.0/branch/entities"
	entities2 "zhasa2.0/manager/entities"
	entities3 "zhasa2.0/sale/entities"
	"zhasa2.0/user/entities"
)

type BranchDirector struct {
	entities.User
	BranchDirectorId BranchDirectorId
	Branch           Branch
}

type BranchDirectorId int32

type SalesManagerGoal struct {
	SalesManagerId entities2.SalesManagerId
	FromDate       time.Time
	ToDate         time.Time
	Amount         entities3.SaleAmount
}
