package openapi

import (
	"github.com/Chise1/openapi/models"
	"github.com/iancoleman/orderedmap"
	"reflect"
	"strconv"
	"strings"
)

// Parameters 路径参数
type Parameters []*models.Parameter
type RouterHelper struct {
	Body       *models.RequestBody //这才是真正应该有的body
	Components Definitions         //Body,包含content forms,files
	Parameters Parameters          //path,query,header,cookie, todo path之后在处理
	Response   map[string]*models.Response
	*models.Schema
	ReqContentType string
}

func NewOpenapiRequest(v RouteStruct) *RouterHelper {
	r := &Reflector{}
	n := &RouterHelper{}
	reqBody := v.GetReqBody()
	if reqBody != nil {
		body, ok := reqBody.(IContentType)
		if ok {
			n.reqBody(r, body, body.GetContentType())
		} else {
			n.reqBody(r, reqBody, Json)
		}
	}
	reqpara := v.GetReqPara()
	if reqpara != nil {
		n.para(r, v.GetReqPara())
	}

	n.Response = map[string]*models.Response{}
	n.WriteRes(v.GetResBody())
	return n
}
func (n *RouterHelper) GetSchema() *models.Schema {
	if n.Schema.Ref != "" {
		return n.GetSchemaStruct(n.Schema)
	}
	return n.Schema
}
func (n *RouterHelper) WriteRes(exceptRes map[int]interface{}) {
	if exceptRes == nil {
		return
	}
	for status, res := range exceptRes {
		if res == nil {
			continue
		}
		exceptSchema := Reflect(res)
		n.updateComponents(exceptSchema.Components)
		body, ok := res.(IContentType)
		t := ""
		if ok {
			t = string(body.GetContentType())
		} else {
			t = "application/json"
		}
		n.Response[strconv.Itoa(status)] = &models.Response{
			Content: map[string]*models.MediaType{
				t: {
					Schema:  exceptSchema.Schema,
					Example: res,
				},
			},
		}
	}
}
func (n *RouterHelper) reqBody(reflector *Reflector, v interface{}, reqType ContentType) {
	body := n.body(reflector, v, reqType)
	for k, v := range body {
		n.ReqContentType = k
		n.Schema = v.Schema
	}
	n.Body = &models.RequestBody{
		Content: body,
	}
}

func (n *RouterHelper) updateComponents(components map[string]*models.Schema) {
	if n.Components == nil {
		n.Components = map[string]*models.Schema{}
	}
	for name, schema := range components {
		n.Components[name] = schema
	}
}
func (n *RouterHelper) body(reflector *Reflector, v interface{}, reqType ContentType) map[string]*models.MediaType {
	var components Definitions = map[string]*models.Schema{}
	vType := reflect.TypeOf(v)
	schema := reflector.reflectTypeToSchema(components, vType)
	n.updateComponents(components)
	rootSchema := n.GetSchemaStruct(schema)
	return map[string]*models.MediaType{
		string(reqType): {
			Schema:  schema,
			Example: GetExample(rootSchema, v),
		},
	}
}
func (n *RouterHelper) para(reflector *Reflector, v interface{}) {
	components := Definitions{}
	st := &models.Schema{
		Version:              Version,
		Type:                 "object",
		Properties:           orderedmap.New(),
		AdditionalProperties: []byte("false"),
	}
	t := reflect.TypeOf(v)
	reflector.reflectStructFields(st, components, t)
	reflector.reflectStruct(components, t)
	delete(components, reflector.TypeName(t))
	for _, name := range st.Properties.Keys() {
		iproperty, _ := st.Properties.Get(name)
		property := iproperty.(*models.Schema)
		field, _ := t.FieldByName(property.FieldName)
		in := ""
		switch field.Tag.Get("in") {
		case "path":
			in = "path"
		case "query":
			in = "query"
		case "header":
			in = "header"
		case "cookie":
			in = "cookie"
		default:
			in = "query"
		}
		required := false
		for i := range st.Required {
			if st.Required[i] == name {
				required = true
				break
			}
		}
		n.Parameters = append(n.Parameters, &models.Parameter{
			Schema:      property,
			Name:        name,
			In:          in,
			Description: property.Description,
			Required:    required,
			Example:     v,
		})
	}
}

func (n *RouterHelper) GetSchemaStruct(schema *models.Schema) *models.Schema {
	schemaPath := strings.Split(schema.Ref, "/")
	s, ok := n.Components[schemaPath[len(schemaPath)-1]]
	if !ok {
		return nil
	}
	return s
}
func GetExample(schema *models.Schema, v interface{}) map[string]interface{} {
	res := map[string]interface{}{}
	value := reflect.ValueOf(v)
	for _, key := range schema.Properties.Keys() {
		prop, _ := schema.Properties.Get(key)
		fieldSchema := prop.(*models.Schema)
		res[fieldSchema.FieldName] = value.FieldByName(fieldSchema.FieldName).Interface()
	}
	return res
}
