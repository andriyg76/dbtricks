package pgdumpsplit

import (
	"flag"
	"os"
	"errors"
	"strings"
)

type params struct {
	destination      string
	chunkSize        int
	isVerbose        bool
	cleanDestination bool
	isHelp           bool
	tail             []string
	error            error
	set              *flag.FlagSet
}

type Params interface {
	Destination() string
	ChunkSize() int
	IsVerbose() bool
	CleanDestination() bool
	IsHelp() bool
	File() string
	Error() error
	PrintUsage(file *os.File)
}

func ParseParams(args []string) (r_val Params) {
	_params := params{
		chunkSize: 2048,
	}
	r_val = &_params

	_params.set = flag.NewFlagSet("pgdumpsplit", flag.ContinueOnError)
	_params.set.StringVar(&_params.destination, "d", "", "Path, where to store splitted files")
	_params.set.IntVar(&_params.chunkSize, "m", 2048, "Max chunk size of database part, in kb")
	_params.set.BoolVar(&_params.isVerbose, "v", false, "Verbose dumping output")
	_params.set.BoolVar(&_params.cleanDestination, "c", false, "Clean destination")
	_params.set.BoolVar(&_params.isHelp, "h", false, "Help")

	_params.error = _params.set.Parse(args[1:])
	if _params.error != nil {
		return
	}
	_params.tail = _params.set.Args()
	if len(_params.tail) > 1 {
		_params.set.PrintDefaults()
		_params.error = errors.New("Unnecessary params given: " + strings.Join(_params.tail, ", "))
		return
	}
	return
}

func (i *params) Destination() string {
	return i.destination
}

func (i *params) ChunkSize() int {
	return i.chunkSize
}

func (i *params) IsVerbose() bool {
	return i.isVerbose
}

func (i *params) CleanDestination() bool {
	return i.cleanDestination
}

func (i *params) IsHelp() bool {
	return i.isHelp
}

func (i *params) PrintUsage(file *os.File) {
	i.set.SetOutput(file)
	i.set.PrintDefaults()
}

func (i *params) Error() error {
	return i.error
}

func (i *params) File() string {
	if len(i.tail) == 1 {
		return i.tail[0]
	} else {
		return ""
	}
}