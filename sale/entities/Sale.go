package entities

import (
	"time"
)

type Sale struct {
	Id              int32
	SaleType        SaleType
	SalesAmount     int64
	SaleDate        time.Time
	SaleDescription string
}

type SaleId int32

type SaleDescription string

type SalesBySaleType map[SaleType]Sale
