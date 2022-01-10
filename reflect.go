package openapi

// schema的生成器,参考项目：https://github.com/alecthomas/jsonschema
// Package jsonschema uses reflection to generate JSON Schemas from Go types [1].
//
// If json tags are present on struct fields, they will be used to infer
// property names and if a property is required (omitempty is present).
//
// [1] http://json-schema.org/latest/json-schema-validation.html

import (
	"encoding/json"
	"github.com/Chise1/openapi/models"
	"github.com/iancoleman/orderedmap"
	"math"
	"net"
	"net/url"
	"reflect"
	"strings"
	"time"
)

// Version is the JSON SchemaChild version.
// If extending JSON SchemaChild with custom values use a custom URI.
// RFC draft-wright-json-schema-00, section 6
var Version = "http://json-schema.org/draft-04/schema#"

// SchemaChild is the root schema.
// RFC draft-wright-json-schema-00, section 4.5
type SchemaChild struct {
	*models.Schema             //根schema
	Components     Definitions //Body,包含content forms,files
}

// customSchemaType is used to detect if the type provides it's own
// custom SchemaChild Schema definition to use instead. Very useful for situations
// where there are custom JSON Marshal and Unmarshal methods.
type customSchemaType interface {
	JSONSchemaType() *models.Schema
}

var customType = reflect.TypeOf((*customSchemaType)(nil)).Elem()

// customSchemaGetFieldDocString
type customSchemaGetFieldDocString interface {
	GetFieldDocString(fieldName string) string
}

type customGetFieldDocString func(fieldName string) string

var customStructGetFieldDocString = reflect.TypeOf((*customSchemaGetFieldDocString)(nil)).Elem()

// Reflect represents a JSON SchemaChild object type.
// Reflect reflects to SchemaChild from a value using the default Reflector
func Reflect(v interface{}) *SchemaChild {
	return ReflectFromType(reflect.TypeOf(v))
}

// ReflectFromType generates root schema using the default Reflector
func ReflectFromType(t reflect.Type) *SchemaChild {
	r := &Reflector{}
	return r.ReflectFromType(t)
}

// A Reflector reflects values into a SchemaChild.
type Reflector struct {
	// AllowAdditionalProperties will cause the Reflector to generate a schema
	// with additionalProperties to 'true' for all struct types. This means
	// the presence of additional keys in JSON objects will not cause validation
	// to fail. Note said additional keys will simply be dropped when the
	// validated JSON is unmarshalled.
	AllowAdditionalProperties bool

	// RequiredFromJSONSchemaTags will cause the Reflector to generate a schema
	// that requires any key tagged with `jsonschema:required`, overriding the
	// default of requiring any key *not* tagged with `json:,omitempty`.
	RequiredFromJSONSchemaTags bool

	// YAMLEmbeddedStructs will cause the Reflector to generate a schema that does
	// not inline embedded structs. This should be enabled if the JSON schemas are
	// used with yaml.Marshal/Unmarshal.
	YAMLEmbeddedStructs bool

	// Prefer yaml: tags over json: tags to generate the schema even if json: tags
	// are present
	PreferYAMLSchema bool

	// ExpandedStruct will cause the toplevel definitions of the schema not
	// be referenced itself to a definition.
	//把根打散
	ExpandedStruct bool

	// Do not reference definitions.
	// All types are still registered under the "definitions" top-level object,
	// but instead of $ref fields in containing types, the entire definition
	// of the contained type is inserted.
	// This will cause the entire structure of types to be output in one tree.
	DoNotReference bool

	// Use package paths as well as type names, to avoid conflicts.
	// Without this setting, if two packages contain a type with the same name,
	// and both are present in a schema, they will conflict and overwrite in
	// the definition map and produce bad output.  This is particularly
	// noticeable when using DoNotReference.
	FullyQualifyTypeNames bool

	// IgnoredTypes defines a slice of types that should be ignored in the schema,
	// switching to just allowing additional properties instead.
	IgnoredTypes []interface{}

	// TypeMapper is a function that can be used to map custom Go types to jsonschema types.
	TypeMapper func(reflect.Type) *models.Schema

	// TypeNamer allows customizing of type names
	TypeNamer func(reflect.Type) string

	// AdditionalFields allows adding structfields for a given type
	AdditionalFields func(reflect.Type) []reflect.StructField
}

// Reflect reflects to SchemaChild from a value.
func (n *Reflector) Reflect(v interface{}) *SchemaChild {
	return n.ReflectFromType(reflect.TypeOf(v))
}

