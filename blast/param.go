package blast

import (
	"bytes"
	"fmt"
	"io"
	"runtime"

	"github.com/TuftsBCB/io/fasta"
	"github.com/TuftsBCB/seq"
)

// Blaster represents values that can execute a BLAST search. This package
// provides some slim implementations of this interface for a couple variations
// of BLAST. Clients requiring access to some of BLAST's more sophisticated
// options should provide their own Blaster.
type Blaster interface {
	// Executable should return the blast executable to run.
	Executable() string

	// CmdArgs should return a list of command line flags to pass to the
	// blast executable. This list must not include the `-outfmt` flag,
	// since clients of this interface may set it in order to retrieve
	// results in an expected format.
	CmdArgs() []string

	// Stdin, when not nil, will be used for the stdin of the blast process.
	Stdin() io.Reader
}

// Query is a generic blaster for any type of BLAST search. It provides a
// thin wrapper around setting command line flags to pass to a BLAST
// executable.
type Query struct {
	// The BLAST executable to use.
	Exec    string
	queries []seq.Sequence
	flags
}

// NewQuery constructs a generic blast search with default parameters.
// Parameters can be overridden using the `SetFlag` method.
//
// Note that `queries` may have length 0. If it does, then the obligation is
// on the caller to set the `-query` flag (or provide some other means of
// giving BLAST a search query).
//
// This also sets the `-num_threads` flag to the number of logical CPUs
// on your machine.
func NewQuery(exec string, queries []seq.Sequence, database string) *Query {
	b := &Query{
		Exec:    exec,
		queries: queries,
		flags:   make(flags, 0),
	}
	b.SetFlag("db", database)
	b.SetFlag("num_threads", runtime.NumCPU())
	return b
}

// NewBlastp calls NewQuery with "blastp" as the executable.
func NewBlastp(queries []seq.Sequence, database string) *Query {
	return NewQuery("blastp", queries, database)
}

// NewBlastn calls NewQuery with "blastn" as the executable.
func NewBlastn(queries []seq.Sequence, database string) *Query {
	return NewQuery("blastn", queries, database)
}

// SetFlag adds a command line switch (without the proceeding "-") to the
// set of blastp arguments. `value` should be a string, integer, float, bool
// or other type with an appropriate `Stringer` implementation that results
// in a valid command line flag value.
//
// If `value` is `false`, then the flag is removed from the blastp arguments.
func (b *Query) SetFlag(name string, value interface{}) {
	b.flags.set(name, value)
}

func (b *Query) Executable() string {
	return b.Exec
}

func (b *Query) Stdin() io.Reader {
	return queryReader(b.queries)
}

type flags map[string]string

func (fs flags) set(name string, v interface{}) {
	switch v := v.(type) {
	case bool:
		if v {
			fs[name] = ""
		} else {
			delete(fs, name)
		}
	default:
		fs[name] = fmt.Sprintf("%v", v)
	}
}

func (fs flags) CmdArgs() []string {
	args := make([]string, 0, len(fs)*2)
	for name, val := range fs {
		args = append(args, "-"+name)
		if len(val) > 0 {
			args = append(args, val)
		}
	}
	return args
}

func queryReader(queries []seq.Sequence) io.Reader {
	if len(queries) == 0 {
		return nil
	}

	buf := new(bytes.Buffer)
	w := fasta.NewWriter(buf)
	if err := w.WriteAll(queries); err != nil {
		// I don't think this is possible unless the underlying byte buffer
		// becomes too big for it to grow any more.
		panic(err)
	}
	return buf
}
