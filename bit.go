package bingo

import "reflect"

type bitField struct {
	value reflect.Value
	bits  uint
}

type bitFields []bitField

func (b bitFields) add(v reflect.Value, bit uint) bitFields {
	switch v.Kind() {
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return append(b, bitField{v, bit})
	}
	return b
}

func (b bitFields) bits() uint {
	var bits uint = 0
	for _, v := range b {
		bits += v.bits
	}
	return bits
}

func byteAlignInBit(bits uint) uint {
	return (bits + 7) / 8 * 8
}

func leftShiftBytes(b []byte, n uint) []byte {
	shiftBytes := n / 8
	shiftBits := n % 8
	length := uint(len(b))
	shifted := make([]byte, length)

	for i := shiftBytes; i < length; i++ {
		indexFrom := i - shiftBytes
		var carry byte
		if indexFrom >= 1 {
			carry = b[indexFrom-1] >> (8 - shiftBits)
		}
		shifted[i] = (b[indexFrom] << shiftBits) | carry
	}

	return shifted
}

func rightShiftBytes(b []byte, n uint) []byte {
	shiftBytes := n / 8
	shiftBits := n % 8
	length := uint(len(b))
	shifted := make([]byte, length)

	for i := uint(0); i < length-shiftBytes; i++ {
		indexFrom := i + shiftBytes
		var carry byte
		if indexFrom+1 < length {
			carry = b[indexFrom+1] << (8 - shiftBits)
		}
		shifted[i] = carry | (b[indexFrom] >> shiftBits)
	}

	return shifted
}
