package entities

type SalesManageId int

type SalesManager struct {
	Id        SalesManageId
	FirstName string
	LastName  string
	AvatarUrl string
}

type SalesManagers []SalesManager

type SalesManagerResponse struct {
	Id SalesManageId `json:"id"`
}
