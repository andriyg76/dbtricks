package main

import (
	"fmt"
	"os"
	"pgdumpsplit/params"
	"log"
	"orders"
	"bufio"
)


type DataHandler interface {

}

type Dumper interface {

}

func main() {
	params := params.ParseParams(os.Args)
	if params.Error() != nil {
		fmt.Fprintln(os.Stderr, "Error parsing params ", params.Error())
		os.Exit(2)
	}

	if params.IsHelp() {
		params.PrintUsage(os.Stdout)
		os.Exit(0)
	}

	var file *os.File
	if params.File() == "" || params.File() == "-" {
		file = os.Stdin
	} else {
		file, err := os.OpenFile(params.File(), os.O_RDONLY, os.ModePerm)
		if err != nil {
			log.Fatal("Can't open file ", params.File(), " for read")
		}
		defer file.Close()
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if scanner.Err() != nil {
		log.Fatal("Error reading input file: ", scanner.Err())
	}

	orders := orders.ReadOrders(params.Destination())

	if !orders.IsEmpty() {
		err := orders.WriteOrders()
		if err != nil {
			panic("Error writing orders: " + err.Error())
		}
	}
}
