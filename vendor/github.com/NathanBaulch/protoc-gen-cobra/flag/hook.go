package flag

import "github.com/spf13/pflag"

func WithPostSetHook(fs *pflag.FlagSet, name string, hook func()) {
	WithPostSetHookE(fs, name, func() error { hook(); return nil })
}

func WithPostSetHookE(fs *pflag.FlagSet, name string, hook func() error) {
	f := fs.Lookup(name)
	f.Value = &postSetHookValue{f.Value, hook}
}

type postSetHookValue struct {
	pflag.Value
	hook func() error
}

func (v *postSetHookValue) Set(s string) error {
	if err := v.Value.Set(s); err != nil {
		return err
	}
	return v.hook()
}
