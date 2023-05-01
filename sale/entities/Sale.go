package entities

import (
	"time"
	"zhasa2.0/manager/entities"
)

type Sale struct {
	SaleManagerId   entities.SalesManagerId
	SalesTypeId     SaleTypeId
	SalesAmount     SaleAmount
	SaleDate        time.Time
	SaleDescription SaleDescription
}

type SaleDescription string

type SalesBySaleType map[SaleType]Sale
