package myjson

import (
	"github.com/gogf/gf/encoding/gjson"
)

type TJson = gjson.Json

func New(source ...interface{}) *TJson {
	if len(source) > 0 {
		return gjson.New(source[0])
	}
	return gjson.New(nil)
}
