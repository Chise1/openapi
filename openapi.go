package openapi

func MakeModelComponents(model RouteStruct) {
	schemas := NewOpenapiRequest(model)
	for name, childSchema := range schemas.Components {
		OPENAPI.Components.Schemas[name] = childSchema
	}
}
