package models

import (
	"encoding/json"
	"github.com/iancoleman/orderedmap"
	"reflect"
	"strconv"
	"strings"
)

type Discriminator struct {
	PropertyName string            `json:"propertyName,omitempty"` //	REQUIRED. The name of the property in the payload that will hold the discriminator value.
	Mapping      map[string]string `json:"mapping,omitempty"`      //An object to hold mappings between payload values and schema names or references.
}
type XML struct {
	Name      string `json:"name,omitempty"`      //	Replaces the name of the element/attribute used for the described schema property. When defined within items, it will affect the name of the individual XML elements within the list. When defined alongside type being array (outside the items), it will affect the wrapping element and only if wrapped is true. If wrapped is false, it will be ignored.
	Namespace string `json:"namespace,omitempty"` // The URI of the namespace definition. This MUST be in the form of an absolute URI.
	Prefix    string `json:"prefix,omitempty"`    //	The prefix to be used for the name.
	Attribute bool   `json:"attribute,omitempty"` //	Declares whether the property definition translates to an attribute instead of an element. Default value is false.
	Wrapped   bool   `json:"wrapped,omitempty"`   //	MAY be used only for an array definition. Signifies whether the array is wrapped (for example, <books><book/><book/></books>) or unwrapped (<book/><book/>). Default value is false. The definition takes effect only when defined alongside type being array (outside the items).
}

type Schema struct {
	Version string        `json:"$schema,omitempty"` // section 6.1
	Ref     string        `json:"$ref,omitempty"`
	Enum    []interface{} `json:"enum,omitempty"` //枚举值：必须是一个数组。这个数组应该 .至少有一个元素。数组中的元素应该是唯一的。如果实例的值等于此关键字数组值中的元素之一，则该实例成功验证此关键字。
	Type    string        `json:"type,omitempty"` //必须是null boolean,object,array,number,string,integer
	// integer number
	MultipleOf       float64  `json:"multipleOf,omitempty"`       //几的倍数
	Maximum          *float64 `json:"maximum,omitempty"`          //数字小于等于
	ExclusiveMaximum bool     `json:"exclusiveMaximum,omitempty"` //数字小于
	Minimum          *float64 `json:"minimum,omitempty"`          //数字大于等于
	ExclusiveMinimum bool     `json:"exclusiveMinimum,omitempty"`
	//string
	MaxLength int    `json:"maxLength,omitempty"` //字符串最长值
	MinLength int    `json:"minLength,omitempty"` // 字符串最短值
	Pattern   string `json:"pattern,omitempty"`   // 正则表达式,注意
	// slice(array)
	MaxItems    uint64  `json:"maxItems,omitempty"`    //数组长度小于等于
	MinItems    uint64  `json:"minItems,omitempty"`    //数组长度大于等于
	UniqueItems bool    `json:"uniqueItems,omitempty"` //是否slice值不能重复
	Items       *Schema `json:"items,omitempty"`       // 验证切片是否属于这个类型
	//object
	MaxProperties        uint64                 `json:"maxProperties,omitempty"`        //object实例属性数量最大值
	MinProperties        uint64                 `json:"minProperties,omitempty"`        // 同上
	Required             []string               `json:"required,omitempty"`             //必须具备的属性
	Properties           *orderedmap.OrderedMap `json:"properties,omitempty"`           //直接定义属性
	AdditionalProperties json.RawMessage        `json:"additionalProperties,omitempty"` //true ： json串可以出现除schema定义之外属性 　　false ：json串不可以出现除schema定义之外属性
	PatternProperties    map[string]*Schema     `json:"patternProperties,omitempty"`    //对字段名称进行正则表达式验证
	//go 结构体可不支持动态结构，所以这几个字段没意义
	AllOf []*Schema `json:"allOf,omitempty"` //必须对所有子模式有效
	OneOf []*Schema `json:"oneOf,omitempty"` //必须仅对其中一个子模式有效
	AnyOf []*Schema `json:"anyOf,omitempty"` //组合类型，只要复合这里面的其中一个验证就行
	Not   *Schema   `json:"not,omitempty"`
	//元数据，所有对象都有的默认值
	Title       string      `json:"title,omitempty"`
	Description string      `json:"description,omitempty"`
	Default     interface{} `json:"default,omitempty"`

	Format   string        `json:"format,omitempty"` // 一些固定类型的匹配，比如email，ip，uuid，datetime等
	Example  *Example      `json:"example,omitempty"`
	Examples []interface{} `json:"examples,omitempty,omitempty"` // 例子
	//go不支持多态（这个字段没啥意义）
	//Discriminator Discriminator         `json:"discriminator,omitempty"`      //    Adds support for polymorphism.The discriminator is an object name that is used to differentiate between other schemas which may satisfy the payload description.See Composition and Inheritance for more details.

	//openapi需要的字段
	Xml          *XML                   `json:"xml,omitempty"`          //    This MAY be used only on properties schemas.It has no effect on root schemas.Adds additional metadata to describe the XML representation of this property.
	ExternalDocs *ExternalDocumentation `json:"externalDocs,omitempty"` //   Additional external documentation for this schema.

	// 以json编码的非json数据（比如base64：{
	//    "type": "string",
	//    "media": "base64",
	//    "contentMediaType": "image/png"
	//}）
	Media          *Schema `json:"media,omitempty"`          //
	BinaryEncoding string  `json:"binaryEncoding,omitempty"` //

	Extras map[string]interface{} `json:"-"`
	// 额外增加的字段,主要解决decode和encode共用一个结构体的时候字段问题
	Nullable   bool   `json:"nullable,omitempty"`
	ReadOnly   bool   `json:"readOnly,omitempty"`
	WriteOnly  bool   `json:"writeOnly,omitempty"`  //重用的时候只写，
	Deprecated bool   `json:"deprecated,omitempty"` //是否为废弃状态
	FieldName  string `json:"-"`                    //字段名称
}