// ReflectFromType generates root schema
func (n *Reflector) ReflectFromType(t reflect.Type) *SchemaChild {
	components := Definitions{}
	if n.ExpandedStruct {
		st := &models.Schema{
			Version:              Version,
			Type:                 "object",
			Properties:           orderedmap.New(),
			AdditionalProperties: []byte("false"),
		}
		if n.AllowAdditionalProperties {
			st.AdditionalProperties = []byte("true")
		}
		n.reflectStructFields(st, components, t)
		n.reflectStruct(components, t)
		delete(components, n.TypeName(t))
		return &SchemaChild{Schema: st, Components: components}
	}

	s := &SchemaChild{
		Schema:     n.reflectTypeToSchema(components, t),
		Components: components,
	}
	return s
}

// Definitions hold schema definitions.
// http://json-schema.org/latest/json-schema-validation.html#rfc.section.5.26
// RFC draft-wright-json-schema-validation-00, section 5.26
//body
type Definitions map[string]*models.Schema

// Available Go defined types for JSON SchemaChild Validation.
// RFC draft-wright-json-schema-validation-00, section 7.3
var (
	timeType = reflect.TypeOf(time.Time{}) // date-time RFC section 7.3.1
	ipType   = reflect.TypeOf(net.IP{})    // ipv4 and ipv6 RFC section 7.3.4, 7.3.5
	uriType  = reflect.TypeOf(url.URL{})   // uri RFC section 7.3.6
)

// Byte slices will be encoded as base64
var byteSliceType = reflect.TypeOf([]byte(nil))

// Except for json.RawMessage
var rawMessageType = reflect.TypeOf(json.RawMessage{})

// Go code generated from protobuf enum types should fulfil this interface.
type protoEnum interface {
	EnumDescriptor() ([]byte, []int)
}

var protoEnumType = reflect.TypeOf((*protoEnum)(nil)).Elem()

func (n *Reflector) reflectTypeToSchema(definitions Definitions, t reflect.Type) *models.Schema {
	// Already added to definitions?
	if _, ok := definitions[n.TypeName(t)]; ok && !n.DoNotReference {
		return &models.Schema{Ref: models.REF_PREFIX + n.TypeName(t)}
	}

	if n.TypeMapper != nil {
		if t := n.TypeMapper(t); t != nil {
			return t
		}
	}

	if rt := n.reflectCustomType(definitions, t); rt != nil {
		return rt
	}

	// jsonpb will marshal protobuf enum options as either strings or integers.
	// It will unmarshal either.
	if t.Implements(protoEnumType) {
		return &models.Schema{OneOf: []*models.Schema{
			{Type: "string"},
			{Type: "integer"},
		}}
	}

	// Defined format types for JSON SchemaChild Validation
	// RFC draft-wright-json-schema-validation-00, section 7.3
	// TODO email RFC section 7.3.2, hostname RFC section 7.3.3, uriref RFC section 7.3.7
	if t == ipType {
		// TODO differentiate ipv4 and ipv6 RFC section 7.3.4, 7.3.5
		return &models.Schema{Type: "string", Format: "ipv4", Title: n.TypeName(t)} // ipv4 RFC section 7.3.4
	}

	switch t.Kind() {
	case reflect.Struct:
		switch t {
		case timeType: // date-time RFC section 7.3.1
			return &models.Schema{Type: "string", Format: "date-time", Title: n.TypeName(t)}
		case uriType: // uri RFC section 7.3.6
			return &models.Schema{Type: "string", Format: "uri", Title: n.TypeName(t)}
		default:
			return n.reflectStruct(definitions, t)
		}

	case reflect.Map:
		switch t.Key().Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			rt := &models.Schema{
				Type: "object",
				PatternProperties: map[string]*models.Schema{
					"^[0-9]+$": n.reflectTypeToSchema(definitions, t.Elem()),
				},
				AdditionalProperties: []byte("false"),
			}
			return rt
		}

		rt := &models.Schema{
			Type: "object",
			PatternProperties: map[string]*models.Schema{
				".*": n.reflectTypeToSchema(definitions, t.Elem()),
			},
		}
		delete(rt.PatternProperties, "additionalProperties")
		return rt

	case reflect.Slice, reflect.Array:
		returnType := &models.Schema{}
		if t == rawMessageType {
			return &models.Schema{
				AdditionalProperties: []byte("true"),
			}
		}
		if t.Kind() == reflect.Array {
			returnType.MinItems = uint64(t.Len())
			returnType.MaxItems = returnType.MinItems
		}
		if t.Kind() == reflect.Slice && t.Elem() == byteSliceType.Elem() {
			returnType.Type = "string"
			returnType.Media = &models.Schema{BinaryEncoding: "base64"}
			return returnType
		}
		returnType.Type = "array"
		returnType.Items = n.reflectTypeToSchema(definitions, t.Elem())
		return returnType

	case reflect.Interface:
		return &models.Schema{
			AdditionalProperties: []byte("true"),
		}
	case reflect.Int, reflect.Int32:
		minimum := float64(math.MinInt32)
		maximum := float64(math.MaxInt32)
		return &models.Schema{Type: "integer", Minimum: &minimum, Maximum: &maximum}
	case reflect.Int8:
		minimum := float64(math.MinInt8)
		maximum := float64(math.MaxInt8)
		return &models.Schema{Type: "integer", Minimum: &minimum, Maximum: &maximum}
	case reflect.Int16:
		minimum := float64(math.MinInt16)
		maximum := float64(math.MaxInt16)
		return &models.Schema{Type: "integer", Minimum: &minimum, Maximum: &maximum}
	case reflect.Int64:
		minimum := float64(math.MinInt64)
		maximum := float64(math.MaxInt64)
		return &models.Schema{Type: "integer", Minimum: &minimum, Maximum: &maximum}
	case reflect.Uint, reflect.Uint32:
		minimum := float64(0)
		maximum := float64(math.MaxUint32)
		return &models.Schema{Type: "integer", Minimum: &minimum, Maximum: &maximum}
	case reflect.Uint8:
		minimum := float64(0)
		maximum := float64(math.MaxUint8)
		return &models.Schema{Type: "integer", Minimum: &minimum, Maximum: &maximum}
	case reflect.Uint16:
		minimum := float64(0)
		maximum := float64(math.MaxUint16)
		return &models.Schema{Type: "integer", Minimum: &minimum, Maximum: &maximum}
	case reflect.Uint64:
		minimum := float64(0)
		maximum := float64(math.MaxUint64)
		return &models.Schema{Type: "integer", Minimum: &minimum, Maximum: &maximum}
	case reflect.Float32:
		max := math.MaxFloat32
		min := -math.MaxFloat32
		return &models.Schema{Type: "number", Minimum: &min, Maximum: &max}
	case reflect.Float64:
		max := math.MaxFloat64
		min := -math.MaxFloat64
		return &models.Schema{Type: "number", Minimum: &min, Maximum: &max}
	case reflect.Bool:
		return &models.Schema{Type: "boolean"}
	case reflect.String:
		return &models.Schema{Type: "string"}
	case reflect.Ptr:
		return n.reflectTypeToSchema(definitions, t.Elem())
	}
	panic("unsupported type " + t.String())
}

