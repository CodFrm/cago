package httputils

// PageRequest 使用UnmarshalJSON的方法来反序列化,避免使用用户传递的page和limit
type PageRequest struct {
	// Deprecated 请使用方法GetPage
	Page int `form:"page" json:"page"`
	// Deprecated 请使用方法GetSize
	Size int `form:"size" json:"size"`
	// Deprecated 请使用方法GetSort
	Sort string `form:"sort" json:"sort"`
	// Deprecated 请使用方法GetOrder
	Order string `form:"order" json:"order"`
}

func (p *PageRequest) GetSort() string {
	return p.Sort
}

func (p *PageRequest) GetOrder() string {
	return p.Order
}

func (p *PageRequest) GetPage() int {
	if p.Page <= 0 {
		return 1
	}
	return p.Page
}

func (p *PageRequest) GetSize() int {
	if p.Size <= 0 || p.Size > 100 {
		return 20
	}
	return p.Size
}

func (p *PageRequest) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *PageRequest) GetLimit() int {
	return p.GetSize()
}

type PageResponse[T any] struct {
	List  []T   `json:"list"`
	Total int64 `json:"total"`
}
