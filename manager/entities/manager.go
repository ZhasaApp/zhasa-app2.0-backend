package entities

type SalesManagerId int

type SalesManager struct {
	Id        SalesManagerId
	FirstName string
	LastName  string
	AvatarUrl string
}

type SalesManagers []SalesManager

type SalesManagerResponse struct {
	Id SalesManagerId `json:"id"`
}
