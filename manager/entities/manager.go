package entities

import "zhasa2.0/branch/entities"

type SalesManagerId int

type SalesManager struct {
	Id        SalesManagerId
	FirstName string
	LastName  string
	AvatarUrl string
	Branch    entities.Branch
}

type SalesManagers []SalesManager

type SalesManagerResponse struct {
	Id SalesManagerId `json:"id"`
}
