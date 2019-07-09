package params

import (
	"errors"
	"flag"
	"fmt"
	"strings"
)

const (
	None  DumpType = -1
	Mysql DumpType = iota
	Pgsql DumpType = iota
)

type DumpType int8

//go:generate command stringer -type DumpType

type Params struct {
	Destination      string
	ChunkSize        int
	Verbose          bool
	Trace            bool
	CleanDestination bool
	tail             []string
	DumpType         DumpType
}

func (i Params) String() string {
	logLevel := "warn"
	if i.Trace {
		logLevel = "trace"
	} else if i.Verbose {
		logLevel = "verbose"
	}
	return fmt.Sprintf("destination=%v chunkSize=%v logLevel=%v cleanDestination=%v dumptype=%v input=%v",
		i.Destination,
		i.ChunkSize,
		logLevel,
		i.CleanDestination,
		i.DumpType,
		i.tail,
	)
}

var OK = errors.New("OK")

func ParseParams(args []string) (error, Params) {
	_params := Params{
		ChunkSize: 2048,
		DumpType:  None,
	}
	var dumptype string
	isHelp := false

	set := flag.NewFlagSet("pgdumpsplit", flag.ContinueOnError)
	set.StringVar(&_params.Destination, "d", "", "Path, where to store splitted files, default - current folder")
	set.IntVar(&_params.ChunkSize, "m", 2048, "Max chunk size of database part, in kb")
	set.BoolVar(&_params.Verbose, "v", false, "Verbose dumping output")
	set.BoolVar(&_params.Trace, "r", false, "Trace dumping output")
	set.BoolVar(&_params.CleanDestination, "c", false, "Clean destination")
	set.StringVar(&dumptype, "t", "", "Dump type PGSQL|MYSQL")
	set.BoolVar(&isHelp, "h", false, "Help")

	var err = set.Parse(args[1:])
	if err != nil {
		return err, _params
	}
	if isHelp {
		set.PrintDefaults()
		return OK, _params
	}

	switch strings.ToUpper(dumptype) {
	case "":
		return errors.New("dumptype is not set"), _params
	case "PGSQL":
		_params.DumpType = Pgsql
		break
	case "MYSQL":
		_params.DumpType = Mysql
		break
	default:
		return errors.New("DumpType " + dumptype + " is not supported"), _params
	}
	_params.tail = set.Args()
	if len(_params.tail) > 1 {
		return errors.New("Unnecessary params given: " + strings.Join(_params.tail, ", ")), _params
	}
	return nil, _params
}

func (i Params) File() string {
	if len(i.tail) == 1 {
		return i.tail[0]
	} else {
		return ""
	}
}
