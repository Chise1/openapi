package models

type Example struct {
	Summary       string      `json:"summary,omitempty"`       //Short description for the example.
	Description   string      `json:"description,omitempty"`   //	Long description for the example. CommonMark syntax MAY be used for rich text representation.
	Value         interface{} `json:"value,omitempty"`         // string   Embedded literal example.The value field and externalValue field are mutually exclusive.To represent examples of media types that cannot naturally represented in JSON or YAML, use a string value to contain the example, escaping where necessary.
	ExternalValue string      `json:"externalValue,omitempty"` //   A URI that points to the literal example.This provides the capability to reference examples that cannot easily be included in JSON or YAML documents.The value field and externalValue field are mutually exclusive.See the rules for resolving Relative References.
	Ref           string      `json:"$ref,omitempty"`
}

type Encoding struct {
	ContentType   string             `json:"contentType,omitempty"`   //The Content-Type for encoding a specific property. Default value depends on the property type: for object - application/json; for array – the default is defined based on the inner type; for all other cases the default is application/octet-stream. The value can be a specific media type (e.g. application/json), a wildcard media type (e.g. image/*), or a comma-separated list of the two types.
	Headers       map[string]*Header `json:"headers,omitempty"`       //[HeaderObject | Reference]    A map allowing additional information to be provided as headers, for example Content-Disposition.Content-Type is described separately and SHALL be ignored in this section.This property SHALL be ignored if the request body media type is not a multipart.
	Style         string             `json:"style,omitempty"`         //Describes how a specific property value will be serialized depending on its type.See Parameter Object for details on the style property.The behavior follows the same values as query parameters, including default values.This property SHALL be ignored if the request body media type is not application/x-www-form-urlencoded or multipart/form-data.If a value is explicitly defined, then the value of contentType (implicit or explicit) SHALL be ignored.
	Explode       bool               `json:"explode,omitempty"`       // When this is true, property values of type array or object generate separate parameters for each value of the array, or key-value-pair of the map.For other types of properties this property has no effect.When style is form, the default value is true.For all other styles, the default value is false.This property SHALL be ignored if the request body media type is not application/x-www-form-urlencoded or multipart/form-data.If a value is explicitly defined, then the value of contentType (implicit or explicit) SHALL be ignored.
	AllowReserved bool               `json:"allowReserved,omitempty"` //    Determines whether the parameter value SHOULD allow reserved characters, as defined by RFC3986:/?#[]@!$&'()*+,;= to be included without percent-encoding. The default value is false. This property SHALL be ignored if the request body media type is not application/x-www-form-urlencoded or multipart/form-data. If a value is explicitly defined, then the value of contentType (implicit or explicit) SHALL be ignored.
}

type MediaType struct {
	Schema   *Schema                `json:"schema,omitempty"` //The schema defining the content of the request, response, or parameter.
	Example  *Example               `json:"example,omitempty"`
	Examples map[string]interface{} `json:"examples,omitempty"` //[Example | Reference]Examples of the media type. Each example object SHOULD match the media type and specified schema if present. The examples field is mutually exclusive of the example field. Furthermore, if referencing a schema which contains an example, the examples value SHALL override the example provided by the schema.
	Encoding map[string]*Encoding   `json:"encoding,omitempty"` //A map between a property name and its encoding information. The key, being the property name, MUST exist in the schema as a property. The encoding object SHALL only apply to requestBody objects when the media type is multipart or application/x-www-form-urlencoded.
}

func (n *MediaType) SetExamples(key string, value interface{}) {
	switch value.(type) {
	case Example:
		n.Examples[key] = value
	case Reference:
		n.Examples[key] = value
	}
}

type RequestBody struct {
	Description string                `json:"description,omitempty"` //A brief description of the request body. This could contain examples of use. CommonMark syntax MAY be used for rich text representation.
	Content     map[string]*MediaType `json:"content,omitempty"`     //REQUIRED. The content of the request body. The key is a media type or media type range and the value describes it. For requests that match multiple keys, only the most specific key is applicable. e.g. text/plain overrides text/*
	Required    bool                  `json:"required,omitempty"`    //Determines if the request body is required in the request. Defaults to false.
	Ref         string                `json:"$ref,omitempty"`
}
type Link struct {
	OperationRef string                `json:"operationRef,omitempty"`
	OperationId  string                `json:"operationId,omitempty"`
	Parameters   map[string]*Parameter `json:"parameters,omitempty"`
	RequestBody  *RequestBody          `json:"requestBody,omitempty"`
	Description  string                `json:"description,omitempty"`
	Server       *Server               `json:"server,omitempty"`
	Ref          string                `json:"$ref,omitempty"` //REQUIRED. The reference identifier. This MUST be in the form of a URI.
}

