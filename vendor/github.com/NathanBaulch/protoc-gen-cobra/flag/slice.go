package flag

import (
	"encoding/json"
	"strings"

	"github.com/spf13/pflag"
)

type sliceValue[T any] struct {
	value   *[]T
	changed bool
	parser  func(val string) (T, error)
}

func SliceVar[T any](fs *pflag.FlagSet, parser func(val string) (T, error), p *[]T, name, usage string) {
	fs.Var(&sliceValue[T]{value: p, parser: parser}, name, usage)
}

func (s *sliceValue[T]) Set(val string) error {
	ss, ok := s.trySetJSON(val)
	if !ok {
		ss = strings.Split(val, ",")
	}

	out := make([]T, len(ss))
	for i, v := range ss {
		var err error
		if out[i], err = s.parser(v); err != nil {
			return err
		}
	}
	if !s.changed {
		*s.value = out
		s.changed = true
	} else {
		*s.value = append(*s.value, out...)
	}
	return nil
}

func (*sliceValue[T]) trySetJSON(val string) ([]string, bool) {
	if len(val) >= 2 {
		if val[0] == '{' && val[len(val)-1] == '}' && json.Valid([]byte(val)) {
			return []string{val}, true
		} else if val[0] == '[' && val[len(val)-1] == ']' {
			var raw []json.RawMessage
			if err := json.Unmarshal([]byte(val), &raw); err != nil {
				return nil, false
			}

			out := make([]string, len(raw))
			for i, v := range raw {
				if v[0] == '"' {
					_ = json.Unmarshal(v, &out[i])
				} else {
					out[i] = string(v)
				}
			}
			return out, true
		}
	}

	return nil, false
}

func (*sliceValue[T]) Type() string { return "slice" }

func (*sliceValue[T]) String() string { return "[]" }

func Uint32SliceVar(fs *pflag.FlagSet, p *[]uint32, name, usage string) {
	SliceVar[uint32](fs, ParseUint32E, p, name, usage)
}

func Uint64SliceVar(fs *pflag.FlagSet, p *[]uint64, name, usage string) {
	SliceVar[uint64](fs, ParseUint64E, p, name, usage)
}
