package flag

import (
	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/NathanBaulch/protoc-gen-cobra/ptypes"
)

func TimestampVar(fs *pflag.FlagSet, p **timestamppb.Timestamp, name, usage string) {
	v := fs.String(name, "", usage)
	WithPostSetHookE(fs, name, func() (err error) { *p, err = ptypes.ToTimestamp(v); return })
}

func TimestampSliceVar(fs *pflag.FlagSet, p *[]*timestamppb.Timestamp, name, usage string) {
	SliceVar[*timestamppb.Timestamp](fs, ParseTimestampE, p, name, usage)
}

func ParseTimestampE(val string) (*timestamppb.Timestamp, error) { return ptypes.ToTimestamp(val) }

// Deprecated
func ParseTimestamp(val string) (interface{}, error) { return ptypes.ToTimestamp(val) }
