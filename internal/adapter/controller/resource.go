package controller

type Link struct {
	Href   string `json:"href"`
	Method string `json:"method,omitempty"`
}

type Resource[T any] struct {
	Data  T               `json:"data"`
	Links map[string]Link `json:"links,omitempty"`
}

type Pagination struct {
	CurrentPage int  `json:"currentPage"`
	PerPage     int  `json:"perPage"`
	Total       int  `json:"total"`
	TotalPages  int  `json:"totalPages"`
	HasPrevious bool `json:"hasPrevious"`
	HasNext     bool `json:"hasNext"`
}

type ResourceCollection[T any] struct {
	Data       []T         `json:"data"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Links      []Link      `json:"links,omitempty"`
}

func NewResource[T any](data T, links map[string]Link) *Resource[T] {
	return &Resource[T]{
		Data:  data,
		Links: links,
	}
}

func NewResourceCollection[T any](data []T, pagination *Pagination, links ...Link) *ResourceCollection[T] {
	return &ResourceCollection[T]{
		Data:       data,
		Pagination: pagination,
		Links:      links,
	}
}
