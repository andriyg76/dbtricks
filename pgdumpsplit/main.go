package main

import (
	"fmt"
	"os"
	"dbtricks/params"
	"log"
	"dbtricks/orders"
	"bufio"
	"pg/dumpsplit"
	"strings"
	"io"
)

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
		var err error
		file, err = os.OpenFile(params.File(), os.O_RDONLY, os.ModePerm)
		if err != nil {
			log.Fatal("Can't open file ", params.File(), " for read")
		}
		defer file.Close()
	}
	reader := bufio.NewReader(file)

	if params.Destination() != "" {
		if err := os.Chdir(params.Destination()); err != nil {
			log.Fatal("Main: Can't change dir to: ", params.Destination())
		}
	}
	orders := orders.ReadOrders(params.Destination())

	splitter, _ := dumpsplit.NewSplitter(orders, params.ChunkSize())
	defer splitter.Close()

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			splitter.Close()
			log.Fatal("Main: Can't read iput file: " + err.Error())
		}
		line = strings.TrimRight(line, "\n\r")
		err = splitter.HandleLine(line)
		if err != nil {
			splitter.Close()
			log.Fatal("Main: Can't handle line: [" + line + "]: " + err.Error())
		}
	}
	if err := splitter.Error(); err != nil {
		splitter.Close()
		log.Fatal("Error reading input file: ", err)
	}

	if !orders.IsEmpty() {
		err := orders.WriteOrders()
		if err != nil {
			panic("Error writing orders: " + err.Error())
		}
	}
}
