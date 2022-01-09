package models

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
)

func (n *Schema) MarshalJSON() ([]byte, error) {
	type Type_ Schema
	b, err := json.Marshal((*Type_)(n))
	if err != nil {
		return nil, err
	}
	if n.Extras == nil || len(n.Extras) == 0 {
		return b, nil
	}
	m, err := json.Marshal(n.Extras)
	if err != nil {
		return nil, err
	}
	if len(b) == 2 {
		return m, nil
	} else {
		b[len(b)-1] = ','
		return append(b, m[1:]...), nil
	}
}

// StructKeywordsFromTags parse structfield's tag
func (n *Schema) StructKeywordsFromTags(f reflect.StructField, parentType *Schema, propertyName string) {
	n.Description = f.Tag.Get(Desc)
	tags := strings.Split(f.Tag.Get(TagName), ",")
	n.genericKeywords(tags, parentType, propertyName)
	//default title
	if n.Title == "" {
		n.Title = f.Name
	}
	n.FieldName = f.Name //deprated
	switch n.Type {
	case "string":
		n.stringKeywords(tags)
	case "number":
		n.numbericKeywords(tags)
	case "integer":
		n.numbericKeywords(tags)
	case "array":
		n.arrayKeywords(tags)
	}
	extras := strings.Split(f.Tag.Get("jsonschema_extras"), ",")
	n.extraKeywords(extras)
}

// read struct tags for generic keyworks
func (n *Schema) genericKeywords(tags []string, parentType *Schema, propertyName string) {
	for _, tag := range tags {
		nameValue := strings.Split(tag, "=")
		if len(nameValue) == 2 {
			name, val := nameValue[0], nameValue[1]
			switch name {
			case "title":
				n.Title = val
			case "description":
				n.Description = val
			case "type":
				n.Type = val
			case "oneof_required":
				var typeFound *Schema
				for i := range parentType.OneOf {
					if parentType.OneOf[i].Title == nameValue[1] {
						typeFound = parentType.OneOf[i]
					}
				}
				if typeFound == nil {
					typeFound = &Schema{
						Title:    nameValue[1],
						Required: []string{},
					}
					parentType.OneOf = append(parentType.OneOf, typeFound)
				}
				typeFound.Required = append(typeFound.Required, propertyName)
			case "oneof_type":
				if n.OneOf == nil {
					n.OneOf = make([]*Schema, 0, 1)
				}
				n.Type = ""
				types := strings.Split(nameValue[1], ";")
				for _, ty := range types {
					n.OneOf = append(n.OneOf, &Schema{
						Type: ty,
					})
				}
			case "enum":
				switch n.Type {
				case "string":
					n.Enum = append(n.Enum, val)
				case "integer":
					i, _ := strconv.Atoi(val)
					n.Enum = append(n.Enum, i)
				case "number":
					f, _ := strconv.ParseFloat(val, 64)
					n.Enum = append(n.Enum, f)
				}
			}
		}
	}
}

// read struct tags for string type keyworks
func (n *Schema) stringKeywords(tags []string) {
	for _, tag := range tags {
		nameValue := strings.Split(tag, "=")
		if len(nameValue) == 2 {
			name, val := nameValue[0], nameValue[1]
			switch name {
			case "minLength":
				i, _ := strconv.Atoi(val)
				n.MinLength = i
			case "maxLength":
				i, _ := strconv.Atoi(val)
				n.MaxLength = i
			case "pattern":
				n.Pattern = val
			case "format":
				switch val {
				case "date-time", "email", "hostname", "ipv4", "ipv6", "uri":
					n.Format = val
					break
				}
			case "default":
				n.Default = val
			case "example":
				n.Examples = append(n.Examples, val)
			}
		}
	}
}

// read struct tags for numberic type keyworks
func (n *Schema) numbericKeywords(tags []string) {
	for _, tag := range tags {
		nameValue := strings.Split(tag, "=")
		if len(nameValue) == 2 {
			name, val := nameValue[0], nameValue[1]
			switch name {
			case "multipleOf":
				i, _ := strconv.ParseFloat(val, 32)
				n.MultipleOf = i
			case "minimum":
				i, _ := strconv.ParseFloat(val, 32)
				n.Minimum = &i
			case "maximum":
				i, _ := strconv.ParseFloat(val, 32)
				n.Maximum = &i
			case "exclusiveMaximum":
				b, _ := strconv.ParseBool(val)
				n.ExclusiveMaximum = b
			case "gte":
				i, _ := strconv.ParseFloat(val, 32)
				n.Minimum = &i
				b, _ := strconv.ParseBool(val)
				n.ExclusiveMinimum = b
			case "exclusiveMinimum":
				b, _ := strconv.ParseBool(val)
				n.ExclusiveMinimum = b
			case "default":
				i, _ := strconv.Atoi(val)
				n.Default = i
			case "example":
				if i, err := strconv.Atoi(val); err == nil {
					n.Examples = append(n.Examples, i)
				}
			}
		}
	}
}

// read struct tags for object type keyworks
// func (t *Schema) objectKeywords(tags []string) {
//     for _, tag := range tags{
//         nameValue := strings.Split(tag, "=")
//         name, val := nameValue[0], nameValue[1]
//         switch name{
//             case "dependencies":
//                 t.Dependencies = val
//                 break;
//             case "patternProperties":
//                 t.PatternProperties = val
//                 break;
//         }
//     }
// }

// read struct tags for array type keyworks
func (n *Schema) arrayKeywords(tags []string) {
	var defaultValues []interface{}
	for _, tag := range tags {
		nameValue := strings.Split(tag, "=")
		if len(nameValue) == 2 {
			name, val := nameValue[0], nameValue[1]
			switch name {
			case "minItems":
				i, _ := strconv.ParseUint(val, 10, 32)
				n.MinItems = i
			case "maxItems":
				i, _ := strconv.ParseUint(val, 10, 32)
				n.MaxItems = i
			case "uniqueItems":
				n.UniqueItems = true
			case "default":
				defaultValues = append(defaultValues, val)
			case "enum":
				switch n.Items.Type {
				case "string":
					n.Items.Enum = append(n.Items.Enum, val)
				case "integer":
					i, _ := strconv.Atoi(val)
					n.Items.Enum = append(n.Items.Enum, i)
				case "number":
					f, _ := strconv.ParseFloat(val, 64)
					n.Items.Enum = append(n.Items.Enum, f)
				}
			}
		}
	}
	if len(defaultValues) > 0 {
		n.Default = defaultValues
	}
}

func (n *Schema) extraKeywords(tags []string) {
	for _, tag := range tags {
		nameValue := strings.Split(tag, "=")
		if len(nameValue) == 2 {
			n.setExtra(nameValue[0], nameValue[1])
		}
	}
}

func (n *Schema) setExtra(key, val string) {
	if n.Extras == nil {
		n.Extras = map[string]interface{}{}
	}
	if existingVal, ok := n.Extras[key]; ok {
		switch existingVal := existingVal.(type) {
		case string:
			n.Extras[key] = []string{existingVal, val}
		case []string:
			n.Extras[key] = append(existingVal, val)
		case int:
			n.Extras[key], _ = strconv.Atoi(val)
		}
	} else {
		n.Extras[key] = val
	}
}
