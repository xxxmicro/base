package model

type FilterType string

const (
	FilterType_EQ       FilterType = "EQ"       //相等
	FilterType_NE       FilterType = "NE"       //不相等
	FilterType_GT       FilterType = "GT"       //大于
	FilterType_GTE      FilterType = "GTE"      //大于等于
	FilterType_LT       FilterType = "LT"       //小于
	FilterType_LTE      FilterType = "LTE"      //小于等于
	FilterType_IN       FilterType = "IN"       //在什么范围内
	FilterType_NOT_IN   FilterType = "NOT_IN"   //不在什么范围内
	FilterType_LIKE     FilterType = "LIKE"     //like
	FilterType_NOT_LIKE FilterType = "NOT_LIKE" //not like
	FilterType_MATCH    FilterType = "MATCH"    //匹配
	FilterType_BETWEEN  FilterType = "BETWEEN"  //匹配
	FilterType_IS_NULL  FilterType = "IS_NULL"  //为空
	FilterType_NOT_NULL FilterType = "NOT_NULL" //不为空
	FilterType_AND      FilterType = "AND"      //AND
	FilterType_OR       FilterType = "OR"       //OR
	FilterType_NOR      FilterType = "NOR"      //NOR

	FilterType_ES_EQ            FilterType = "EQ"            // 等于
	FilterType_ES_OR            FilterType = "OR"            //
	FilterType_ES_AND           FilterType = "AND"           //
	FilterType_ES_NESTED        FilterType = "NESTED"        // 嵌套查询
	FilterType_ES_TERMS_SCORE   FilterType = "TERMS_SCORE"   //
	FilterType_ES_EQ_SCORE      FilterType = "EQ_SCORE"      //
	FilterType_ES_TERM_FILTER   FilterType = "TERM_FILTER"   // 查找单个value
	FilterType_ES_TERMS_FILTER  FilterType = "TERMS_FILTER"  // 查找多个value
	FilterType_ES_RANGER_FILTER FilterType = "RANGER_FILTER" // 左开右闭
	FilterType_ES_RANGEL_FILTER FilterType = "RANGEL_FILTER" // 左闭右开
	FilterType_ES_RANGE_FILTER  FilterType = "RANGE_FILTER"  // 左闭右闭
	FilterType_ES_LTE_FILTER    FilterType = "LTE_FILTER"    // 小于等于
	FilterType_ES_LT_FILTER     FilterType = "LT_FILTER"     // 小于
	FilterType_ES_GTE_FILTER    FilterType = "GTE_FILTER"    // 大于等于
	FilterType_ES_GT_FILTER     FilterType = "GT_FILTER"     // 大于
)

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
