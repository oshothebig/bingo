package bingo

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"reflect"
	"strconv"
	"strings"
)

func boolToUint8(v bool) uint8 {
	if v == true {
		return 0x01
	} else {
		return 0x00
	}
}

func toWritablePrimitiveData(orig interface{}) (conv interface{}, suceeded bool) {
	switch p := orig.(type) {
	case int8, int16, int32, int64, *int8, *int16, *int32, *int64,
		uint8, uint16, uint32, uint64, *uint8, *uint16, *uint32, *uint64,
		float32, float64, *float32, *float64,
		complex64, complex128, *complex64, *complex128:
		conv = orig
	case string:
		conv = []byte(p)
	case *string:
		conv = []byte(*p)
	case bool:
		conv = boolToUint8(p)
	case *bool:
		conv = boolToUint8(*p)
	default:
		return nil, false
	}
	return conv, true
}

func Marshal(data interface{}) (b []byte, err error) {
	buf := new(bytes.Buffer)
	// fast path for primitive type
	if data, ok := toWritablePrimitiveData(data); ok {
		err = binary.Write(buf, binary.BigEndian, data)
		return buf.Bytes(), err
	}

	// reflection based encoding
	v := reflect.Indirect(reflect.ValueOf(data))
	//	map and func type are not supported
	if v.Kind() == reflect.Map || v.Kind() == reflect.Func {
		err = errors.New("Unsupported type")
		return []byte{}, err
	}

	enc := &refEncoder{binary.BigEndian, new(bytes.Buffer)}
	err = enc.encode(v, nil)
	return enc.buf.Bytes(), err
}

type refEncoder struct {
	order binary.ByteOrder
	buf   *bytes.Buffer
}

func (e *refEncoder) encode(v reflect.Value, opt *tagOptions) error {
	v = reflect.Indirect(v)
	switch v.Kind() {
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		return binary.Write(e.buf, e.order, v.Interface())
	case reflect.String:
		s := v.String()
		return binary.Write(e.buf, e.order, []byte(s))
	case reflect.Bool:
		b := v.Bool()
		return binary.Write(e.buf, e.order, boolToUint8(b))
	case reflect.Array, reflect.Slice:
		l := v.Len()
		for i := 0; i < l; i++ {
			if err := e.encode(v.Index(i), nil); err != nil {
				return err
			}
		}
	case reflect.Struct:
		t := v.Type()
		l := v.NumField()
		for i := 0; i < l; i++ {
			var field reflect.Value
			if fv := v.Field(i); fv.CanSet() || t.Field(i).Name != "_" {
				field = fv
			} else {
				field = reflect.Zero(fv.Type())
			}

			if err := e.encode(field, nil); err != nil {
				return err
			}
		}
	case reflect.Interface:
		if err := e.encode(v.Elem(), nil); err != nil {
			return err
		}
	case reflect.Map, reflect.Func, reflect.Int, reflect.Uint, reflect.Chan:
		return errors.New("Unsupported type")
	}
	return nil
}

type tagOptions struct {
	bits uint
}

func parseTagOptions(tag string) *tagOptions {
	options := map[string]string{}
	tokens := strings.Split(tag, ",")
	for _, token := range tokens {
		token = strings.TrimSpace(token)
		keyval := strings.Split(token, ":")
		options[keyval[0]] = keyval[1]
	}

	return newTagOptions(options)
}

func newTagOptions(options map[string]string) *tagOptions {
	topts := new(tagOptions)
	for k, v := range options {
		switch k {
		case "bits":
			num, err := strconv.ParseUint(v, 10, 32)
			if err != nil {
				continue
			}
			topts.bits = uint(num)
		}
	}
	return topts
}

type Encoder struct {
	writer io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{writer: w}
}

func (enc *Encoder) Encode(v interface{}) error {
	b, err := Marshal(v)
	if err != nil {
		return err
	}
	_, err = enc.writer.Write(b)
	return err
}

func Unmarshal(data []byte, v interface{}) error {
	buf := bytes.NewBuffer(data)
	switch v.(type) {
	case *int8, *int16, *int32, *int64,
		*uint8, *uint16, *uint32, *uint64,
		*float32, *float64, *complex64, *complex128:
		return binary.Read(buf, binary.BigEndian, v)
	default:
		return errors.New("Unsupported type: not implemented")
	}
}

type Decoder struct {
	reader io.Reader
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{reader: r}
}

func (enc *Decoder) Decode(v interface{}) error {
	return nil
}
