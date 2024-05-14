package flag

import (
	"strconv"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func ParseBoolE(val string) (bool, error) { return strconv.ParseBool(val) }

func ParseInt32E(val string) (int32, error) {
	if i, err := strconv.ParseInt(val, 10, 32); err != nil {
		return 0, err
	} else {
		return int32(i), nil
	}
}

func ParseInt64E(val string) (int64, error) { return strconv.ParseInt(val, 10, 64) }

func ParseUint32E(val string) (uint32, error) {
	if i, err := strconv.ParseUint(val, 10, 32); err != nil {
		return 0, err
	} else {
		return uint32(i), nil
	}
}

func ParseUint64E(val string) (uint64, error) { return strconv.ParseUint(val, 10, 64) }

func ParseFloat32E(val string) (float32, error) {
	if i, err := strconv.ParseFloat(val, 32); err != nil {
		return 0, err
	} else {
		return float32(i), nil
	}
}

func ParseFloat64E(val string) (float64, error) { return strconv.ParseFloat(val, 64) }

func ParseStringE(val string) (string, error) { return val, nil }

func ParseMessageE[T proto.Message](val string) (T, error) {
	var t T
	t = t.ProtoReflect().New().Interface().(T)
	err := protojson.Unmarshal([]byte(val), t)
	return t, err
}
