# go_query_dsl
--
    import "github.com/rekki/go-query/util/go_query_dsl"


## Usage

```go
var (
	ErrInvalidLengthDsl = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowDsl   = fmt.Errorf("proto: integer overflow")
)
```

```go
var Query_Type_name = map[int32]string{
	0: "TERM",
	1: "AND",
	2: "OR",
	3: "DISMAX",
}
```

```go
var Query_Type_value = map[string]int32{
	"TERM":   0,
	"AND":    1,
	"OR":     2,
	"DISMAX": 3,
}
```

#### type Query

```go
type Query struct {
	Queries    []*Query   `protobuf:"bytes,1,rep,name=queries,proto3" json:"queries,omitempty"`
	Type       Query_Type `protobuf:"varint,2,opt,name=type,proto3,enum=go.query.dsl.Query_Type" json:"type,omitempty"`
	Field      string     `protobuf:"bytes,3,opt,name=field,proto3" json:"field,omitempty"`
	Value      string     `protobuf:"bytes,4,opt,name=value,proto3" json:"value,omitempty"`
	Not        *Query     `protobuf:"bytes,5,opt,name=not,proto3" json:"not,omitempty"`
	Tiebreaker float32    `protobuf:"fixed32,6,opt,name=tiebreaker,proto3" json:"tiebreaker,omitempty"`
	Boost      float32    `protobuf:"fixed32,7,opt,name=boost,proto3" json:"boost,omitempty"`
}
```


#### func (*Query) Descriptor

```go
func (*Query) Descriptor() ([]byte, []int)
```

#### func (*Query) GetBoost

```go
func (m *Query) GetBoost() float32
```

#### func (*Query) GetField

```go
func (m *Query) GetField() string
```

#### func (*Query) GetNot

```go
func (m *Query) GetNot() *Query
```

#### func (*Query) GetQueries

```go
func (m *Query) GetQueries() []*Query
```

#### func (*Query) GetTiebreaker

```go
func (m *Query) GetTiebreaker() float32
```

#### func (*Query) GetType

```go
func (m *Query) GetType() Query_Type
```

#### func (*Query) GetValue

```go
func (m *Query) GetValue() string
```

#### func (*Query) Marshal

```go
func (m *Query) Marshal() (dAtA []byte, err error)
```

#### func (*Query) MarshalTo

```go
func (m *Query) MarshalTo(dAtA []byte) (int, error)
```

#### func (*Query) MarshalToSizedBuffer

```go
func (m *Query) MarshalToSizedBuffer(dAtA []byte) (int, error)
```

#### func (*Query) ProtoMessage

```go
func (*Query) ProtoMessage()
```

#### func (*Query) Reset

```go
func (m *Query) Reset()
```

#### func (*Query) Size

```go
func (m *Query) Size() (n int)
```

#### func (*Query) String

```go
func (m *Query) String() string
```

#### func (*Query) Unmarshal

```go
func (m *Query) Unmarshal(dAtA []byte) error
```

#### func (*Query) XXX_DiscardUnknown

```go
func (m *Query) XXX_DiscardUnknown()
```

#### func (*Query) XXX_Marshal

```go
func (m *Query) XXX_Marshal(b []byte, deterministic bool) ([]byte, error)
```

#### func (*Query) XXX_Merge

```go
func (m *Query) XXX_Merge(src proto.Message)
```

#### func (*Query) XXX_Size

```go
func (m *Query) XXX_Size() int
```

#### func (*Query) XXX_Unmarshal

```go
func (m *Query) XXX_Unmarshal(b []byte) error
```

#### type Query_Type

```go
type Query_Type int32
```


```go
const (
	Query_TERM   Query_Type = 0
	Query_AND    Query_Type = 1
	Query_OR     Query_Type = 2
	Query_DISMAX Query_Type = 3
)
```

#### func (Query_Type) EnumDescriptor

```go
func (Query_Type) EnumDescriptor() ([]byte, []int)
```

#### func (Query_Type) String

```go
func (x Query_Type) String() string
```