func (n *Reflector) reflectCustomType(definitions Definitions, t reflect.Type) *models.Schema {
	if t.Kind() == reflect.Ptr {
		return n.reflectCustomType(definitions, t.Elem())
	}

	if t.Implements(customType) {
		v := reflect.New(t)
		o := v.Interface().(customSchemaType)
		st := o.JSONSchemaType()
		definitions[n.TypeName(t)] = st
		if n.DoNotReference {
			return st
		} else {
			return &models.Schema{
				Version: Version,
				Ref:     models.REF_PREFIX + n.TypeName(t),
			}
		}
	}

	return nil
}

// Reflects a struct to a JSON SchemaChild type.
func (n *Reflector) reflectStruct(definitions Definitions, t reflect.Type) *models.Schema {
	if st := n.reflectCustomType(definitions, t); st != nil {
		return st
	}

	for _, ignored := range n.IgnoredTypes {
		if reflect.TypeOf(ignored) == t {
			st := &models.Schema{
				Type:                 "object",
				Properties:           orderedmap.New(),
				AdditionalProperties: []byte("true"),
			}
			definitions[n.TypeName(t)] = st

			if n.DoNotReference {
				return st
			} else {
				return &models.Schema{
					Version: Version,
					Ref:     models.REF_PREFIX + n.TypeName(t),
				}
			}
		}
	}

	st := &models.Schema{
		Type:                 "object",
		Properties:           orderedmap.New(),
		AdditionalProperties: []byte("false"),
		Title:                n.TypeName(t),
	}
	if n.AllowAdditionalProperties {
		st.AdditionalProperties = []byte("true")
	}
	definitions[n.TypeName(t)] = st
	n.reflectStructFields(st, definitions, t)

	if n.DoNotReference {
		return st
	} else {
		return &models.Schema{
			Version: Version,
			Ref:     models.REF_PREFIX + n.TypeName(t),
		}
	}
}

