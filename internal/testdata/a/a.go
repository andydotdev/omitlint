package a

import (
	"testdata/a/b"
)

type (
	aliasBasic           = int
	underlyingBasic      int
	aliasUnderlyingBasic = underlyingBasic
)

type (
	aliasStruct           = struct{}
	structType            struct{}
	aliasUnderlyingStruct = structType
)

type (
	underlyingExternalStructType       b.StructType
	alisasUnderlyingExternalStructType = b.StructType
)

// Good should have no suggested changes.
type Good struct {
	// basic types have omittable zero values
	Bool                 bool                 `json:"bool,omitempty"`
	Int                  int                  `json:"int,omitempty"`
	Int8                 int8                 `json:"int8,omitempty"`
	Int16                int16                `json:"int16,omitempty"`
	Int32                int32                `json:"int32,omitempty"`
	Int64                int64                `json:"int64,omitempty"`
	Uint                 uint                 `json:"uint,omitempty"`
	Uint8                uint8                `json:"uint8,omitempty"`
	Uint16               uint16               `json:"uint16,omitempty"`
	Uint32               uint32               `json:"uint32,omitempty"`
	Uint64               uint64               `json:"uint64,omitempty"`
	Float32              float32              `json:"float32,omitempty"`
	Float64              float64              `json:"float64,omitempty"`
	String               string               `json:"string,omitempty"` // Strings can have len 0
	Array                [0]structType        `json:"array,omitempty"`  // Arrays of len 0 can be omitted.
	Slice                []string             `json:"slice,omitempty"`  // Slices can have len 0 or be explicitly nil.
	Map                  map[any]structType   `json:"map,omitempty"`    // Maps can have len 0 or be explicitly nil.
	UnderlyingBasic      underlyingBasic      `json:"underlyingbasic,omitempty"`
	AliasBasic           aliasBasic           `json:"aliasbasic,omitempty"`
	AliasUnderlyingBasic aliasUnderlyingBasic `json:"aliasunderlyingbasic,omitempty"`
	unexported           structType           `json"unexported,omitempty"`   // Unexported
	Excluded             structType           `json:"-"`                     // Excluded from encoding
	NotOmitempty         structType           `json:"notOmitempty"`          // No omitempty option"
	NamedOmitEmpty       structType           `json:"omitempty"`             // named "omitempty" (not an encoding option)
	Interface            any                  `json:"interface,omitempty"`   // Interfaces can be nil; however, see TODO in omitempty.go
	SliceCap0            []int                `json:"slice_cap_0,omitempty"` // (for assert testing)
}

// Unencodable has fields with types that cannot be encoded, but would already
// be reported as a json.UnsupportedTypeError, and is outside the scope of this
// linter, so they are left alone for now.
type Unencodable struct {
	Complex64  complex64  `json:"complex64,omitempty"`
	Complex128 complex128 `json:"complex128,omitempty"`
	ChanInt    chan int   `json:"chanint,omitempty"`
	Func       func()     `json:"func,omitempty"`
}

// All of the following fields would otherwise always be encoded and should be
// changed to pointer types by the analysis suggestion.
type Bad struct {
	StructType                        structType                         `json:"struct_type,omitempty"`                           // want `field "StructType" is marked "omitempty", but cannot be omitted`
	ExternalStructType                b.StructType                       `json:"external_struct_type,omitempty"`                  // want `field "ExternalStructType" is marked "omitempty", but cannot be omitted`
	UnderlyingExternalStructType      underlyingExternalStructType       `json:"underlying_external_struct_type,omitempty"`       // want `field "UnderlyingExternalStructType" is marked "omitempty", but cannot be omitted`
	AliasUnderlyingExternalStructType alisasUnderlyingExternalStructType `json:"alias_underlying_external_struct_type,omitempty"` // want `field "AliasUnderlyingExternalStructType" is marked "omitempty", but cannot be omitted`
}
