package flag

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/spf13/pflag"
)

type mapValue[K comparable, V any] struct {
	value     *map[K]V
	changed   bool
	keyParser func(val string) (K, error)
	valParser func(val string) (V, error)
}

func MapVar[K comparable, V any](fs *pflag.FlagSet, keyParser func(val string) (K, error), valParser func(val string) (V, error), p *map[K]V, name, usage string) {
	fs.Var(&mapValue[K, V]{value: p, keyParser: keyParser, valParser: valParser}, name, usage)
}

func (m *mapValue[K, V]) Set(val string) error {
	sm, ok := m.trySetJSON(val)
	if !ok {
		ss := strings.Split(val, ",")
		sm = make(map[string]string, len(ss))
		for _, pair := range ss {
			kv := strings.SplitN(pair, "=", 2)
			if len(kv) != 2 {
				return fmt.Errorf("%s must be comma separated key=value or a json object", pair)
			}
			sm[kv[0]] = kv[1]
		}
	}

	out := make(map[K]V, len(sm))
	for k, v := range sm {
		if k, err := m.keyParser(k); err != nil {
			return err
		} else if v, err := m.valParser(v); err != nil {
			return err
		} else {
			out[k] = v
		}
	}
	if !m.changed {
		*m.value = out
		m.changed = true
	} else {
		for k, v := range out {
			(*m.value)[k] = v
		}
	}
	return nil
}

func (*mapValue[K, V]) trySetJSON(val string) (map[string]string, bool) {
	if len(val) >= 2 && val[0] == '{' && val[len(val)-1] == '}' {
		var raw map[string]json.RawMessage
		if err := json.Unmarshal([]byte(val), &raw); err != nil {
			return nil, false
		}

		out := make(map[string]string, len(raw))
		for k, v := range raw {
			var str string
			if v[0] == '"' {
				_ = json.Unmarshal(v, &str)
			} else {
				str = string(v)
			}
			out[k] = str
		}
		return out, true
	}

	return nil, false
}

func (*mapValue[K, V]) Type() string { return "map" }

func (*mapValue[K, V]) String() string { return "{}" }

// Deprecated
func ReflectMapVar(fs *pflag.FlagSet, keyParser, valParser func(val string) (interface{}, error), typ string, p interface{}, name, usage string) {
	v := reflect.ValueOf(p)
	if !v.IsValid() || v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Map {
		panic("must be a pointer to a map")
	}
	fs.Var(&reflectMapValue{value: v, typ: typ, keyParser: keyParser, valParser: valParser}, name, usage)
}

// Deprecated
func ParseBool(val string) (interface{}, error) { return strconv.ParseBool(val) }

// Deprecated
func ParseInt32(val string) (interface{}, error) {
	if i, err := strconv.ParseInt(val, 10, 32); err != nil {
		return nil, err
	} else {
		return int32(i), nil
	}
}

// Deprecated
func ParseInt64(val string) (interface{}, error) { return strconv.ParseInt(val, 10, 64) }

// Deprecated
func ParseUint32(val string) (interface{}, error) {
	if i, err := strconv.ParseUint(val, 10, 32); err != nil {
		return nil, err
	} else {
		return uint32(i), nil
	}
}

// Deprecated
func ParseUint64(val string) (interface{}, error) { return strconv.ParseUint(val, 10, 64) }

// Deprecated
func ParseFloat32(val string) (interface{}, error) {
	if i, err := strconv.ParseFloat(val, 32); err != nil {
		return nil, err
	} else {
		return float32(i), nil
	}
}

// Deprecated
func ParseFloat64(val string) (interface{}, error) { return strconv.ParseFloat(val, 64) }

// Deprecated
func ParseString(val string) (interface{}, error) { return val, nil }

// Deprecated
func ParseBytesBase64(val string) (interface{}, error) {
	return base64.RawStdEncoding.DecodeString(strings.TrimRight(val, "="))
}

type reflectMapValue struct {
	value     reflect.Value
	typ       string
	changed   bool
	keyParser func(val string) (interface{}, error)
	valParser func(val string) (interface{}, error)
}

func (s *reflectMapValue) Set(val string) error {
	ss := strings.Split(val, ",")
	v := s.value.Elem()
	out := reflect.MakeMapWithSize(v.Type(), len(ss))
	for _, pair := range ss {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			return fmt.Errorf("%s must be formatted as key=value", pair)
		}
		if k, err := s.keyParser(kv[0]); err != nil {
			return err
		} else if v, err := s.valParser(kv[1]); err != nil {
			return err
		} else {
			out.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v))
		}
	}
	if !s.changed {
		v.Set(out)
		s.changed = true
	} else {
		iter := out.MapRange()
		for iter.Next() {
			v.SetMapIndex(iter.Key(), iter.Value())
		}
	}
	return nil
}

func (s *reflectMapValue) Type() string { return s.typ }

func (*reflectMapValue) String() string { return "<nil>" }
