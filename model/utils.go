package model

import (
	"bytes"
	"crypto/rand"
	"encoding/base32"
	"github.com/pborman/uuid"
	goi18n "github.com/nicksnyder/go-i18n/i18n"
	"time"
	"encoding/json"
)

var encoding = base32.NewEncoding("ybndrfg8ejkmcpqxot1uwisza345h769")


type AppError struct {
	Id            string `json:"id"`
	Message       string `json:"message"`               // Message to be display to the end user without debugging information
	DetailedError string `json:"detailed_error"`        // Internal error string to help the developer
	RequestId     string `json:"request_id,omitempty"`  // The RequestId that's also set in the header
	StatusCode    int    `json:"status_code,omitempty"` // The http status code
	Where         string `json:"-"`                     // The function where it happened in the form of Struct.Func
	IsOAuth       bool   `json:"is_oauth,omitempty"`    // Whether the error is OAuth specific
	params        map[string]interface{}
}


type StringInterface map[string]interface{}
type StringMap map[string]string
type StringArray []string

func (er *AppError) Error() string {
	return er.Where + ": " + er.Message + ", " + er.DetailedError
}

// change string interface to string json
func StringInterfaceToJson(objmap map[string]interface{}) string {
	b, _ := json.Marshal(objmap)
	return string(b)
}

func MapToJson(objmap map[string]string) string {
	b, _ := json.Marshal(objmap)
	return string(b)
}

func ArrayToJson(objmap []string) string {
	b, _ := json.Marshal(objmap)
	return string(b)
}




func NewAppError(where string, id string, params map[string]interface{}, details string, status int) *AppError {
	ap := &AppError{}
	ap.Id = id
	ap.params = params
	ap.Message = id
	ap.Where = where
	ap.DetailedError = details
	ap.StatusCode = status
	ap.IsOAuth = false
	return ap
}

func (er *AppError) SystemMessage(T goi18n.TranslateFunc) string {
	if er.params == nil {
		return T(er.Id)
	} else {
		return T(er.Id, er.params)
	}
}


func NewId() string {
	var b bytes.Buffer
	encoder := base32.NewEncoder(encoding, &b)
	encoder.Write(uuid.NewRandom())
	encoder.Close()
	b.Truncate(26) // removes the '==' padding
	return b.String()
}

func NewRandomString(length int) string {
	var b bytes.Buffer
	str := make([]byte, length+8)
	rand.Read(str)
	encoder := base32.NewEncoder(encoding, &b)
	encoder.Write(str)
	encoder.Close()
	b.Truncate(length) // removes the '==' padding
	return b.String()
}

func GetMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func (er *AppError) ToJson() string {
	b, _ := json.Marshal(er)
	return string(b)
}

func CopyStringMap(originalMap map[string]string) map[string]string {
	copyMap := make(map[string]string)
	for k, v := range originalMap {
		copyMap[k] = v
	}
	return copyMap
}


func (er *AppError) Translate(T goi18n.TranslateFunc) {
	if T == nil {
		er.Message = er.Id
		return
	}

	if er.params == nil {
		er.Message = T(er.Id)
	} else {
		er.Message = T(er.Id, er.params)
	}
}