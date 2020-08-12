package model

type PageQuery struct {
	Filters map[string]interface{}		`json:"filters"`
	PageNo int											`json:"pageNo"`
	PageSize int										`json:"pageSize"`
	Sort 		[]*SortSpec								`json:"sort"`
}

type Page struct {
	Content interface{}	`json:"content"`
	Total int `json:"total"`
	PageNo int	`json:"pageNo"`
	PageSize int	`json:"pageSize"`
	PageCount int `json:"pageCount"`
}
