package utils

import "github.com/nicksnyder/go-i18n/i18n"


var TDefault i18n.TranslateFunc


type TranslateFunc func(translationID string, params ...interface{}) string


func T(inp string, params ...interface{})  string {
	return inp
}