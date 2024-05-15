package ptypes

import (
	"encoding/base64"
	"strings"

	"github.com/spf13/cast"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func ToTimestamp(v interface{}) (*timestamppb.Timestamp, error) {
	if t, ok := v.(*timestamppb.Timestamp); ok {
		return t, nil
	}
	if t, err := cast.ToTimeE(v); err != nil {
		return nil, err
	} else {
		return timestamppb.New(t), nil
	}
}

func ToDuration(v interface{}) (*durationpb.Duration, error) {
	if d, ok := v.(*durationpb.Duration); ok {
		return d, nil
	}
	if d, err := cast.ToDurationE(v); err != nil {
		return nil, err
	} else {
		return durationpb.New(d), nil
	}
}

func ToDoubleWrapper(v interface{}) (*wrapperspb.DoubleValue, error) {
	if d, ok := v.(*wrapperspb.DoubleValue); ok {
		return d, nil
	}
	if d, err := cast.ToFloat64E(v); err != nil {
		return nil, err
	} else {
		return wrapperspb.Double(d), nil
	}
}

func ToFloatWrapper(v interface{}) (*wrapperspb.FloatValue, error) {
	if f, ok := v.(*wrapperspb.FloatValue); ok {
		return f, nil
	}
	if f, err := cast.ToFloat32E(v); err != nil {
		return nil, err
	} else {
		return wrapperspb.Float(f), nil
	}
}

func ToInt64Wrapper(v interface{}) (*wrapperspb.Int64Value, error) {
	if i, ok := v.(*wrapperspb.Int64Value); ok {
		return i, nil
	}
	if i, err := cast.ToInt64E(v); err != nil {
		return nil, err
	} else {
		return wrapperspb.Int64(i), nil
	}
}

func ToUInt64Wrapper(v interface{}) (*wrapperspb.UInt64Value, error) {
	if i, ok := v.(*wrapperspb.UInt64Value); ok {
		return i, nil
	}
	if i, err := cast.ToUint64E(v); err != nil {
		return nil, err
	} else {
		return wrapperspb.UInt64(i), nil
	}
}

func ToInt32Wrapper(v interface{}) (*wrapperspb.Int32Value, error) {
	if i, ok := v.(*wrapperspb.Int32Value); ok {
		return i, nil
	}
	if i, err := cast.ToInt32E(v); err != nil {
		return nil, err
	} else {
		return wrapperspb.Int32(i), nil
	}
}

func ToUInt32Wrapper(v interface{}) (*wrapperspb.UInt32Value, error) {
	if i, ok := v.(*wrapperspb.UInt32Value); ok {
		return i, nil
	}
	if i, err := cast.ToUint32E(v); err != nil {
		return nil, err
	} else {
		return wrapperspb.UInt32(i), nil
	}
}

func ToBoolWrapper(v interface{}) (*wrapperspb.BoolValue, error) {
	if b, ok := v.(*wrapperspb.BoolValue); ok {
		return b, nil
	}
	if b, err := cast.ToBoolE(v); err != nil {
		return nil, err
	} else {
		return wrapperspb.Bool(b), nil
	}
}

func ToStringWrapper(v interface{}) (*wrapperspb.StringValue, error) {
	if s, ok := v.(*wrapperspb.StringValue); ok {
		return s, nil
	}
	if s, err := cast.ToStringE(v); err != nil {
		return nil, err
	} else {
		return wrapperspb.String(s), nil
	}
}

func ToBytesWrapper(v interface{}) (*wrapperspb.BytesValue, error) {
	if b, ok := v.(*wrapperspb.BytesValue); ok {
		return b, nil
	}
	if b, ok := v.([]byte); ok {
		return wrapperspb.Bytes(b), nil
	}
	if s, err := cast.ToStringE(v); err != nil {
		return nil, err
	} else if v, err := base64.RawStdEncoding.DecodeString(strings.TrimRight(s, "=")); err != nil {
		return nil, err
	} else {
		return wrapperspb.Bytes(v), nil
	}
}
