package utils

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func CapitalizeFirstOnly(s string) string {
	if len(s) == 0 {
		return s
	}

	// Извлечение первого символа и его длины в байтах
	firstRune, size := utf8.DecodeRuneInString(s)

	// Преобразуем первый символ к верхнему регистру
	first := string(unicode.ToUpper(firstRune))

	lower := strings.ToLower(s)
	// Создаем строку с первым символом верхнего регистра и остальными символами без изменений
	return first + lower[size:]
}
