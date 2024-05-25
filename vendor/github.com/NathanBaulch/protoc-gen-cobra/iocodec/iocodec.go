package iocodec

import (
	"encoding/xml"
	"errors"
	"io"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/NathanBaulch/protoc-gen-cobra/ptypes"
)

var NoOp = func(interface{}) error { return nil }

type (
	DecoderMaker func(io.Reader) Decoder
	Decoder      func(interface{}) error
	EncoderMaker func(io.Writer) Encoder
	Encoder      func(interface{}) error
)

func JSONDecoderMaker() DecoderMaker {
	return func(r io.Reader) Decoder {
		return func(v interface{}) error {
			if b, err := io.ReadAll(r); err != nil {
				return err
			} else {
				return protojson.Unmarshal(b, v.(proto.Message))
			}
		}
	}
}

func JSONEncoderMaker(pretty bool) EncoderMaker {
	m := &protojson.MarshalOptions{}
	if pretty {
		m.Indent = "  "
	}
	return func(w io.Writer) Encoder {
		return func(v interface{}) error {
			if b, err := m.Marshal(v.(proto.Message)); err != nil {
				return err
			} else {
				_, err := w.Write(b)
				return err
			}
		}
	}
}

func XMLDecoderMaker() DecoderMaker {
	return func(r io.Reader) Decoder { return xml.NewDecoder(r).Decode }
}

func XMLEncoderMaker(pretty bool) EncoderMaker {
	return func(w io.Writer) Encoder {
		return func(v interface{}) error {
			e := xml.NewEncoder(w)
			defer e.Flush()
			if pretty {
				e.Indent("", "  ")
			}
			if err := e.Encode(v); err != nil {
				return err
			}
			_, err := w.Write([]byte("\n"))
			return err
		}
	}
}

func DecodeKnownTypes(d Decoder) Decoder {
	return func(v interface{}) error {
		var i interface{}
		if err := d(&i); err != nil {
			return err
		}

		return decodeValue(i, v)
	}
}

func decodeValue(v interface{}, p interface{}) error {
	cfg := &mapstructure.DecoderConfig{
		Result:  p,
		TagName: "json",
		DecodeHook: func(from reflect.Type, to reflect.Type, v interface{}) (interface{}, error) {
			to = indirect(to)
			if to.Kind() == reflect.Struct {
				switch reflect.PtrTo(to) {
				case reflect.TypeOf((*timestamppb.Timestamp)(nil)):
					return ptypes.ToTimestamp(v)
				case reflect.TypeOf((*durationpb.Duration)(nil)):
					return ptypes.ToDuration(v)
				case reflect.TypeOf((*wrapperspb.DoubleValue)(nil)):
					return ptypes.ToDoubleWrapper(v)
				case reflect.TypeOf((*wrapperspb.FloatValue)(nil)):
					return ptypes.ToFloatWrapper(v)
				case reflect.TypeOf((*wrapperspb.Int64Value)(nil)):
					return ptypes.ToInt64Wrapper(v)
				case reflect.TypeOf((*wrapperspb.UInt64Value)(nil)):
					return ptypes.ToUInt64Wrapper(v)
				case reflect.TypeOf((*wrapperspb.Int32Value)(nil)):
					return ptypes.ToInt32Wrapper(v)
				case reflect.TypeOf((*wrapperspb.UInt32Value)(nil)):
					return ptypes.ToUInt32Wrapper(v)
				case reflect.TypeOf((*wrapperspb.BoolValue)(nil)):
					return ptypes.ToBoolWrapper(v)
				case reflect.TypeOf((*wrapperspb.StringValue)(nil)):
					return ptypes.ToStringWrapper(v)
				case reflect.TypeOf((*wrapperspb.BytesValue)(nil)):
					return ptypes.ToBytesWrapper(v)
				}
			}
			return v, nil
		},
	}
	if decoder, err := mapstructure.NewDecoder(cfg); err != nil {
		return err
	} else if err := decoder.Decode(v); err != nil {
		return err
	} else {
		return nil
	}
}

