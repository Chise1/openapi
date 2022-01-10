# openapi实验项目

支持tag:

```text
openapi //默认的tag标签名
openapi_desc //如果描述里面带','则可以用这个tag
openapi_extras //额外的字段，和openapi一样，防止出现','用的
```

示例

```go
package openapi

import (
	"encoding/json"
	"github.com/Chise1/openapi/models"
	"github.com/iancoleman/orderedmap"
	"github.com/stretchr/testify/require"
	"net"
	"net/url"
	"time"
)

type GrandfatherType struct {
	FamilyName string `json:"family_name" openapi:"required"`
}

type SomeBaseType struct {
	SomeBaseProperty     int `json:"some_base_property"`
	SomeBasePropertyYaml int `yaml:"some_base_property_yaml"`
	// The openapi required tag is nonsensical for private and ignored properties.
	// Their presence here tests that the fields *will not* be required in the output
	// schema, even if they are tagged required.
	somePrivateBaseProperty   string          `openapi:"required"`
	SomeIgnoredBaseProperty   string          `json:"-" openapi:"required"`
	SomeSchemaIgnoredProperty string          `openapi:"-,required"`
	Grandfather               GrandfatherType `json:"grand"`

	SomeUntaggedBaseProperty           bool `openapi:"required"`
	someUnexportedUntaggedBaseProperty bool
}

type MapType map[string]interface{}

type nonExported struct {
	PublicNonExported  int
	privateNonExported int
}

type ProtoEnum int32

func (ProtoEnum) EnumDescriptor() ([]byte, []int) { return []byte(nil), []int{0} }

const (
	Unset ProtoEnum = iota
	Great
)

type TestUser struct {
	SomeBaseType
	nonExported
	MapType

	ID      int                    `json:"id" openapi:"required"`
	Name    string                 `json:"name" openapi:"required,minLen=1,maxLen=20,pattern=.*,desc=this is a property,title=the name,example=joe,example=lucy,default=alex"`
	Friends []int                  `json:"friends,omitempty" openapi_desc:"list of IDs, omitted when empty"`
	Tags    map[string]interface{} `json:"tags,omitempty"`

	TestFlag       bool
	IgnoredCounter int `json:"-"`

	// Tests for RFC draft-wright-json-schema-validation-00, section 7.3
	BirthDate time.Time `json:"birth_date,omitempty"`
	Website   url.URL   `json:"website,omitempty"`
	IPAddress net.IP    `json:"network_address,omitempty"`

	// Tests for RFC draft-wright-json-schema-hyperschema-00, section 4
	Photo  []byte `json:"photo,omitempty" openapi:"required"`
	Photo2 Bytes  `json:"photo2,omitempty" openapi:"required"`

	// Tests for jsonpb enum support
	Feeling ProtoEnum `json:"feeling,omitempty"`
	Age     int       `json:"age" openapi:"gt=18,lt=120"`
	Email   string    `json:"email" openapi:"format=email"`

	// Test for "extras" support
	Baz string `openapi_extras:"foo=bar,hello=world,foo=bar1"`

	// Tests for simple enum Tags
	Color      string  `json:"color" openapi:"enum=red,enum=green,enum=blue"`
	Rank       int     `json:"rank,omitempty" openapi:"enum=1,enum=2,enum=3"`
	Multiplier float64 `json:"mult,omitempty" openapi:"enum=1.0,enum=1.5,enum=2.0"`

	// Tests for enum Tags on slices
	Roles      []string  `json:"roles" openapi:"enum=admin,enum=moderator,enum=user"`
	Priorities []int     `json:"priorities,omitempty" openapi:"enum=-1,enum=0,enum=1,enun=2"`
	Offsets    []float64 `json:"offsets,omitempty" openapi:"enum=1.570796,enum=3.141592,enum=6.283185"`

	// Test for raw JSON
	Raw json.RawMessage `json:"raw"`
}
```

int支持:lt,lte,gt,gte,multi,default,example,enum,required enum的简单写法是：直接用|进行分割。例如

```go
type E struct{
Offsets    []float64 `json:"offsets,omitempty" openapi:"1.570796|3.141592|6.283185"`
}
```

简写和enum两种方法都支持

# tips

来源：jsonschema