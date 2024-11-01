package pager

const (
	_defaultPage     = 1
	_defaultPageSize = 20
)

type Request struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

func New(page, pageSize int) Request {
	if page <= 0 {
		page = _defaultPage
	}
	if pageSize <= 0 {
		pageSize = _defaultPageSize
	}
	return Request{
		Page:     page,
		PageSize: pageSize,
	}
}

func (req Request) Offset() int {
	return (req.Page - 1) * req.PageSize
}

func (req Request) Limit() int {
	return req.PageSize
}

func (req Request) Response(totalCount int64) *Response {
	return &Response{
		Page:      req.Page,
		PageSize:  req.PageSize,
		TotalPage: (int(totalCount) + req.PageSize - 1) / req.PageSize,
	}
}

type Response struct {
	Page      int `json:"page"`
	PageSize  int `json:"page_size"`
	TotalPage int `json:"total_page"`
}
