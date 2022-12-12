package httputils

type PageRequest struct {
	Page  int `form:"page" binding:"default=1"`
	Limit int `form:"limit" binding:"default=20"`
}

func (p *PageRequest) GetPage() int {
	if p.Page == 0 {
		return 1
	}
	return p.Page
}

func (p *PageRequest) GetOffset() int {
	return (p.GetPage() - 1) * p.Limit
}

func (p *PageRequest) GetLimit() int {
	if p.Limit == 0 {
		return 20
	}
	return p.Limit
}

type PageResponse[T any] struct {
	List  []T   `json:"list"`
	Total int64 `json:"total"`
}