func (n *Reflector) reflectStructFields(st *models.Schema, definitions Definitions, t reflect.Type) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return
	}

	var getFieldDocString customGetFieldDocString
	if t.Implements(customStructGetFieldDocString) {
		v := reflect.New(t)
		o := v.Interface().(customSchemaGetFieldDocString)
		getFieldDocString = o.GetFieldDocString
	}

	handleField := func(f reflect.StructField) {
		name, shouldEmbed, required, nullable := n.reflectFieldName(f)
		// if anonymous and exported type should be processed recursively
		// current type should inherit properties of anonymous one
		if name == "" {
			if shouldEmbed {
				n.reflectStructFields(st, definitions, f.Type)
			}
			return
		}

		property := n.reflectTypeToSchema(definitions, f.Type)
		property.StructKeywordsFromTags(f, st, name)
		if getFieldDocString != nil {
			property.Description = getFieldDocString(f.Name)
		}

		if nullable {
			property = &models.Schema{
				OneOf: []*models.Schema{
					property,
					{
						Type: "null",
					},
				},
			}
		}

		st.Properties.Set(name, property)
		if required {
			st.Required = append(st.Required, name)
		}
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		handleField(f)
	}
	if n.AdditionalFields != nil {
		if af := n.AdditionalFields(t); af != nil {
			for _, sf := range af {
				handleField(sf)
			}
		}
	}
}

func requiredFromJSONTags(tags []string) bool {
	if ignoredByJSONTags(tags) {
		return false
	}

	for _, tag := range tags[1:] {
		if tag == "omitempty" {
			return false
		}
	}
	return true
}

func requiredFromJSONSchemaTags(tags []string) bool {
	if ignoredByJSONSchemaTags(tags) {
		return false
	}
	for _, tag := range tags {
		if tag == "required" {
			return true
		}
	}
	return false
}

func nullableFromJSONSchemaTags(tags []string) bool {
	if ignoredByJSONSchemaTags(tags) {
		return false
	}
	for _, tag := range tags {
		if tag == "nullable" {
			return true
		}
	}
	return false
}

func inlineYAMLTags(tags []string) bool {
	for _, tag := range tags {
		if tag == "inline" {
			return true
		}
	}
	return false
}

func ignoredByJSONTags(tags []string) bool {
	return tags[0] == "-"
}

func ignoredByJSONSchemaTags(tags []string) bool {
	return tags[0] == "-"
}

func (n *Reflector) reflectFieldName(f reflect.StructField) (string, bool, bool, bool) {
	jsonTags, exist := f.Tag.Lookup("json")
	yamlTags, yamlExist := f.Tag.Lookup("yaml")
	if !exist || n.PreferYAMLSchema {
		jsonTags = yamlTags
		exist = yamlExist
	}

	jsonTagsList := strings.Split(jsonTags, ",")
	yamlTagsList := strings.Split(yamlTags, ",")

	if ignoredByJSONTags(jsonTagsList) {
		return "", false, false, false
	}

	jsonSchemaTags := strings.Split(f.Tag.Get(models.TagName), ",")
	if ignoredByJSONSchemaTags(jsonSchemaTags) {
		return "", false, false, false
	}

	name := f.Name
	required := requiredFromJSONTags(jsonTagsList)

	if n.RequiredFromJSONSchemaTags {
		required = requiredFromJSONSchemaTags(jsonSchemaTags)
	}

	nullable := nullableFromJSONSchemaTags(jsonSchemaTags)

	if jsonTagsList[0] != "" {
		name = jsonTagsList[0]
	}

	// field not anonymous and not export has no export name
	if !f.Anonymous && f.PkgPath != "" {
		name = ""
	}

	embed := false

	// field anonymous but without json tag should be inherited by current type
	if f.Anonymous && !exist {
		if !n.YAMLEmbeddedStructs {
			name = ""
			embed = true
		} else {
			name = strings.ToLower(name)
		}
	}

	if yamlExist && inlineYAMLTags(yamlTagsList) {
		name = ""
		embed = true
	}

	return name, embed, required, nullable
}

func (n *SchemaChild) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(n.Schema)
	if err != nil {
		return nil, err
	}
	if n.Components == nil || len(n.Components) == 0 {
		return b, nil
	}
	d, err := json.Marshal(struct {
		Components Definitions `json:"components,omitempty"`
	}{n.Components})
	if err != nil {
		return nil, err
	}
	if len(b) == 2 {
		return d, nil
	} else {
		b[len(b)-1] = ','
		return append(b, d[1:]...), nil
	}
}

func (n *Reflector) TypeName(t reflect.Type) string {
	if n.TypeNamer != nil {
		if name := n.TypeNamer(t); name != "" {
			return name
		}
	}
	if n.FullyQualifyTypeNames {
		return t.PkgPath() + "." + t.Name()
	}
	return t.Name()
}