func (t *Schema) MarshalJSON() ([]byte, error) {
	type Type_ Schema
	b, err := json.Marshal((*Type_)(t))
	if err != nil {
		return nil, err
	}
	if t.Extras == nil || len(t.Extras) == 0 {
		return b, nil
	}
	m, err := json.Marshal(t.Extras)
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

// StructKeywordsFromTags 解析结构体字段的tag
func (t *Schema) StructKeywordsFromTags(f reflect.StructField, parentType *Schema, propertyName string) {
	t.Description = f.Tag.Get(Description)
	tags := strings.Split(f.Tag.Get(TagName), ",")
	t.genericKeywords(tags, parentType, propertyName)
	//default title
	if t.Title == "" {
		t.Title = f.Name
	}
	t.FieldName = f.Name
	switch t.Type {
	case "string":
		t.stringKeywords(tags)
	case "number":
		t.numbericKeywords(tags)
	case "integer":
		t.numbericKeywords(tags)
	case "array":
		t.arrayKeywords(tags)
	}
	extras := strings.Split(f.Tag.Get(Extras), ",")
	t.extraKeywords(extras)
}

// read struct tags for generic keyworks
func (t *Schema) genericKeywords(tags []string, parentType *Schema, propertyName string) {
	for _, tag := range tags {
		nameValue := strings.Split(tag, "=")
		if len(nameValue) == 2 {
			name, val := nameValue[0], nameValue[1]
			switch name {
			case "title":
				t.Title = val
			case "desc":
				t.Description = val
			case "type":
				t.Type = val
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
				if t.OneOf == nil {
					t.OneOf = make([]*Schema, 0, 1)
				}
				t.Type = ""
				types := strings.Split(nameValue[1], ";")
				for _, ty := range types {
					t.OneOf = append(t.OneOf, &Schema{
						Type: ty,
					})
				}
			case "enum":
				switch t.Type {
				case "string":
					t.Enum = append(t.Enum, val)
				case "integer":
					i, _ := strconv.Atoi(val)
					t.Enum = append(t.Enum, i)
				case "number":
					f, _ := strconv.ParseFloat(val, 64)
					t.Enum = append(t.Enum, f)
				}
			}
		} else if len(nameValue) == 1 {
			if strings.Contains(tag, "|") {
				vales := strings.Split(tag, "|")
				for _, val := range vales {
					switch t.Items.Type {
					case "string":
						t.Items.Enum = append(t.Items.Enum, val)
					case "integer":
						i, _ := strconv.Atoi(val)
						t.Items.Enum = append(t.Items.Enum, i)
					case "number":
						f, _ := strconv.ParseFloat(val, 64)
						t.Items.Enum = append(t.Items.Enum, f)
					}
				}
			}

		}

	}
}

// read struct tags for string type keyworks
func (t *Schema) stringKeywords(tags []string) {
	for _, tag := range tags {
		nameValue := strings.Split(tag, "=")
		if len(nameValue) == 2 {
			name, val := nameValue[0], nameValue[1]
			switch name {
			case "len":
				i, _ := strconv.Atoi(val)
				t.MinLength = i
				t.MaxLength = i
			case "minLen":
				i, _ := strconv.Atoi(val)
				t.MinLength = i
			case "maxLen":
				i, _ := strconv.Atoi(val)
				t.MaxLength = i
			case "pattern":
				t.Pattern = val
			case "format":
				switch val {
				case "date-time", "email", "hostname", "ipv4", "ipv6", "uri":
					t.Format = val
					break
				}
			case "default":
				t.Default = val
			case "example":
				t.Examples = append(t.Examples, val)
			}
		}
	}
}

// read struct tags for numberic type keyworks
func (t *Schema) numbericKeywords(tags []string) {
	for _, tag := range tags {
		nameValue := strings.Split(tag, "=")
		if len(nameValue) == 2 {
			name, val := nameValue[0], nameValue[1]
			switch name {
			case "multi":
				i, _ := strconv.ParseFloat(val, 32)
				t.MultipleOf = i
			case "gte":
				i, _ := strconv.ParseFloat(val, 32)
				t.Minimum = &i
			case "gt":
				i, _ := strconv.ParseFloat(val, 32)
				t.Minimum = &i
				t.ExclusiveMinimum = true
			case "lte":
				i, _ := strconv.ParseFloat(val, 32)
				t.Maximum = &i
			case "lt":
				i, _ := strconv.ParseFloat(val, 32)
				t.Maximum = &i
				t.ExclusiveMaximum = true
			//case "eq": //todo 没啥用
			//case "ne":
			case "default":
				i, _ := strconv.Atoi(val)
				t.Default = i
			case "example":
				if i, err := strconv.Atoi(val); err == nil {
					t.Examples = append(t.Examples, i)
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
func (t *Schema) arrayKeywords(tags []string) {
	var defaultValues []interface{}
	for _, tag := range tags {
		nameValue := strings.Split(tag, "=")
		if len(nameValue) == 2 {
			name, val := nameValue[0], nameValue[1]
			switch name {
			case "minLen":
				i, _ := strconv.ParseUint(val, 10, 32)
				t.MinItems = i
			case "maxLen":
				i, _ := strconv.ParseUint(val, 10, 32)
				t.MaxItems = i
			case "len":
				i, _ := strconv.ParseUint(val, 10, 32)
				t.MinItems = i
				t.MaxItems = i
			case "unique":
				t.UniqueItems = true
			case "default":
				defaultValues = append(defaultValues, val)
			case "enum":
				switch t.Items.Type {
				case "string":
					t.Items.Enum = append(t.Items.Enum, val)
				case "integer":
					i, _ := strconv.Atoi(val)
					t.Items.Enum = append(t.Items.Enum, i)
				case "number":
					f, _ := strconv.ParseFloat(val, 64)
					t.Items.Enum = append(t.Items.Enum, f)
				}
			}
		} else if len(nameValue) == 1 {
			if strings.Contains(tag, "|") {
				vales := strings.Split(tag, "|")
				for _, val := range vales {
					switch t.Items.Type {
					case "string":
						t.Items.Enum = append(t.Items.Enum, val)
					case "integer":
						i, _ := strconv.Atoi(val)
						t.Items.Enum = append(t.Items.Enum, i)
					case "number":
						f, _ := strconv.ParseFloat(val, 64)
						t.Items.Enum = append(t.Items.Enum, f)
					}
				}
			}
		}
	}
	if len(defaultValues) > 0 {
		t.Default = defaultValues
	}
}

func (t *Schema) extraKeywords(tags []string) {
	for _, tag := range tags {
		nameValue := strings.Split(tag, "=")
		if len(nameValue) == 2 {
			t.setExtra(nameValue[0], nameValue[1])
		}
	}
}

func (t *Schema) setExtra(key, val string) {
	if t.Extras == nil {
		t.Extras = map[string]interface{}{}
	}
	if existingVal, ok := t.Extras[key]; ok {
		switch existingVal := existingVal.(type) {
		case string:
			t.Extras[key] = []string{existingVal, val}
		case []string:
			t.Extras[key] = append(existingVal, val)
		case int:
			t.Extras[key], _ = strconv.Atoi(val)
		}
	} else {
		t.Extras[key] = val
	}
}
