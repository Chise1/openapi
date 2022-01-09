package models

type Tag struct {
	Name         string                `json:"name"`         //	REQUIRED. The name of the tag.
	Description  string                `json:"description"`  //	A description for the tag. CommonMark syntax MAY be used for rich text representation.
	ExternalDocs ExternalDocumentation `json:"externalDocs"` //	Additional external documentation for this tag.
}
