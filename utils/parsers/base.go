package parsers

import (
	"github.com/shopspring/decimal"
	"strconv"
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
