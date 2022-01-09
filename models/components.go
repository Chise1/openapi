package models

type Components struct {
	Schemas         map[string]*Schema      `json:"schemas,omitempty"`
	Responses       map[string]*Response    `json:"responses,omitempty"`       //    An object to hold reusable Response Objects.
	Parameters      map[string]*Parameter   `json:"parameters,omitempty"`      //  An object to hold reusable Parameter Objects.
	Examples        map[string]*Example     `json:"examples,omitempty"`        //   An object to hold reusable Example Objects.
	RequestBodies   map[string]*RequestBody `json:"requestBodies,omitempty"`   //    An object to hold reusable Request Body Objects.
	Headers         map[string]*Header      `json:"headers,omitempty"`         //    An object to hold reusable Header Objects.
	SecuritySchemes map[string]interface{}  `json:"securitySchemes,omitempty"` //    An object to hold reusable Security Scheme Objects.
	Links           map[string]*Link        `json:"links,omitempty"`           //  An object to hold reusable Link Objects.
	Callbacks       map[string]*PathItem    `json:"callbacks,omitempty"`       //An object to hold reusable Callback Objects.
	PathItems       map[string]*PathItem    `json:"pathItems,omitempty"`       //    An object to hold reusable Path Item Object.
}

func (n *Components) SetSecuritySchemes(key string, value interface{}) {
	switch value.(type) {
	case APIKey:
		n.SecuritySchemes[key] = value
	case HTTPBase:
		n.SecuritySchemes[key] = value
	case OAuth2:
		n.SecuritySchemes[key] = value
	case OpenIdConnect:
		n.SecuritySchemes[key] = value
	case HTTPBearer:
		n.SecuritySchemes[key] = value
	default:
		panic("error value!")
	}
}
