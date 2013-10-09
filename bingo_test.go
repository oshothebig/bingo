package bingo

import (
	"bytes"
	"reflect"
	"testing"
)

var (
	uint8Sample  uint8  = 0x12
	uint16Sample uint16 = 0x1234
	uint32Sample uint32 = 0x12345678
	uint64Sample uint64 = 0x123456789abcdef0
	int8Sample   int8   = 0x12
	int16Sample  int16  = 0x1234
	int32Sample  int32  = 0x12345678
	int64Sample  int64  = 0x123456789abcdef0
)

var numberSamples = []struct {
	Number interface{}
	Bytes  []byte
}{
	{uint8Sample, []byte{0x12}},
	{uint16Sample, []byte{0x12, 0x34}},
	{uint32Sample, []byte{0x12, 0x34, 0x56, 0x78}},
	{uint64Sample, []byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0}},
	{int8Sample, []byte{0x12}},
	{int16Sample, []byte{0x12, 0x34}},
	{int32Sample, []byte{0x12, 0x34, 0x56, 0x78}},
	{int64Sample, []byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0}},
}

var numberPtrSamples = []struct {
	Number interface{}
	Bytes  []byte
}{
	{&uint8Sample, []byte{0x12}},
	{&uint16Sample, []byte{0x12, 0x34}},
	{&uint32Sample, []byte{0x12, 0x34, 0x56, 0x78}},
	{&uint64Sample, []byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0}},
	{&int8Sample, []byte{0x12}},
	{&int16Sample, []byte{0x12, 0x34}},
	{&int32Sample, []byte{0x12, 0x34, 0x56, 0x78}},
	{&int64Sample, []byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0}},
}

func checkMarshal(t *testing.T, v interface{}, expected []byte) {
	name := reflect.TypeOf(v).Name()

	actual, err := Marshal(v)
	if err != nil {
		t.Errorf("Error (%v) occurred when marshaling: %s", err, name)
		return
	}
	if !bytes.Equal(actual, expected) {
		t.Errorf("Bytes of %s not matched", name)
		t.Logf("Actual  : %x\n", actual)
		t.Logf("Expected: %x\n", expected)
	}
}

func checkUnmarshal(t *testing.T, b []byte, expected interface{}) {
	actual := reflect.New(reflect.TypeOf(expected).Elem()).Interface()
	name := reflect.TypeOf(expected).Name()

	err := Unmarshal(b, actual)
	if err != nil {
		t.Errorf("Error (%v) occurred when unmarshaling: %s", err, name)
		return
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Fields of %s not matched", name)
		t.Logf("Actual  :%v\n", actual)
		t.Logf("Expected:%v\n", expected)
	}
}

func checkUnsupportedType(t *testing.T, v interface{}) {
	name := reflect.TypeOf(v).Name()

	_, err := Marshal(v)
	if err == nil {
		t.Errorf("Should fail marshal, but succeeded: %s", name)
	}
}

func TestMarshalNumber(t *testing.T) {
	for _, v := range numberSamples {
		checkMarshal(t, v.Number, v.Bytes)
	}
}

func TestMarshalNumberPrt(t *testing.T) {
	for _, v := range numberPtrSamples {
		checkMarshal(t, v.Number, v.Bytes)
	}
}

func TestUnmarshalNumberPtr(t *testing.T) {
	for _, v := range numberPtrSamples {
		checkUnmarshal(t, v.Bytes, v.Number)
	}
}

func TestMarshalString(t *testing.T) {
	s := "abcde"

	checkMarshal(t, s, []byte(s))
}

func TestMarshalStringPtr(t *testing.T) {
	s := "abcde"

	checkMarshal(t, s, []byte(s))
}

type flatStructOnlyWithNumber struct {
	Age    uint8
	Height uint16
}

func TestMarshalPrimitiveOnlyStruct(t *testing.T) {
	v := flatStructOnlyWithNumber{31, 171}
	b := []byte{0x1f, 0x00, 0xab}
	checkMarshal(t, v, b)
	checkMarshal(t, &v, b)
}

type flatStructWithNumberPtr struct {
	Age    *uint8
	Height *uint16
}

func TestMarshalFlatStructWithNumberPrt(t *testing.T) {
	var age uint8 = 31
	var height uint16 = 171
	v := flatStructWithNumberPtr{&age, &height}
	b := []byte{0x1f, 0x00, 0xab}
	checkMarshal(t, v, b)
	checkMarshal(t, &v, b)

}

type flatStructWithString struct {
	Age  uint8
	Name string
}

func TestMarshalFlatStructWithString(t *testing.T) {
	name := "Sho Shimizu"
	v := flatStructWithString{31, name}
	b := []byte{0x1f}
	b = append(b, []byte(name)...)
	checkMarshal(t, v, b)
	checkMarshal(t, &v, b)
}

type flatStructWithMap struct {
	Age       uint8
	Attribute map[string]string
}

func TestMarshalFlatStructWithMap(t *testing.T) {
	v := flatStructWithMap{31, map[string]string{"Name": "Bob"}}
	checkUnsupportedType(t, v)
}

type flatStructWithFunc struct {
	Age  uint8
	Func func()
}

func TestMarshalFlatStructWithFunc(t *testing.T) {
	v := flatStructWithFunc{31, func() {}}
	checkUnsupportedType(t, v)
}

type flatStructWithInt struct {
	Age int
}

func TestMarshalFlatStructWithInt(t *testing.T) {
	v := flatStructWithInt{31}
	checkUnsupportedType(t, v)
}

type flatStructWithUint struct {
	Age uint
}

func TestMarshalFlatStructWithUint(t *testing.T) {
	v := flatStructWithUint{31}
	checkUnsupportedType(t, v)
}

type Human interface {
	Name() string
}

type Person string

func (p Person) Name() string {
	return string(p)
}

type flatStructWithInterface struct {
	Age   uint8
	Human Human
}

func TestMarshalflatStructWithInterface(t *testing.T) {
	name := "Sho Shimizu"
	var p Person = Person(name)
	v := flatStructWithInterface{31, p}
	b := append([]byte{0x1f}, []byte(name)...)
	checkMarshal(t, v, b)
}

type flatStructWithBlankField struct {
	Age uint8
	_   [8]uint8
}

func TestMarshalFlatStructWithBlankField(t *testing.T) {
	v := flatStructWithBlankField{Age: 31}
	b := []byte{0x1f, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	checkMarshal(t, v, b)
}

func TestMarshalUnsupportedType(t *testing.T) {
	m := map[string]uint32{
		"One": 1,
		"Two": 2,
	}
	checkUnsupportedType(t, m)

	f := func() {}
	checkUnsupportedType(t, f)
}
