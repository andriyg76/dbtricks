package main

import (
	"fmt"
	"dbtricks"
	"os"
	"pgdumpsplit"
)


type DataHandler interface {

}

type Dumper interface {

}

func NewDumper(output_file string) (d Dumper, e error) {
	return
}

func main() {
	params, err := pgdumpsplit.ParseParams(os.Args)
	if err != nil {
		fmt.Fprint(os.Stderr, "Error parsing params", err.Error())
		params.PrintUsage()
		os.Exit(2)
	}

	if params.IsHelp() {
		params.PrintUsage()
		os.Exit(0)
	}

	orders := dbtricks.ReadOrders()

	if !orders.IsEmpty() {
		err = orders.WriteOrders()
		if err != nil {
			panic("Error writing orders: " + err.Error())
		}
	}
}
