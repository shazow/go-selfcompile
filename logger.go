package selfcompile

import (
	"io"
	"io/ioutil"
)
import stdlog "log"

var logger *stdlog.Logger

// SetLogger replaces the package's log writer
func SetLogger(w io.Writer) {
	flags := stdlog.Flags()
	prefix := "[selfcompile] "
	logger = stdlog.New(w, prefix, flags)
}

func init() {
	SetLogger(ioutil.Discard)
}
