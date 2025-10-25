package controller

type Link struct {
	Href   string `json:"href"`
	Method string `json:"method,omitempty"`
}

type Resource struct {
	Data  any             `json:"data"`
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

type ResourceCollection struct {
	Data       []any       `json:"data"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Links      []Link      `json:"links,omitempty"`
}

func NewResource(data any, links map[string]Link) *Resource {
	return &Resource{
		Data:  data,
		Links: links,
	}
}

func NewResourceCollection(data []any, pagination *Pagination, links ...Link) *ResourceCollection {
	return &ResourceCollection{
		Data:       data,
		Pagination: pagination,
		Links:      links,
	}
}
