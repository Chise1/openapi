package openapi

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
)

type Handler func(context.Context, interface{}, interface{}) (interface{}, interface{})

type TestRouter struct {
	Description string
	Method      string
	Path        string
	Param       interface{}
	ResParam    interface{}
	ReqStruct   interface{}         //请求参数结构体,不能为指针，这里传入的参数会作为example
	RepStruct   interface{}         //返回正确处理的结构体
	ErrStruct   map[int]interface{} //返回报错结构体 status:struct
	Endpoint    Handler
}

func (n *TestRouter) GetReqPara() interface{} {
	return n.Param
}
func (n *TestRouter) GetReqBody() interface{} {
	return n.ReqStruct
} //interface可以为IContentType

func (n *TestRouter) GetResPara() interface{} {
	return n.ResParam
}

func (n *TestRouter) GetDescription() string {
	return n.Description
}

func (n *TestRouter) GetPath() string {
	return n.Path
}

func (n *TestRouter) GetMethod() string {
	return n.Method
}
func (n *TestRouter) GetResBody() map[int]interface{} {
	return map[int]interface{}{
		200: ReqStruct{
			Hello: "world",
		},
	}
}

type ReqStruct struct {
	Hello string `json:"hello"`
}
type ReqParam struct {
	Name   string  `json:"name" in:"path"`
	Old    uint    `json:"old"`
	Old2   uint    `json:"old2" in:"query"`
	Cook   float64 `json:"cook" in:"cookie"`
	Header string  `json:"header" in:"header"`
}

func TestRegister2Openapi(t *testing.T) {
	route := TestRouter{
		Method:      "GET",
		Description: "test",
		Path:        "/hello/{name}",
		Param:       ReqParam{Name: "张三"},
		ReqStruct:   ReqStruct{Hello: "HJ"},
		RepStruct:   ReqStruct{Hello: "HJ-b"},
		Endpoint: func(context context.Context, i interface{}, i2 interface{}) (interface{}, interface{}) {
			a := i2.(*ReqStruct)
			a.Hello = i.(*ReqParam).Name
			return nil, a
		},
	}
	Register2Openapi(&route)
	r, _ := json.Marshal(OPENAPI)
	fmt.Println(string(r))
}
