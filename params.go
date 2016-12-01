package main

import (
	"flag"
	"errors"
	"strings"
	"fmt"
)

const (
	DUMPTYPE_NONE Dumptype = -1
	DUMPTYPE_MYSQL Dumptype = iota
	DUMPTYPE_PGSQL Dumptype = iota
)

type Dumptype int8
//go:generate command stringer -type Dumptype

type params struct {
	destination      string
	chunkSize        int
	verbose          bool
	trace            bool
	cleanDestination bool
	tail             []string
	dumptype         Dumptype
}

func (params params) String() string {
	logLevel := "warn"
	if params.trace  {
		logLevel = "trace"
	} else if params.verbose  {
		logLevel = "verbose"
	}
	return fmt.Sprintf("destination=%v chunkSize=%v logLevel=%v cleanDestination=%v dumptype=%v input=%v",
		params.destination,
		params.chunkSize,
		logLevel,
		params.cleanDestination,
		params.dumptype,
		params.tail,
	)
}

var OK = errors.New("OK")

func parseParams(args []string) (error, params) {
	_params := params{
		chunkSize: 2048,
		dumptype: -1,
	}
	var dumptype string
	isHelp := false

	set := flag.NewFlagSet("pgdumpsplit", flag.ContinueOnError)
	set.StringVar(&_params.destination, "d", "", "Path, where to store splitted files, default - current folder")
	set.IntVar(&_params.chunkSize, "m", 2048, "Max chunk size of database part, in kb")
	set.BoolVar(&_params.verbose, "v", false, "Verbose dumping output")
	set.BoolVar(&_params.trace, "r", false, "Trace dumping output")
	set.BoolVar(&_params.cleanDestination, "c", false, "Clean destination")
	set.StringVar(&dumptype, "t", "", "Dumptype PGSQL|MYSQL")
	set.BoolVar(&isHelp, "h", false, "Help")

	var error = set.Parse(args[1:])
	if error != nil {
		return error, _params
	}
	if (isHelp) {
		set.PrintDefaults()
		return OK, _params
	}

	switch strings.ToUpper(dumptype) {
	case "":
		return errors.New("Dumptype is not set"), _params
	case "PGSQL":
		_params.dumptype = DUMPTYPE_PGSQL
		break
	case "MYSQL":
		_params.dumptype = DUMPTYPE_MYSQL
		break
	default:
		return errors.New("Dumptype " + dumptype + " is not supported"), _params
	}
	_params.tail = set.Args()
	if len(_params.tail) > 1 {
		return errors.New("Unnecessary params given: " + strings.Join(_params.tail, ", ")), _params
	}
	return nil, _params
}

func (i params) File() string {
	if len(i.tail) == 1 {
		return i.tail[0]
	} else {
		return ""
	}
}
