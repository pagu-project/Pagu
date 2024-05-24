package flag

import (
	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/NathanBaulch/protoc-gen-cobra/ptypes"
)

func DurationVar(fs *pflag.FlagSet, p **durationpb.Duration, name, usage string) {
	v := fs.String(name, "", usage)
	WithPostSetHookE(fs, name, func() (err error) { *p, err = ptypes.ToDuration(v); return })
}

func DurationSliceVar(fs *pflag.FlagSet, p *[]*durationpb.Duration, name, usage string) {
	SliceVar[*durationpb.Duration](fs, ParseDurationE, p, name, usage)
}

func ParseDurationE(val string) (*durationpb.Duration, error) { return ptypes.ToDuration(val) }

// Deprecated
func ParseDuration(val string) (interface{}, error) { return ptypes.ToDuration(val) }
