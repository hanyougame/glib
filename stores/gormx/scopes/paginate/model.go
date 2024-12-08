package paginate

// Pagination 分页
type Pagination[T any] struct {
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
	Rows     T     `json:"rows"`
}

func (p *Pagination[T]) Offset() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	return (p.Page - 1) * p.PageSize
}

func (p *Pagination[T]) Limit() int {
	if p.PageSize < 10 {
		p.PageSize = 10
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
	return p.PageSize
}

func (p *Pagination[T]) GetPage() int {
	return p.Page
}

func (p *Pagination[T]) GetPageSize() int {
	return p.PageSize
}
