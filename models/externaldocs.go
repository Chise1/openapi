package models

type ExternalDocumentation struct {
	Description string `json:"description"` //A description of the target documentation. CommonMark syntax MAY be used for rich text representation.
	Url         string `json:"url"`         //	REQUIRED. The URL for the target documentation. This MUST be in the form of a URL.
}
