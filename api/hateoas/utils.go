package hateoas

import (
	"bytes"
	"encoding/json"
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
func JSONBFilter(db *gorm.DB, criterias map[string]interface{}) *gorm.DB {
	for k, v := range criterias {
		values := strings.Split(v.(string), ",")
		switch len(values) {
		case 1:
			db = db.Where(k+" (?)", values[0])
			break
		case 2:
			db = db.Where(k+" (?,?)", values[0], values[1])
			break
		}
	}
	return db
}

// InlineFilter filter inline filters
func InlineFilter(db *gorm.DB, criterias map[string]interface{}) *gorm.DB {
	for k, v := range criterias {
		values := strings.Split(v.(string), ",")
		switch len(values) {
		case 1:
			db = db.Where(k, values[0])
			break
		case 2:
			db = db.Where(k, values[0], values[1])
			break
		}
	}
	return db
}

// CheckFilter analyse filters
func CheckFilter(criterias map[string]interface{}) (map[string]interface{}, map[string]interface{}, map[string]interface{}) {
	// Analyse critarias for extract inline, standard and JSONB ones
	standardCriterias := make(map[string]interface{})
	inlineCriterias := make(map[string]interface{})
	jsonbCriterias := make(map[string]interface{})
	for k, v := range criterias {
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

// unCamelCase discover physical column name in database
func unCamelCase(value string) string {
	if strings.Contains(value, "->>") {
		// on jsonb column don't un camel case field
		return value
	}
	accu := ""
	for i := 0; i < len(value); i = i + 1 {
		switch value[i] {
		case '_':
			break
		default:
			// is already lower
			if bytes.ToLower([]byte{value[i]})[0] == value[i] {
				accu = accu + string(value[i])
			} else {
				accu = accu + "_" + string(bytes.ToLower([]byte{value[i]})[0])
			}
		}
	}
	return accu
}

// BaseURL returns the base path that has been used to access current resource
func BaseURL(c *gin.Context) string {
	basePath, ok := c.Get(hateoasBasePathKey)
	if ok {
		return basePath.(string)
	}
	return c.Request.URL.EscapedPath()
}
