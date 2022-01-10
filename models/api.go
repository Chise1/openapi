package models

type APIKeyIn string

const (
	APIKeyInQuery  = "query"
	APIKeyInHeader = "header"
	APIKeyInCookie = "cookie"
)

type APIKey struct {
	Type        SecuritySchemeType `json:"type"` //默认值为APIKeyInQuery
	In          APIKeyIn           `json:"in"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
}

func NewAPIKey() APIKey {
	return APIKey{
		Type: ApiKey,
	}
}
