package models

type Contact struct {
	Name  string `json:"name,omitempty"`  //The identifying name of the contact person/organization.
	Url   string `json:"url,omitempty"`   //The URL pointing to the contact information. This MUST be in the form of a URL.
	Email string `json:"email,omitempty"` //The email address of the contact person/organization. This MUST be in the form of an email address.
}

type License struct {
	Name string `json:"name,omitempty"` //REQUIRED. The license name used for the API.
	Url  string `json:"url,omitempty"`  //A URL to the license used for the API. This MUST be in the form of a URL. The url field is mutually exclusive of the identifier field.
}

type Info struct {
	Title          string   `json:"title,omitempty"`          //REQUIRED. The title of the API.
	Summary        string   `json:"summary,omitempty"`        //A short summary of the API.
	Description    string   `json:"description,omitempty"`    //A description of the API. CommonMark syntax MAY be used for rich text representation.
	TermsOfService string   `json:"termsOfService,omitempty"` //A URL to the Terms of Service for the API. This MUST be in the form of a URL.
	Contact        *Contact `json:"contact,omitempty"`        //The contact information for the exposed API.
	License        *License `json:"license,omitempty"`        //The license information for the exposed API.
	Version        string   `json:"version,omitempty"`        //REQUIRED. The version of the OpenAPI document (which is distinct from the OpenAPI Specification version or the API implementation version).
}
