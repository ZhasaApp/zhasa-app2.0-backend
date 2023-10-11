package base

type Pagination struct {
	PageSize int32
	Page     int32
}

func (p Pagination) GetOffset() int32 {
	return p.PageSize * (p.Page - 1)
}
