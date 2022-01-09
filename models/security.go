package models

type SecurityRequirementObject struct {
	names map[string][]string //Each name MUST correspond to a security scheme which is declared in the Security Schemes under the Components Object. If the security scheme is of type "oauth2" or "openIdConnect", then the value is a list of scope names required for the execution, and the list MAY be empty if authorization does not require a specified scope. For other security scheme types, the array MAY contain a list of role names which are required for the execution, but are not otherwise defined or exchanged in-band.
}

type SecuritySchemeType string

const (
	ApiKey        = "apiKey"
	Http          = "http"
	Oauth2        = "oauth2"
	Openidconnect = "openIdConnect"
)

type SecurityBase struct {
	Type        SecuritySchemeType `json:"type"`
	Description string             `json:"description"`
}
