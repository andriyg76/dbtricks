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

func main() {
	params := pgdumpsplit.ParseParams(os.Args)
	if params.Error() != nil {
		fmt.Fprintln(os.Stderr, "Error parsing params ", params.Error())
		os.Exit(2)
	}

	if params.IsHelp() {
		params.PrintUsage(os.Stdout)
		os.Exit(0)
	}

	orders := dbtricks.ReadOrders()

	if !orders.IsEmpty() {
		err := orders.WriteOrders()
		if err != nil {
			panic("Error writing orders: " + err.Error())
		}
	}
}
