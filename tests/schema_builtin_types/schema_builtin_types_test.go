package path_added

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sabahtalateh/oapi/tests"
)

//go:generate oapi schema SchemaA
type A struct {
	Bool       bool       `json:"bool"`
	Uint8      uint8      `json:"uint8"`
	Uint16     uint16     `json:"uint16"`
	Uint32     uint32     `json:"uint32"`
	Uint64     uint64     `json:"uint64"`
	Int8       int8       `json:"int8"`
	Int16      int16      `json:"int16"`
	Int32      int32      `json:"int32"`
	Int64      int64      `json:"int64"`
	Float32    float32    `json:"float32"`
	Float64    float64    `json:"float64"`
	Complex64  complex64  `json:"complex64"`
	Complex128 complex128 `json:"complex128"`
	String     string     `json:"string"`
	Int        int        `json:"int"`
	Uint       uint       `json:"uint"`
	Uintptr    uintptr    `json:"uintptr"`
	Byte       byte       `json:"byte"`
	Rune       rune       `json:"rune"`
	Any        any        `json:"any"`
	Bytes      []byte     `json:"bytes"`
}

var expected = tests.TrimOApi(`
openapi: 3.0.3
info:
  title: << YOUR TEXT >>
  version: 0.0.1
paths:
components:
  schemas:
    SchemaA:
      type: object
      required:
        - bool
        - uint8
        - uint16
        - uint32
        - uint64
        - int8
        - int16
        - int32
        - int64
        - float32
        - float64
        - complex64
        - complex128
        - string
        - int
        - uint
        - uintptr
        - byte
        - rune
        - any
        - bytes
      properties:
        bool:
          type: boolean
        uint8:
          type: integer
          format: int32
        uint16:
          type: integer
          format: int32
        uint32:
          type: integer
          format: int32
        uint64:
          type: integer
          format: int64
        int8:
          type: integer
          format: int32
        int16:
          type: integer
          format: int32
        int32:
          type: integer
          format: int32
        int64:
          type: integer
          format: int64
        float32:
          type: number
          format: float
        float64:
          type: number
          format: double
        complex64:
          type: string
        complex128:
          type: string
        string:
          type: string
        int:
          type: integer
          format: int64
        uint:
          type: integer
          format: int64
        uintptr:
          type: integer
          format: int64
        byte:
          type: string
          format: byte
        rune:
          type: string
        any:
          type: string
        bytes:
          type: string
          format: byte
`)

func Test(t *testing.T) {
	err := tests.GoGenerate(t, tests.Dir())
	require.NoError(t, err)
	defer tests.CleanGenerated(t, tests.Dir())

	oapi, _ := tests.ReadOAPI(t)
	require.Equal(t, expected, oapi)
}
