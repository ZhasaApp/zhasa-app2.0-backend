package base

type Pagination struct {
	PageSize int32
	Page     int32
}

func (p Pagination) GetOffset() int32 {
	if p.Page < 0 {
		return 0
	}
	return p.PageSize * (p.Page)
}