type Response struct {
	Description string                `json:"description,omitempty"` //REQUIRED. A description of the response. CommonMark syntax MAY be used for rich text representation.
	Headers     map[string]*Header    `json:"headers,omitempty"`     //[HeaderObject | Reference Object]	Maps a header name to its definition. RFC7230 states header names are case insensitive. If a response header is defined with the name "Content-Type", it SHALL be ignored.
	Content     map[string]*MediaType `json:"content,omitempty"`     //	A map containing descriptions of potential response payloads. The key is a media type or media type range and the value describes it. For responses that match multiple keys, only the most specific key is applicable. e.g. text/plain overrides text/*
	Links       map[string]*Link      `json:"links,omitempty"`       //[LinkObject | Reference Object]	A map of operations links that can be followed from the response. The key of the map is a short name for the link, following the naming constraints of the names for Component Objects.
	Ref         string                `json:"$ref,omitempty"`        //REQUIRED. The reference identifier. This MUST be in the form of a URI.
}

type Operation struct {
	Tags         []string               `json:"tags,omitempty"`         //A list of tags for API documentation control. Tags can be used for logical grouping of operations by resources or any other qualifier.
	Summary      string                 `json:"summary,omitempty"`      //A short summary of what the operation does.
	Description  string                 `json:"description,omitempty"`  //A verbose explanation of the operation behavior. CommonMark syntax MAY be used for rich text representation.
	ExternalDocs *ExternalDocumentation `json:"externalDocs,omitempty"` //Additional external documentation for this operation.
	OperationId  string                 `json:"operationId,omitempty"`  //Unique string used to identify the operation. The id MUST be unique among all operations described in the API. The operationId value is case-sensitive. Tools and libraries MAY use the operationId to uniquely identify an operation, therefore, it is RECOMMENDED to follow common programming naming conventions.
	Parameters   []*Parameter           `json:"parameters,omitempty"`   //Parameter or Reference    //A list of parameters that are applicable for this operation. If a parameter is already defined at the Path Item, the new definition will override it but can never remove it. The list MUST NOT include duplicated parameters. A unique parameter is defined by a combination of a name and location. The list can use the Reference to link to parameters that are defined at the OpenAPI Object's components/parameters.
	RequestBody  *RequestBody           `json:"requestBody,omitempty"`  //The request body applicable for this operation. The requestBody is fully supported in HTTP methods where the HTTP 1.1 specification RFC7231 has explicitly defined semantics for request bodies. In other cases where the HTTP spec is vague (such as GET, HEAD and DELETE), requestBody is permitted but does not have well-defined semantics and SHOULD be avoided if possible.
	Responses    map[string]*Response   `json:"responses,omitempty"`
	Callbacks    map[string]*PathItem   `json:"callbacks,omitempty"`  //str?
	Deprecated   bool                   `json:"deprecated,omitempty"` //Specifies that a parameter is deprecated and SHOULD be transitioned out of usage. Default value is false.
	Security     []map[string][]string  `json:"security,omitempty"`
	Servers      []*Server              `json:"servers,omitempty"`
}
type PathItem struct {
	Ref         string     `json:"$ref,omitempty"`        //Allows for a referenced definition of this path item. The referenced structure MUST be in the form of a Path Item Object. In case a Path Item Object field appears both in the defined object and the referenced object, the behavior is undefined. See the rules for resolving Relative References.
	Summary     string     `json:"summary,omitempty"`     //An optional, string summary, intended to apply to all operations in this path.
	Description string     `json:"description,omitempty"` //An optional, string description, intended to apply to all operations in this path. CommonMark syntax MAY be used for rich text representation.
	Get         *Operation `json:"get,omitempty"`
	Put         *Operation `json:"put,omitempty"`
	Post        *Operation `json:"post,omitempty"`
	Delete      *Operation `json:"delete,omitempty"`
	Options     *Operation `json:"options,omitempty"`
	Head        *Operation `json:"head,omitempty"`
	Patch       *Operation `json:"patch,omitempty"`
	Trace       *Operation `json:"trace,omitempty"`
	Servers     []*Server  `json:"servers,omitempty"`
	Parameters  *Parameter `json:"parameters,omitempty"`
}
