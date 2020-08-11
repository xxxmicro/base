package model


type TimeType string // 数据库的时间类型
const (
        DATETIME      TimeType = "datetime" // 时间类型 time.Time
        TIMESTAMP 		TimeType = "timestamp" // 时间戳 int64
)

type SortType string

const (
        SortType_DEFAULT SortType = "DEFAULT"
        SortType_ASC     SortType = "ASC" // 升序
        SortType_DSC     SortType = "DSC" // 降序
)

type SortSpec struct {
	Property   string   `json:"property"`   // 属性名
	Type       SortType `json:"type"`       // 排序类型
	IgnoreCase bool     `json:"ignoreCase"` // 忽略大小写
}

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