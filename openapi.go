package openapi

func MakeModelComponents(model RouteStruct) *RouterHelper {
	schemas := NewOpenapiRequest(model)
	for name, childSchema := range schemas.Components {
		OPENAPI.Components.Schemas[name] = childSchema
	}
	return schemas
}
