package entities

type Pagination struct {
	PageNumber int64 `query:"page"`
	PageSize   int64 `query:"size"`
}

type PagedResponse struct {
	HasMore bool        `json:"has_more"`
	Total   int64       `json:"total"`
	Object  string      `json:"object"`
	Data    interface{} `json:"data"`
}

func NewDefaultPagination() Pagination {
	return Pagination{
		PageSize:   25,
		PageNumber: 1}
}

func (p *Pagination) GetPageStartIndex() int64 {
	p.preventZeroPageValue()
	return p.PageSize * (p.PageNumber - 1)
}

func (p *Pagination) HasMorePages(total int64) bool {
	return p.GetPageStartIndex()+p.PageSize < total
}

func (p *Pagination) preventZeroPageValue() {
	if p.PageNumber == 0 {
		p.PageNumber = 1
	}
}

func NewPagedResponse(result interface{}, hasMorePages bool, count int64) PagedResponse {
	return PagedResponse{
		HasMore: hasMorePages,
		Total:   count,
		Data:    result,
		Object:  "list",
	}
}

func (r *PagedResponse) IsEmpty() bool {
	return r.Total == 0
}
