package openapi

import (
	"github.com/Chise1/openapi/models"
	"reflect"
)

func struct2Example(schema *models.Schema, v interface{}) map[string]interface{} {
	res := map[string]interface{}{}
	value := reflect.ValueOf(v)
	for _, key := range schema.Properties.Keys() {
		prop, _ := schema.Properties.Get(key)
		fieldSchema := prop.(*models.Schema)
		res[fieldSchema.FieldName] = value.FieldByName(fieldSchema.FieldName).Interface()
	}
	return res
}
