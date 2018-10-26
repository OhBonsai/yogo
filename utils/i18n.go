package utils



var TDefault TranslateFunc


type TranslateFunc func(translationID string, params ...interface{}) string



func T(inp string, params ...interface{})  string {
	return inp
}