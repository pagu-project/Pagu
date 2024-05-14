package flag

import (
	"encoding/base64"
	"io"
	"os"
	"strings"

	"github.com/spf13/pflag"
)

var stdin io.Reader = os.Stdin

func BytesBase64Var(fs *pflag.FlagSet, p *[]byte, name, usage string) {
	v := fs.String(name, "", usage)
	hook := func() (err error) {
		if *v == "-" {
			*p, err = io.ReadAll(stdin)
		} else {
			*p, err = ParseBytesBase64E(*v)
		}
		return
	}
	WithPostSetHookE(fs, name, hook)
}

func BytesBase64SliceVar(fs *pflag.FlagSet, p *[][]byte, name, usage string) {
	SliceVar[[]byte](fs, ParseBytesBase64E, p, name, usage)
}

func ParseBytesBase64E(val string) ([]byte, error) {
	return base64.RawStdEncoding.DecodeString(strings.TrimRight(val, "="))
}
