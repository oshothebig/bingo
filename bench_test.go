package bingo

import "testing"

func benchMarshal(b *testing.B, v interface{}) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Marshal(v)
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
