package models

type Header struct {
	Description   string               `json:"description"` //A brief description of the parameter. This could contain examples of use. CommonMark syntax MAY be used for rich text representation.
	Required      bool                 `json:"required"`    //Determines whether this parameter is mandatory. If the parameter location is "path", this property is REQUIRED and its value MUST be true. Otherwise, the property MAY be included and its default value is false.
	Deprecated    bool                 `json:"deprecated"`  //Specifies that a parameter is deprecated and SHOULD be transitioned out of usage. Default value is false.
	Style         string               `json:"style"`
	Explode       bool                 `json:"explode"`
	AllowReserved bool                 `json:"allowReserved"`
	Schema        *Schema              `json:"schema"`
	Example       interface{}          `json:"example"`
	Examples      map[string]Example   `json:"examples"`
	Content       map[string]MediaType `json:"content"`
	*Reference
}
