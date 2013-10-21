package bingo

import (
	"reflect"
	"testing"
)

func benchMarshal(b *testing.B, v interface{}) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Marshal(v)
	}
}

func benchShiftBytes(b *testing.B, data []byte, n uint, shiftFunc func([]byte, uint) []byte) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		shiftFunc(data, n)
	}
}

func BenchmarkMarshalUint8(b *testing.B) {
	var v uint8
	benchMarshal(b, v)
}

func BenchmarkMarshalUint16(b *testing.B) {
	var v uint16
	benchMarshal(b, v)
}

func BenchmarkMarshalUint32(b *testing.B) {
	var v uint32
	benchMarshal(b, v)
}

func BenchmarkMarshalUint64(b *testing.B) {
	var v uint64
	benchMarshal(b, v)
}

func BenchmarkMarshalInt8(b *testing.B) {
	var v int8
	benchMarshal(b, v)
}

func BenchmarkMarshalInt16(b *testing.B) {
	var v int16
	benchMarshal(b, v)
}

func BenchmarkMarshalInt32(b *testing.B) {
	var v int32
	benchMarshal(b, v)
}

func BenchmarkMarshalInt64(b *testing.B) {
	var v int64
	benchMarshal(b, v)
}

func BenchmarkMarshalFloat32(b *testing.B) {
	var v float32
	benchMarshal(b, v)
}

func BenchmarkMarshalFloat64(b *testing.B) {
	var v float64
	benchMarshal(b, v)
}

func BenchmarkMarshalComplex54(b *testing.B) {
	var v complex64
	benchMarshal(b, v)
}

func BenchmarkMarshalComplex128(b *testing.B) {
	var v complex128
	benchMarshal(b, v)
}

func BenchmarkMarshalString0(b *testing.B) {
	var v string
	benchMarshal(b, v)
}

func BenchmarkMarshalString32(b *testing.B) {
	v := createStringWithLength(32)
	benchMarshal(b, v)
}

func BenchmarkMarshalString1024(b *testing.B) {
	v := createStringWithLength(1024)
	benchMarshal(b, v)
}

func createStringWithLength(n int) string {
	var s string
	for i := 0; i < n; i++ {
		s += "a"
	}
	return s
}

func BenchmarkLeftShift64BitsAsBytes(b *testing.B) {
	data := make([]byte, 8)
	benchShiftBytes(b, data, 1, leftShiftBytes)
}

func BenchmarkLeftShift64BitsAsUint64(b *testing.B) {
	var data uint64
	for i := 0; i < b.N; i++ {
		_ = data << 1
	}
}

func BenchmarkRightShift64BitsAsBytes(b *testing.B) {
	data := make([]byte, 8)
	benchShiftBytes(b, data, 1, rightShiftBytes)
}

func BenchmarkRightShift64BitsAsUint64(b *testing.B) {
	var data uint64
	for i := 0; i < b.N; i++ {
		_ = data >> 1
	}
}

func BenchmarkLeftShift32(b *testing.B) {
	n := uint(32)
	data := make([]byte, n)
	benchShiftBytes(b, data, 1, leftShiftBytes)
}

func BenchmarkLeftShift1024(b *testing.B) {
	n := uint(1024)
	data := make([]byte, n)
	benchShiftBytes(b, data, 1, leftShiftBytes)
}

func BenchmarkLeftShift32768(b *testing.B) {
	n := uint(32768)
	data := make([]byte, n)
	benchShiftBytes(b, data, 1, leftShiftBytes)
}

type reflectStruct struct {
	A uint32
	B uint8
}

func BenchmarkNormalSet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := reflectStruct{}
		v.A = 100
		v.B = 200
	}
}

func BenchmarkReflectSet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := reflectStruct{}
		rv := reflect.ValueOf(&v)
		rv.Elem().Field(0).SetUint(100)
		rv.Elem().Field(1).SetUint(200)
	}
}
