package hateoas

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// ToJSON simple indented json
func ToJSON(data interface{}) string {
	byt, _ := json.MarshalIndent(data, "", "\t")
	return string(byt)
}

// JSONBFilter filter jsonb
func JSONBFilter(db *gorm.DB, criteria map[string]interface{}) *gorm.DB {
	for k, v := range criteria {
		values := strings.Split(v.(string), ",")
		switch len(values) {
		case 1:
			db = db.Where(k+" (?)", values[0])
		case 2:
			db = db.Where(k+" (?,?)", values[0], values[1])
		}
	}
	return db
}

// InlineFilter filter inline filters
func InlineFilter(db *gorm.DB, criteria map[string]interface{}) *gorm.DB {
	for k, v := range criteria {
		values := strings.Split(v.(string), ",")
		switch len(values) {
		case 1:
			db = db.Where(k, values[0])
		case 2:
			db = db.Where(k, values[0], values[1])
		}
	}
	return db
}

// CheckFilter analyse filters
func CheckFilter(criteria map[string]interface{}) (map[string]interface{}, map[string]interface{}, map[string]interface{}) {
	// Analyse critarias for extract inline, standard and JSONB ones
	standardCriterias := make(map[string]interface{})
	inlineCriterias := make(map[string]interface{})
	jsonbCriterias := make(map[string]interface{})
	for k, v := range criteria {
		if strings.Contains(k, ".") {
			var values = strings.Split(k, ".")
			key := values[0] + "->>'" + values[1] + "' in "
			jsonbCriterias[key] = v
		} else if strings.Contains(k, "@>") {
			var values = strings.Split(k, "@")
			key := values[0] + " @>"
			jsonbCriterias[key] = v
		} else {
			if strings.Contains(k, "?") {
				inlineCriterias[k] = v
			} else {
				standardCriterias[k] = v
			}
		}
	}
	return standardCriterias, inlineCriterias, jsonbCriterias
}

// BaseURL returns the base path that has been used to access current resource
func BaseURL(c *gin.Context) string {
	basePath, ok := c.Get(hateoasBasePathKey)
	if ok {
		return basePath.(string)
	}
	return c.Request.URL.EscapedPath()
}

// GetGormSortClause returns the SQL-escaped sort clause
func (p Pageable) GetGormSortClause() interface{} {
	if sortClause := p.GetSortClause(); sortClause != "1" {
		return sortClause
	}
	// wrap column pointer in a gorm expression to avoid quote-surrounding
	return gorm.Expr("1")
}

// GetSortClause returns the SQL-escaped sort clause
func (p Pageable) GetSortClause() string {
	// if not sort column was specified, optimistically use first column to preserve page consistency
	if p.Sort == "" {
		return "1"
	}
	fields := strings.Split(p.Sort, ",")
	for i, field := range fields {
		// for each field to sort, read sort direction after a semicolon, asc will be used as default
		direction := directionAsc
		if fieldClause := strings.Split(field, ";"); len(fieldClause) == 2 {
			field = fieldClause[0]
			if strings.ToLower(fieldClause[1]) == directionDesc {
				direction = directionDesc
			}
		}
		// %q sanitizes the field double-quotes to prevent sql injections
		fields[i] = fmt.Sprintf("%q %s", field, direction)
	}
	return strings.Join(fields, ", ")
}

// GetOffset returns the page offset
func (p Pageable) GetOffset() int {
	return p.Page * p.Size
}

// NewPage initialize an empty resource page
func NewPage(pageable Pageable, defaultPageSize int, basePath string) Page {
	if pageable.Size == 0 {
		pageable.Size = defaultPageSize
	}
	return Page{Pageable: pageable, BasePath: basePath}
}
