package base_controller

import (
	"cookie_supply_management/pkg/base/base_model"
	"cookie_supply_management/pkg/base/base_repository"
	"cookie_supply_management/utils/parsers"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"reflect"

	"strconv"
	"strings"
)

func GetQuery(c *gin.Context, a interface{}) base_repository.Searcher {
	search := c.Query("search")
	if search == "" {
		return nil
	}

	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(search), &raw); err != nil {
		return nil
	}

	data := reflect.ValueOf(a)
	myModel := fmt.Sprintf("%v", data.Type())[10:] + "s"

	var (
		query      string
		queryJoin  string
		isJoin     bool
		JoinModels string
		i          int
	)

	lenRaw := len(raw)
	for key, v := range raw {
		value := fmt.Sprintf("%v", v)
		typeValue := detectType(value, key)
		if typeValue == "map[string]interface {}" {
			JoinModels = key
			subQueryJoin, subIsJoin := processMap(v, lenRaw)
			queryJoin += subQueryJoin
			isJoin = subIsJoin
		} else {
			queryJoin += buildQueryPart(myModel, key, value, typeValue, lenRaw, i)
		}
		i++
	}

	return base_repository.NewSearcher(query, isJoin, JoinModels, queryJoin)
}

func detectType(value, key string) string {
	if strings.Contains(key, "json") {
		return "json"
	}
	if parsers.CheckDate("2006-01-02", value) {
		return "date"
	}
	if strings.Contains(value, "and") {
		return "date"
	}
	return fmt.Sprintf("%T", value)
}

func processMap(v interface{}, lenRaw int) (string, bool) {
	var (
		subQueryJoin string
		subIsJoin    bool
		j            int
	)

	byteV, _ := json.Marshal(v)
	var rawV map[string]interface{}
	if err := json.Unmarshal(byteV, &rawV); err != nil {
		return "", false
	}

	for k, val := range rawV {
		subQueryJoin += buildQueryPart("", k, fmt.Sprintf("%v", val), fmt.Sprintf("%T", val), lenRaw, j)
		j++
	}
	subIsJoin = true
	return subQueryJoin, subIsJoin
}

func buildQueryPart(model, key, value, typeValue string, lenRaw, i int) string {
	var (
		queryPart string
		from, to  string
	)

	if strings.Contains(value, "and") {
		parts := strings.Split(value, " and ")
		from, to = parts[0], parts[1]
	}

	switch typeValue {
	case "bool":
		queryPart = fmt.Sprintf("%s.%s = %s", model, key, value)
	case "string":
		queryPart = fmt.Sprintf("%s.%s ILIKE '%%%s%%'", model, key, value)
	case "int", "float64":
		queryPart = fmt.Sprintf("%s.%s = %s", model, key, value)
	case "date":
		if strings.Contains(value, "and") {
			queryPart = fmt.Sprintf("%s BETWEEN '%s' AND '%s'", key, from, to)
		} else {
			queryPart = fmt.Sprintf("%s.%s::text LIKE '%%%s%%'", model, key, value)
		}
	case "json":
		lang, word := splitJson(value)
		queryPart = fmt.Sprintf("%s.%s ->> '%s' ILIKE '%%%s%%'", model, key, lang, word)
	}

	if value == "not_null" {
		queryPart = fmt.Sprintf("%s.%s IS NOT NULL", model, key)
	}

	if i < lenRaw-1 {
		queryPart += " and "
	}

	return queryPart
}

func splitJson(value string) (string, string) {
	parts := strings.Split(value, " = ")
	return parts[0], parts[1]
}

func GetOrder(ctx *gin.Context, model base_model.HasId) base_repository.OrderFilter {
	var values []string = ctx.Request.URL.Query()["order_by"]
	return base_repository.NewOrder(values, model)
}

func GetPager(ctx *gin.Context) base_repository.Pager {
	page, _ := strconv.Atoi(ctx.Query("page"))
	if page == 0 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(ctx.Query("page_size"))
	switch {
	case pageSize > 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	return base_repository.NewPager(page, pageSize, offset)
}

func SelectedFields(ctx *gin.Context) string {
	fields := ctx.Request.URL.Query().Get("fields")
	if fields == "" {
		return "*"
	}

	fieldsList := strings.Split(fields, ",")

	var selectedFields string
	for _, field := range fieldsList {
		selectedFields += field + ","
	}
	selectedFields = selectedFields[:len(selectedFields)-1]

	return selectedFields
}
