package flag

import (
	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/NathanBaulch/protoc-gen-cobra/ptypes"
)

func BoolWrapperVar(fs *pflag.FlagSet, p **wrapperspb.BoolValue, name, usage string) {
	v := fs.String(name, "", usage)
	WithPostSetHookE(fs, name, func() (err error) { *p, err = ptypes.ToBoolWrapper(v); return })
}

func BoolWrapperSliceVar(fs *pflag.FlagSet, p *[]*wrapperspb.BoolValue, name, usage string) {
	v := fs.BoolSlice(name, nil, usage)
	hook := func() {
		out := make([]*wrapperspb.BoolValue, len(*v))
		for i, item := range *v {
			out[i] = wrapperspb.Bool(item)
		}
		*p = out
	}
	WithPostSetHook(fs, name, hook)
}

func ParseBoolWrapperE(val string) (*wrapperspb.BoolValue, error) { return ptypes.ToBoolWrapper(val) }

// Deprecated
func ParseBoolWrapper(val string) (interface{}, error) { return ptypes.ToBoolWrapper(val) }

func Int32WrapperVar(fs *pflag.FlagSet, p **wrapperspb.Int32Value, name, usage string) {
	v := fs.Int32(name, 0, usage)
	WithPostSetHook(fs, name, func() { *p = wrapperspb.Int32(*v) })
}

func Int32WrapperSliceVar(fs *pflag.FlagSet, p *[]*wrapperspb.Int32Value, name, usage string) {
	v := fs.Int32Slice(name, nil, usage)
	hook := func() {
		out := make([]*wrapperspb.Int32Value, len(*v))
		for i, item := range *v {
			out[i] = wrapperspb.Int32(item)
		}
		*p = out
	}
	WithPostSetHook(fs, name, hook)
}

func ParseInt32WrapperE(val string) (*wrapperspb.Int32Value, error) {
	return ptypes.ToInt32Wrapper(val)
}

// Deprecated
func ParseInt32Wrapper(val string) (interface{}, error) { return ptypes.ToInt32Wrapper(val) }

func Int64WrapperVar(fs *pflag.FlagSet, p **wrapperspb.Int64Value, name, usage string) {
	v := fs.Int64(name, 0, usage)
	WithPostSetHook(fs, name, func() { *p = wrapperspb.Int64(*v) })
}

func Int64WrapperSliceVar(fs *pflag.FlagSet, p *[]*wrapperspb.Int64Value, name, usage string) {
	v := fs.Int64Slice(name, nil, usage)
	hook := func() {
		out := make([]*wrapperspb.Int64Value, len(*v))
		for i, item := range *v {
			out[i] = wrapperspb.Int64(item)
		}
		*p = out
	}
	WithPostSetHook(fs, name, hook)
}

func ParseInt64WrapperE(val string) (*wrapperspb.Int64Value, error) {
	return ptypes.ToInt64Wrapper(val)
}

// Deprecated
func ParseInt64Wrapper(val string) (interface{}, error) { return ptypes.ToInt64Wrapper(val) }

func UInt32WrapperVar(fs *pflag.FlagSet, p **wrapperspb.UInt32Value, name, usage string) {
	v := fs.Uint32(name, 0, usage)
	WithPostSetHook(fs, name, func() { *p = wrapperspb.UInt32(*v) })
}

func UInt32WrapperSliceVar(fs *pflag.FlagSet, p *[]*wrapperspb.UInt32Value, name, usage string) {
	var v []uint32
	Uint32SliceVar(fs, &v, name, usage)
	hook := func() {
		out := make([]*wrapperspb.UInt32Value, len(v))
		for i, item := range v {
			out[i] = wrapperspb.UInt32(item)
		}
		*p = out
	}
	WithPostSetHook(fs, name, hook)
}

func ParseUInt32WrapperE(val string) (*wrapperspb.UInt32Value, error) {
	return ptypes.ToUInt32Wrapper(val)
}

// Deprecated
func ParseUInt32Wrapper(val string) (interface{}, error) { return ptypes.ToUInt32Wrapper(val) }

func UInt64WrapperVar(fs *pflag.FlagSet, p **wrapperspb.UInt64Value, name, usage string) {
	v := fs.Uint64(name, 0, usage)
	WithPostSetHook(fs, name, func() { *p = wrapperspb.UInt64(*v) })
}

