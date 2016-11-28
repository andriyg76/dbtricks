package main

import (
	"flag"
	"os"
	"errors"
	"strings"
)

const (
	DUMPTYPE_NONE = Dumptype(-1)
	DUMPTYPE_MYSQL = Dumptype(iota)
	DUMPTYPE_PGSQL = Dumptype(iota)
)

type Dumptype int8

type params struct {
	destination      string
	chunkSize        int
	isVerbose        bool
	cleanDestination bool
	tail             []string
	dumptype         Dumptype
}

type Params interface {
	Destination() string
	ChunkSize() int
	IsVerbose() bool
	CleanDestination() bool
	File() string
	PrintUsage(file *os.File)
	Dumptype() Dumptype
}

func parseParams(args []string) (error, *params) {
	_params := params{
		chunkSize: 2048,
		dumptype: -1,
	}
	var dumptype string
	isHelp := false

	set := flag.NewFlagSet("pgdumpsplit", flag.ContinueOnError)
	set.StringVar(&_params.destination, "d", "", "Path, where to store splitted files, default - current folder")
	set.IntVar(&_params.chunkSize, "m", 2048, "Max chunk size of database part, in kb")
	set.BoolVar(&_params.isVerbose, "v", false, "Verbose dumping output")
	set.BoolVar(&_params.cleanDestination, "c", false, "Clean destination")
	set.StringVar(&dumptype, "t", "", "Dumptype PGSQL|MYSQL")
	set.BoolVar(&isHelp, "h", false, "Help")

	var error = set.Parse(args[1:])
	if error != nil {
		return error, nil
	}
	if (isHelp) {
		set.PrintDefaults()
		return nil, nil
	}

	switch strings.ToUpper(dumptype) {
	case "":
		return errors.New("Dumptype is not set"), nil
	case "PGSQL":
		_params.dumptype = DUMPTYPE_PGSQL
		break
	case "MYSQL":
		_params.dumptype = DUMPTYPE_PGSQL
		break
	default:
		return errors.New("Dumptype " + dumptype + " is not supported"), nil
	}
	_params.tail = set.Args()
	if len(_params.tail) > 1 {
		return errors.New("Unnecessary params given: " + strings.Join(_params.tail, ", ")), nil
	}
	return nil, &_params
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

func (i *params) File() string {
	if len(i.tail) == 1 {
		return i.tail[0]
	} else {
		return ""
	}
}

func (i *params) Dumptype() Dumptype {
	return i.dumptype
}