package models

type HTTPBase struct {
	Type        SecuritySchemeType `json:"type"`
	Description string             `json:"description"`
	Schema      string             `json:"schema"`
}

type HTTPBearer struct {
	Type         SecuritySchemeType `json:"type"`
	Description  string             `json:"description"`
	Schema       string             `json:"schema"`
	BearerFormat string             `json:"bearerFormat"`
}

func NewHTTPBearer() HTTPBearer {
	return HTTPBearer{
		Schema: "bearer",
		Type:   Http,
	}
}
