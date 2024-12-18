package paginate

// Pagination 分页
type Pagination struct {
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
	Rows     any   `json:"rows"`
	Extend   any   `json:"extend,omitempty"`
}

func (p *Pagination) Offset() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	return (p.Page - 1) * p.PageSize
}

func (p *Pagination) Limit() int {
	if p.PageSize < 10 {
		p.PageSize = 10
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
	return p.PageSize
}

func (p *Pagination) GetPage() int {
	return p.Page
}

func (p *Pagination) GetPageSize() int {
	return p.PageSize
}
