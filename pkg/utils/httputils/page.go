package httputils

type PageRequest[T any] struct {
	Page  int `form:"page,default=1"`
	Limit int `form:"limit,default=20"`
}

func (p *PageRequest[T]) GetPage() int {
	if p.Page == 0 {
		return 1
	}
	return p.Page
}

func (p *PageRequest[T]) GetOffset() int {
	return (p.GetPage() - 1) * p.Limit
}

func (p *PageRequest[T]) GetLimit() int {
	if p.Limit == 0 {
		return 20
	}
	return p.Limit
}

func (p *PageRequest[T]) Response(list []T, total int64) (*PageResponse[T], error) {
	return &PageResponse[T]{
		List:  list,
		Total: total,
	}, nil
}

type IPageResponse interface {
	GetList() interface{}
	GetTotal() int64
}

type IPageDataResponse interface {
	IPageResponse
	GetData() interface{}
}

type PageResponse[T any] struct {
	List  []T   `json:"list"`
	Total int64 `json:"total"`
}

func (p *PageResponse[T]) GetList() interface{} {
	return p.List
}

func (p *PageResponse[T]) GetTotal() int64 {
	return p.Total
}