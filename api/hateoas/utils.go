package hateoas

import (
	"encoding/json"
	"strings"

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
