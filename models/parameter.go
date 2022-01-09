package models

type Parameter struct {
	// Parameter
	Name          string                `json:"name,omitempty"`
	In            string                `json:"in,omitempty"`          //REQUIRED. The location of the parameter. Possible values are "query", "header", "path" or "cookie".
	Description   string                `json:"description,omitempty"` //A brief description of the parameter. This could contain examples of use. CommonMark syntax MAY be used for rich text representation.
	Required      bool                  `json:"required,omitempty"`    //Determines whether this parameter is mandatory. If the parameter location is "path", this property is REQUIRED and its value MUST be true. Otherwise, the property MAY be included and its default value is false.
	Deprecated    bool                  `json:"deprecated,omitempty"`  //Specifies that a parameter is deprecated and SHOULD be transitioned out of usage. Default value is false.
	Style         string                `json:"style,omitempty"`
	Explode       bool                  `json:"explode,omitempty"`
	AllowReserved bool                  `json:"allowReserved,omitempty"`
	Schema        *Schema               `json:"schema,omitempty,omitempty"`
	Example       interface{}           `json:"example,omitempty"`
	Examples      map[string]*Example   `json:"examples,omitempty"` //暂不支持
	Content       map[string]*MediaType `json:"content,omitempty"`
	Ref           string                `json:"$ref,omitempty"` //REQUIRED. The reference identifier. This MUST be in the form of a URI.
}
