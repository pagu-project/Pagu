package flag

import "github.com/spf13/pflag"

func BoolPointerVar(fs *pflag.FlagSet, p **bool, name, usage string) {
	v := fs.Bool(name, false, usage)
	WithPostSetHook(fs, name, func() { *p = v })
}

func Int32PointerVar(fs *pflag.FlagSet, p **int32, name, usage string) {
	v := fs.Int32(name, 0, usage)
	WithPostSetHook(fs, name, func() { *p = v })
}

func Int64PointerVar(fs *pflag.FlagSet, p **int64, name, usage string) {
	v := fs.Int64(name, 0, usage)
	WithPostSetHook(fs, name, func() { *p = v })
}

func Uint32PointerVar(fs *pflag.FlagSet, p **uint32, name, usage string) {
	v := fs.Uint32(name, 0, usage)
	WithPostSetHook(fs, name, func() { *p = v })
}

func Uint64PointerVar(fs *pflag.FlagSet, p **uint64, name, usage string) {
	v := fs.Uint64(name, 0, usage)
	WithPostSetHook(fs, name, func() { *p = v })
}

func Float32PointerVar(fs *pflag.FlagSet, p **float32, name, usage string) {
	v := fs.Float32(name, 0, usage)
	WithPostSetHook(fs, name, func() { *p = v })
}

func Float64PointerVar(fs *pflag.FlagSet, p **float64, name, usage string) {
	v := fs.Float64(name, 0, usage)
	WithPostSetHook(fs, name, func() { *p = v })
}

func StringPointerVar(fs *pflag.FlagSet, p **string, name, usage string) {
	v := fs.String(name, "", usage)
	WithPostSetHook(fs, name, func() { *p = v })
}
