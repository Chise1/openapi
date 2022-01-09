package models

type IOpenapiStruct interface {
	Map() map[string]interface{}
}
type IReference interface {
	GetRef() struct{}
}

const (
	REF_PREFIX = "#/components/schemas/"
)