var (
	errUnchanged       = errors.New("unchanged")
	emptyInterfaceType = reflect.TypeOf((*interface{})(nil)).Elem()
)

func EncodeKnownTypes(e Encoder) Encoder {
	return func(v interface{}) error {
		if v, err := encodeValue(v); err != nil && err != errUnchanged {
			return err
		} else {
			return e(v)
		}
	}
}

func encodeValue(v interface{}) (interface{}, error) {
	switch v := v.(type) {
	case *timestamppb.Timestamp:
		return v.AsTime(), nil
	case *durationpb.Duration:
		return v.AsDuration(), nil
	case *wrapperspb.BoolValue:
		return v.Value, nil
	case *wrapperspb.BytesValue:
		return v.Value, nil
	case *wrapperspb.DoubleValue:
		return v.Value, nil
	case *wrapperspb.FloatValue:
		return v.Value, nil
	case *wrapperspb.Int32Value:
		return v.Value, nil
	case *wrapperspb.UInt32Value:
		return v.Value, nil
	case *wrapperspb.Int64Value:
		return v.Value, nil
	case *wrapperspb.UInt64Value:
		return v.Value, nil
	case *wrapperspb.StringValue:
		return v.Value, nil
	default:
		vv := reflect.Indirect(reflect.ValueOf(v))
		switch vv.Kind() {
		case reflect.Struct:
			return encodeStruct(v)
		case reflect.Map:
			return encodeMap(vv)
		case reflect.Array, reflect.Slice:
			return encodeSlice(vv)
		default:
			return v, errUnchanged
		}
	}
}

func encodeStruct(v interface{}) (interface{}, error) {
	m := make(map[string]interface{})
	cfg := &mapstructure.DecoderConfig{Result: &m, TagName: "json"}
	if decoder, err := mapstructure.NewDecoder(cfg); err != nil {
		return nil, err
	} else if err := decoder.Decode(v); err != nil {
		return nil, err
	}

	for k, e := range m {
		if encoded, err := encodeValue(e); err == nil {
			m[k] = encoded
		} else if err != errUnchanged {
			return nil, err
		}
	}
	return m, nil
}

func encodeMap(vv reflect.Value) (interface{}, error) {
	switch indirect(vv.Type().Elem()).Kind() {
	case reflect.Struct, reflect.Map, reflect.Array, reflect.Slice, reflect.Interface:
		mv := reflect.MakeMap(reflect.MapOf(vv.Type().Key(), emptyInterfaceType))
		for _, kv := range vv.MapKeys() {
			iv := vv.MapIndex(kv)
			if encoded, err := encodeValue(iv.Interface()); err != nil && err != errUnchanged {
				return nil, err
			} else {
				if encoded == nil {
					// workaround: SetMapIndex treats untyped nil as a map delete
					encoded = (*int)(nil)
				}
				mv.SetMapIndex(kv, reflect.ValueOf(encoded))
			}
		}
		return mv.Interface(), nil
	default:
		return vv.Interface(), errUnchanged
	}
}

func encodeSlice(vv reflect.Value) (interface{}, error) {
	switch indirect(vv.Type().Elem()).Kind() {
	case reflect.Struct, reflect.Map, reflect.Array, reflect.Slice, reflect.Interface:
		l := vv.Len()
		sv := reflect.MakeSlice(reflect.SliceOf(emptyInterfaceType), l, l)
		for i := 0; i < l; i++ {
			if encoded, err := encodeValue(vv.Index(i).Interface()); err != nil && err != errUnchanged {
				return nil, err
			} else {
				if encoded == nil {
					// workaround: Set doesn't like untyped nil
					encoded = (*int)(nil)
				}
				sv.Index(i).Set(reflect.ValueOf(encoded))
			}
		}
		return sv.Interface(), nil
	default:
		return vv.Interface(), errUnchanged
	}
}

func indirect(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		return t.Elem()
	}
	return t
}
