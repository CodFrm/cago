package httputils

type PageRequest[T any] struct {
	Page  int `form:"page" binding:"default=1"`
	Limit int `form:"limit" binding:"default=20"`
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

type PageResponse[T any] struct {
	List  []T   `json:"list"`
	Total int64 `json:"total"`
}
