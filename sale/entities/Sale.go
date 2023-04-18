package entities

import (
	"time"
)

type Sale struct {
	SalesType   SaleType
	SalesAmount SaleAmount
	SaleDate    time.Time
}

type SalesBySaleType map[SaleType]Sale
