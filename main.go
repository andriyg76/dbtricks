package main

import (
	"bufio"
	"fmt"
	"github.com/andriyg76/dbtricks/orders"
	"os"
	"log"
	pg_dumpsplit "github.com/andriyg76/dbtricks/pg/dumpsplit"
	mysql_dumpsplit "github.com/andriyg76/dbtricks/mysql/dumpsplit"
	"io"
	"strings"
	"github.com/andriyg76/dbtricks/dumper"
)

func main() {
	err, params := parseParams(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing params ", err)
		os.Exit(2)
	}

	if params == nil {
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

	var splitter dumper.Dumper
	var error error
	if (params.Dumptype() == DUMPTYPE_PGSQL) {
		splitter, error = pg_dumpsplit.NewSplitter(orders, params.ChunkSize())
	} else if (params.Dumptype() == DUMPTYPE_MYSQL) {
		splitter, error = mysql_dumpsplit.NewSplitter(orders, params.ChunkSize())
	} else {
		log.Fatal("Unsupported dumptype: ", params.Dumptype())
	}
	if error != nil {
		log.Fatal("Can't initialize datasplitter type: ", params.Dumptype())
	}
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
