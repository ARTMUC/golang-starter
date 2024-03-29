package crud

const (
	AND           = "$and"
	OR            = "$or"
	SEPARATOR     = "||"
	SortSeparator = ","
)

type PaginationResponse[T any] struct {
	Data       []*T  `json:"data"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"totalPages"`
}

type GetAllRequest struct {
	Joins    []string               `json:"-" form:"-"`
	Page     int                    `json:"page" form:"page"`
	Limit    int                    `json:"limit" form:"limit"`
	Preloads string                 `json:"join" form:"join"`
	S        string                 `json:"s" form:"s"`
	C        map[string]interface{} `json:"c" form:"c"`
	Fields   string                 `json:"fields" form:"fields"`
	Filter   []string               `json:"filter" form:"filter"`
	Sort     []string               `json:"sort" form:"sort"`
}

var filterConditions = map[string]string{
	"eq":      "=",
	"ne":      "!=",
	"gt":      ">",
	"lt":      "<",
	"gte":     ">=",
	"lte":     "<=",
	"$in":     "in",
	"cont":    "ILIKE",
	"isnull":  "IS NULL",
	"notnull": "IS NOT NULL",
}

type ById struct {
	ID string `uri:"id" binding:"required"`
}
