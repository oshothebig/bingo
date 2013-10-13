package bingo

import (
	"encoding/binary"
	"errors"
	"reflect"
)

type bitField struct {
	value reflect.Value
	bits  uint
}

type bitFields struct {
	fields []bitField
	bits   uint
}

func newBitFields() *bitFields {
	return &bitFields{make([]bitField, 0, 16), 0}
}

func (b *bitFields) add(v reflect.Value, bit uint) *bitFields {
	switch v.Kind() {
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		b.fields = append(b.fields, bitField{v, bit})
		b.bits += bit
	}
	return b
}

func (b *bitFields) bytes() ([]byte, error) {
	if b.bits > 64 {
		return []byte{}, errors.New("bingo: over 64 bit field not supported")
	}

	var result uint64
	var shiftRemaining uint = 64
	for _, field := range b.fields {
		var num uint64
		switch field.value.Kind() {
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			num = uint64(field.value.Int())
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			num = field.value.Uint()
		}

		// get lower {bits} bits
		num = num & ((1 << field.bits) - 1)
		shiftRemaining -= field.bits
		result |= num << shiftRemaining
	}

	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, result)
	return data[:byteAlign(b.bits)], nil
}

func byteAlign(bits uint) uint {
	return (bits + 7) / 8
}

func leftShiftBytes(b []byte, n uint) []byte {
	shiftBytes := n / 8
	shiftBits := n % 8
	length := uint(len(b))
	shifted := make([]byte, length)

	shifted[shiftBytes] = b[0] << shiftBits
	for i := shiftBytes + 1; i < length; i++ {
		srcIndex := i - shiftBytes
		carry := b[srcIndex-1] >> (8 - shiftBits)
		shifted[i] = (b[srcIndex] << shiftBits) | carry
	}

	return shifted
}

func rightShiftBytes(b []byte, n uint) []byte {
	shiftBytes := n / 8
	shiftBits := n % 8
	length := uint(len(b))
	shifted := make([]byte, length)

	last := length - 1
	shifted[last-shiftBytes] = b[last] >> shiftBits
	for i := uint(0); i < length-shiftBytes-1; i++ {
		srcIndex := i + shiftBytes
		carry := b[srcIndex+1] << (8 - shiftBits)
		shifted[i] = carry | (b[srcIndex] >> shiftBits)
	}

	return shifted
}
