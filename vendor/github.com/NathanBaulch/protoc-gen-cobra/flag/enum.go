package flag

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/pflag"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func EnumVar[T protoreflect.Enum](fs *pflag.FlagSet, p *T, name, usage string) {
	v := fs.String(name, "", usage)
	WithPostSetHookE(fs, name, func() (err error) { *p, err = ParseEnumE[T](*v); return })
}

func EnumPointerVar[T protoreflect.Enum](fs *pflag.FlagSet, p **T, name, usage string) {
	v := fs.String(name, "", usage)
	WithPostSetHookE(fs, name, func() error {
		if e, err := ParseEnumE[T](*v); err != nil {
			return err
		} else {
			*p = &e
			return nil
		}
	})
}

func EnumSliceVar[T protoreflect.Enum](fs *pflag.FlagSet, p *[]T, name, usage string) {
	SliceVar[T](fs, ParseEnumE[T], p, name, usage)
}

func ParseEnumE[T protoreflect.Enum](val string) (T, error) {
	var t T
	if v := parseEnum[T](val); v != nil {
		return t.Type().New(v.Number()).(T), nil
	} else {
		return t, fmt.Errorf("unable to parse enum: %s", val)
	}
}

func parseEnum[T protoreflect.Enum](val string) protoreflect.EnumValueDescriptor {
	var t T
	vals := t.Descriptor().Values()
	if i, err := strconv.ParseInt(val, 0, 32); err == nil {
		if v := vals.ByNumber(protoreflect.EnumNumber(i)); v != nil {
			return v
		}
	} else if v := vals.ByName(protoreflect.Name(val)); v != nil {
		return v
	} else {
		for i := 0; i < vals.Len(); i++ {
			if v := vals.Get(i); strings.EqualFold(string(v.Name()), val) {
				return v
			}
		}
	}
	return nil
}
