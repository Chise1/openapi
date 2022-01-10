package models

type IOpenapiStruct interface {
	Map() map[string]interface{}
}
type IReference interface {
	GetRef() struct{}
}

const (
	REF_PREFIX  = "#/components/schemas/"
	TagName     = "openapi"
	Description = "openapi_desc"   //把description单独写
	Extras      = "openapi_extras" //额外的字段，和tagname一致
)
