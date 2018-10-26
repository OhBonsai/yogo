package model

import (
	"encoding/json"
	"io"
)

const (
	PLAN_MESSAGE_MAX_RUNES_V1 = 2000
	PLAN_NAME_MAX_BYTES_V1 = 50
	PLAN_PROPS_DELETE_BY = "deleteBy"
)


type Plan struct {
	Id         string `json:"id"`
	CreateAt   int64  `json:"create_at"`
	UpdateAt   int64  `json:"update_at"`
	DeleteAt   int64  `json:"delete_at"`

	Name       string          `json:"name"`
	Props      StringInterface `json:"props"`
}



func (o *Plan) IsValid() *AppError{
	return nil
}


func (o *Plan) PreSave() {
	if o.Id == "" {
		o.Id = NewId()
	}

	if o.CreateAt == 0 {
		o.CreateAt = GetMillis()
	}

	o.UpdateAt = o.CreateAt
	o.PreCommit()
}


func (o *Plan) PreCommit() {
	if o.Props == nil {
		o.Props = make(map[string]interface{})
	}
}

func (o *Plan) MakeNonNil() {
	if o.Props == nil {
		o.Props = make(map[string]interface{})
	}
}

func (o *Plan) ToJson() string {
	copy := *o
	b, _ := json.Marshal(&copy)
	return string(b)
}


func (o *Plan) AddProp(key string, value interface{}) {
	o.MakeNonNil()
	o.Props[key] = value
}


func PlanFromJson(data io.Reader) *Plan {
	var o *Plan
	json.NewDecoder(data).Decode(&o)
	return o
}
