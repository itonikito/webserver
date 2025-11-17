package utils

import "strings"

func IsValidParameter(param string) bool {
	if param == "" {
		return false
	}

	// Список игнорируемых путей
	ignoredPaths := []string{
		"favicon.ico",
		"robots.txt",
		"apple-touch-icon.png",
		"android-chrome",
	}

	for _, ignored := range ignoredPaths {
		if strings.HasPrefix(param, ignored) {
			return false
		}
	}

	// Если параметр содержит точку (скорее всего это файл), игнорируем
	if strings.Contains(param, ".") {
		return false
	}

	// Проверяем что параметр содержит только разрешенные символы
	for _, char := range param {
		if !isAllowedChar(char) {
			return false
		}
	}

	return true
}

func isAllowedChar(char rune) bool {
	// Разрешаем буквы, цифры, подчеркивания и дефисы
	return (char >= 'a' && char <= 'z') ||
		(char >= 'A' && char <= 'Z') ||
		(char >= '0' && char <= '9') ||
		char == '_' ||
		char == '-'
}
