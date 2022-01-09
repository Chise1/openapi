package models

type OAuthFlowImplicit struct {
	RefreshUrl       string            `json:"refreshUrl"`
	Scopes           map[string]string `json:"scopes"`
	AuthorizationUrl string            `json:"authorizationUrl"`
}
type OAuthFlowPassword struct {
	RefreshUrl string            `json:"refreshUrl"`
	Scopes     map[string]string `json:"scopes"`
	TokenUrl   string            `json:"tokenUrl"`
}
type OAuthFlowClientCredentials struct {
	RefreshUrl string            `json:"refreshUrl"`
	Scopes     map[string]string `json:"scopes"`
	TokenUrl   string            `json:"tokenUrl"`
}
type OAuthFlowAuthorizationCode struct {
	RefreshUrl       string            `json:"refreshUrl"`
	Scopes           map[string]string `json:"scopes"`
	AuthorizationUrl string            `json:"authorizationUrl"`
	TokenUrl         string            `json:"tokenUrl"`
}
type OAuthFlows struct {
	Implicit          *OAuthFlowImplicit          `json:"implicit"`
	Password          *OAuthFlowPassword          `json:"password"`
	ClientCredentials *OAuthFlowClientCredentials `json:"clientCredentials"`
	AuthorizationCode *OAuthFlowAuthorizationCode `json:"authorizationCode"`
}

type OAuth2 struct {
	Type        SecuritySchemeType `json:"type"`
	Description string             `json:"description"`
	Flows       OAuthFlows         `json:"flows"`
}

func NewOAuth2() OAuth2 {
	return OAuth2{Type: Oauth2}
}

type OpenIdConnect struct {
	Type        SecuritySchemeType `json:"type"`
	Description string             `json:"description"`
	Flows       OAuthFlows         `json:"flows"`
}

func NewOpenIdConnect() OAuth2 {
	return OAuth2{Type: Openidconnect}
}
