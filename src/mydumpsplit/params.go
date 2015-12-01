package mydumpsplit

import "flag"

type params struct {
	destination string
	chunk_size  int
	is_verbose  bool
	clean_dest  bool
	is_help     bool
	params      *flag.FlagSet
}

type Params interface {
	Destination() string
	ChunkSize() int
	IsVerbose() bool
	CleanDestination() bool
	IsHelp() bool
	PrintUsage()
}

func ParseParams(args []string) (Params, error) {
	_params := params{
		chunk_size: 2048,
	}

	set := flag.NewFlagSet("Split database dump file to a chunks.", flag.ContinueOnError)
	set.StringVar(&_params.destination, "d", "", "Path, where to store splitted files")
	set.IntVar(&_params.chunk_size, "m", 2048, "Max chunk size of database part, in kb")
	set.BoolVar(&_params.is_verbose, "v", false, "Verbose dumping output")
	set.BoolVar(&_params.clean_dest, "c", false, "Clean destination")
	set.BoolVar(&_params.is_help, "h", false, "Help")

	err := set.Parse(args)
	if err != nil {
		return nil, err
	}
	_params.params = set
	return &_params, nil
}

func (i* params) Destination() string {
	return i.destination
}

func (i* params) ChunkSize() int {
	return i.chunk_size
}

func (i* params) IsVerbose() bool {
	return i.is_verbose
}

func (i* params) CleanDestination() bool {
	return i.clean_dest
}

func (i* params) IsHelp() bool {
	return i.is_help
}

func (i* params) PrintUsage() {
	i.params.Usage()
}