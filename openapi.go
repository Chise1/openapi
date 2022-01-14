package openapi

import (
	"github.com/Chise1/openapi/models"
	"reflect"
)

func Register2Openapi(n RouteStruct) *RouterHelper {
	schemas := NewOpenapiRequest(n)
	for name, childSchema := range schemas.Components {
		OPENAPI.Components.Schemas[name] = childSchema
	}
	apiRef := reflect.TypeOf(n.GetReqBody())
	endpointName := reflect.TypeOf(n).Name()
	pathItem := &models.PathItem{}
	method := n.GetMethod()
	if method == "" {
		method = "GET"
	}
	path := n.GetPath()
	oper := &models.Operation{
		Description: n.GetDescription(),
		Summary:     "",
		OperationId: apiRef.PkgPath() + apiRef.Name() + endpointName + method,
		RequestBody: schemas.Body,
		Parameters:  schemas.Parameters,
	}
	oper.Responses = map[string]*models.Response{"200": &DefaultRes}
	oper.Responses = schemas.Response
	if method == "GET" {
		pathItem.Get = oper
	} else if method == "POST" {
		pathItem.Post = oper
	} else if method == "PUT" {
		pathItem.Put = oper
	} else if method == "DELETE" {
		pathItem.Delete = oper
	} //todo 需要支持其他method.
	OPENAPI.Paths[path] = pathItem
	return schemas
}
