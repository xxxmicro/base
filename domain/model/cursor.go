package model

type CursorQuery struct {
	Filters    map[string]interface{} 	`json:"filters"`    // 筛选条件
	Cursor     interface{}            	`json:"cursor"`     // 游标值
	CursorSort *SortSpec              	`json:"cursorSort"` // 游标字段&排序
	Size       int                  	`json:"size"`       // 数据量
	Direction  byte                   	`json:"direction"`  // 查询方向 0：游标前；1：游标后
}

type CursorList struct {
	Extra CursorExtra			`json:"extra"`
	Items   interface{} 		`json:"items"`   // 数据列表指针
}

type CursorExtra struct {
	Direction byte        		`json:"direction"` // 查询方向 0：游标前；1：游标后
	Size      int       		`json:"size"`      // 数据量
	HasMore   bool        		`json:"hasMore"`   // 是否有更多数据
	MaxCursor interface{} 		`json:"maxCursor"` // 结果集中的起始游标值
	MinCursor interface{} 		`json:"minCursor"` // 结果集中的结束游标值
}
