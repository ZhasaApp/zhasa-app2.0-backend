package entities

type SaleTypeId int32

type CreateSaleTypeBody struct {
	Title       string
	Description string
}

type SaleType struct {
	Id          SaleTypeId
	Title       string
	Description string
}
