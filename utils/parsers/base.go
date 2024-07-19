package parsers

import (
	"github.com/shopspring/decimal"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ParamInt(p string) int {
	val, err := strconv.Atoi(p)
	if err != nil {
		return 0
	}
	return val
}

func ParamDecimal(p string) decimal.Decimal {
	val, err := strconv.ParseFloat(p, 64)
	if err != nil {
		return decimal.NewFromInt(0)
	}
	return decimal.NewFromFloat(val)
}

func ParamFloat32(p string) float32 {
	val, _ := strconv.ParseFloat(p, 32)
	return float32(val)
}

func ParamUint(p string) uint {
	return uint(ParamInt(p))
}

func ParamBool(p string) bool {
	val, _ := strconv.ParseBool(p)
	return val
}

func ParseFieldName(input string) string {
	i := strings.Index(input, "(")
	if i >= 0 {
		j := strings.Index(input, ")")
		if j >= 0 {
			return toCamelCase(input[i+1 : j])
		}
	}
	return ""
}

func toCamelCase(input string) string {
	words := strings.Split(input, "_")
	key := strings.Title(words[0])
	for _, word := range words[1:] {
		key += strings.Title(word)
	}
	return key
}

func CheckDate(format, date string) bool {
	t, err := time.Parse(format, date)
	if err != nil {
		return false
	}
	return t.Format(format) == date
}

var matchFirstCap = regexp.MustCompile("([A-Z])([A-Z][a-z])")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
