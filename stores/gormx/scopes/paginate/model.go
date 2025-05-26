package paginate

// Pagination 分页
type Pagination struct {
	Page          int   `json:"page"`
	PageSize      int   `json:"page_size"`
	Total         int64 `json:"total"`
	Rows          any   `json:"rows"`
	Extend        any   `json:"extend,omitempty"`
	TotalPage     int64 `json:"total_page"`
	NoQueryTotal  bool  `json:"no_query_total"`  // 是否不查询总数
	ForcePageSize bool  `json:"force_page_size"` // 是否强制获取多条数据
}

func (p *Pagination) Offset() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	return (p.Page - 1) * p.PageSize
}

func (p *Pagination) Limit() int {
	if p.PageSize == 0 {
		p.PageSize = 10
	}
	if !p.ForcePageSize && p.PageSize > 100 {
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