func UInt64WrapperSliceVar(fs *pflag.FlagSet, p *[]*wrapperspb.UInt64Value, name, usage string) {
	var v []uint64
	Uint64SliceVar(fs, &v, name, usage)
	hook := func() {
		out := make([]*wrapperspb.UInt64Value, len(v))
		for i, item := range v {
			out[i] = wrapperspb.UInt64(item)
		}
		*p = out
	}
	WithPostSetHook(fs, name, hook)
}

func ParseUInt64WrapperE(val string) (*wrapperspb.UInt64Value, error) {
	return ptypes.ToUInt64Wrapper(val)
}

// Deprecated
func ParseUInt64Wrapper(val string) (interface{}, error) { return ptypes.ToUInt64Wrapper(val) }

func FloatWrapperVar(fs *pflag.FlagSet, p **wrapperspb.FloatValue, name, usage string) {
	v := fs.Float32(name, 0, usage)
	WithPostSetHook(fs, name, func() { *p = wrapperspb.Float(*v) })
}

func FloatWrapperSliceVar(fs *pflag.FlagSet, p *[]*wrapperspb.FloatValue, name, usage string) {
	v := fs.Float32Slice(name, nil, usage)
	hook := func() {
		out := make([]*wrapperspb.FloatValue, len(*v))
		for i, item := range *v {
			out[i] = wrapperspb.Float(item)
		}
		*p = out
	}
	WithPostSetHook(fs, name, hook)
}

func ParseFloatWrapperE(val string) (*wrapperspb.FloatValue, error) {
	return ptypes.ToFloatWrapper(val)
}

// Deprecated
func ParseFloatWrapper(val string) (interface{}, error) { return ptypes.ToFloatWrapper(val) }

func DoubleWrapperVar(fs *pflag.FlagSet, p **wrapperspb.DoubleValue, name, usage string) {
	v := fs.Float64(name, 0, usage)
	WithPostSetHook(fs, name, func() { *p = wrapperspb.Double(*v) })
}

func DoubleWrapperSliceVar(fs *pflag.FlagSet, p *[]*wrapperspb.DoubleValue, name, usage string) {
	v := fs.Float64Slice(name, nil, usage)
	hook := func() {
		out := make([]*wrapperspb.DoubleValue, len(*v))
		for i, item := range *v {
			out[i] = wrapperspb.Double(item)
		}
		*p = out
	}
	WithPostSetHook(fs, name, hook)
}

func ParseDoubleWrapperE(val string) (*wrapperspb.DoubleValue, error) {
	return ptypes.ToDoubleWrapper(val)
}

// Deprecated
func ParseDoubleWrapper(val string) (interface{}, error) { return ptypes.ToDoubleWrapper(val) }

func StringWrapperVar(fs *pflag.FlagSet, p **wrapperspb.StringValue, name, usage string) {
	v := fs.String(name, "", usage)
	WithPostSetHook(fs, name, func() { *p = wrapperspb.String(*v) })
}

func StringWrapperSliceVar(fs *pflag.FlagSet, p *[]*wrapperspb.StringValue, name, usage string) {
	v := fs.StringSlice(name, nil, usage)
	hook := func() {
		out := make([]*wrapperspb.StringValue, len(*v))
		for i, item := range *v {
			out[i] = wrapperspb.String(item)
		}
		*p = out
	}
	WithPostSetHook(fs, name, hook)
}

func ParseStringWrapperE(val string) (*wrapperspb.StringValue, error) {
	return ptypes.ToStringWrapper(val)
}

// Deprecated
func ParseStringWrapper(val string) (interface{}, error) { return ptypes.ToStringWrapper(val) }

func BytesBase64WrapperVar(fs *pflag.FlagSet, p **wrapperspb.BytesValue, name, usage string) {
	var v []byte
	BytesBase64Var(fs, &v, name, usage)
	WithPostSetHook(fs, name, func() { *p = wrapperspb.Bytes(v) })
}

func BytesBase64WrapperSliceVar(fs *pflag.FlagSet, p *[]*wrapperspb.BytesValue, name, usage string) {
	var v [][]byte
	BytesBase64SliceVar(fs, &v, name, usage)
	hook := func() {
		out := make([]*wrapperspb.BytesValue, len(v))
		for i, item := range v {
			out[i] = wrapperspb.Bytes(item)
		}
		*p = out
	}
	WithPostSetHook(fs, name, hook)
}

func ParseBytesBase64WrapperE(val string) (*wrapperspb.BytesValue, error) {
	return ptypes.ToBytesWrapper(val)
}

// Deprecated
func ParseBytesBase64Wrapper(val string) (interface{}, error) { return ptypes.ToBytesWrapper(val) }
