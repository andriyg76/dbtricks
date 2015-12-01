package pgdumpsplit

import (
	"flag"
	"encoding/json"
)

type params struct {
	V_Destination      string `json:"Destination"`
	V_ChunkSize        int    `json:"ChunkSize"`
	V_IsVerbose        bool   `json:"IsVerbose"`
	V_CleanDestination bool   `json:"CleanDestination"`
	V_IsHelp           bool   `json:"IsHelp"`
	params             *flag.FlagSet
}

type Params interface {
	Destination() string
	ChunkSize() int
	IsVerbose() bool
	CleanDestination() bool
	IsHelp() bool
	PrintUsage()
	String() string
}

func ParseParams(args []string) (Params, error) {
	_params := params{
		V_ChunkSize: 2048,
	}

	set := flag.NewFlagSet("Split database dump file to a chunks.", flag.ContinueOnError)
	set.StringVar(&_params.V_Destination, "d", "", "Path, where to store splitted files")
	set.IntVar(&_params.V_ChunkSize, "m", 2048, "Max chunk size of database part, in kb")
	set.BoolVar(&_params.V_IsVerbose, "v", false, "Verbose dumping output")
	set.BoolVar(&_params.V_CleanDestination, "c", false, "Clean destination")
	set.BoolVar(&_params.V_IsHelp, "h", false, "Help")

	err := set.Parse(args)
	if err != nil {
		return nil, err
	}
	_params.params = set
	return &_params, nil
}

func (i* params) Destination() string {
	return i.V_Destination
}

func (i* params) ChunkSize() int {
	return i.V_ChunkSize
}

func (i* params) IsVerbose() bool {
	return i.V_IsVerbose
}

func (i* params) CleanDestination() bool {
	return i.V_CleanDestination
}

func (i* params) IsHelp() bool {
	return i.V_IsHelp
}

func (i* params) PrintUsage() {
	i.params.Usage()
}

func (i* params) String() string {
	val, _ := json.Marshal(i)
	return string(val)
}