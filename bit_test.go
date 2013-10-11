package bingo

import (
	"bytes"
	"reflect"
	"testing"
)

var leftShiftData = []struct {
	orig     []byte
	shift    uint
	expected []byte
}{
	{[]byte{0x34, 0x12}, 1, []byte{0x68, 0x24}},
	{[]byte{0x34, 0x12}, 2, []byte{0xd0, 0x48}},
	{[]byte{0x34, 0x12}, 3, []byte{0xa0, 0x91}},
	{[]byte{0x34, 0x12}, 4, []byte{0x40, 0x23}},
	{[]byte{0x34, 0x12}, 5, []byte{0x80, 0x46}},
	{[]byte{0x34, 0x12}, 6, []byte{0x00, 0x8d}},
	{[]byte{0x34, 0x12}, 7, []byte{0x00, 0x1a}},
	{[]byte{0x34, 0x12}, 8, []byte{0x00, 0x34}},
	{[]byte{0x34, 0x12}, 9, []byte{0x00, 0x68}},
}

var rightShiftData = []struct {
	orig     []byte
	shift    uint
	expected []byte
}{
	{[]byte{0x34, 0x12}, 1, []byte{0x1a, 0x09}},
	{[]byte{0x34, 0x12}, 2, []byte{0x8d, 0x04}},
	{[]byte{0x34, 0x12}, 3, []byte{0x46, 0x02}},
	{[]byte{0x34, 0x12}, 4, []byte{0x23, 0x01}},
	{[]byte{0x34, 0x12}, 5, []byte{0x91, 0x00}},
	{[]byte{0x34, 0x12}, 6, []byte{0x48, 0x00}},
	{[]byte{0x34, 0x12}, 7, []byte{0x24, 0x00}},
	{[]byte{0x34, 0x12}, 8, []byte{0x12, 0x00}},
	{[]byte{0x34, 0x12}, 9, []byte{0x09, 0x00}},
}

func TestBitFieldsAdd(t *testing.T) {
	f := make(bitFields, 0)
	f = f.add(reflect.ValueOf(uint8(0)), 5)
	if len(f) != 1 {
		t.Error("add() failed, but expected to succeed")
	}

	f = f.add(reflect.ValueOf(complex64(1+1i)), 5)
	if len(f) != 1 {
		t.Error("add() succeeded, but expected to fail")
	}
}

func TestBitFieldBits(t *testing.T) {
	f := make(bitFields, 0)
	f = f.add(reflect.ValueOf(uint8(0x12)), 5)
	f = f.add(reflect.ValueOf(uint16(0x34)), 7)

	if f.bits() != 12 {
		t.Error("Total bits: %d, want %d", f.bits(), 12)
	}
}

func TestLeftShiftBytes(t *testing.T) {
	for _, v := range leftShiftData {
		actual := leftShiftBytes(v.orig, v.shift)
		if !bytes.Equal(actual, v.expected) {
			t.Errorf("%d bit left shift of %x: %x, want %x", v.shift, v.orig, actual, v.expected)
		}
	}
}

func TestRightShiftBytes(t *testing.T) {
	for _, v := range rightShiftData {
		actual := rightShiftBytes(v.orig, v.shift)
		if !bytes.Equal(actual, v.expected) {
			t.Errorf("%d bit left shift of %x: %x, want %x", v.shift, v.orig, actual, v.expected)
		}
	}
}
