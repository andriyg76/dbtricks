package main

import (
	"bufio"
	"fmt"
	"github.com/andriyg76/dbtricks/orders"
	"os"
	log "github.com/andriyg76/glogger"
	pg_dumpsplit "github.com/andriyg76/dbtricks/pg/dumpsplit"
	mysql_dumpsplit "github.com/andriyg76/dbtricks/mysql/dumpsplit"
	"io"
	"strings"
	"github.com/andriyg76/dbtricks/splitter"
)

func main() {
	err, params := parseParams(os.Args)
	if err == OK {
		os.Exit(0)
	} else if err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing params ", err)
		os.Exit(2)
	}

	if params.trace {
		log.SetLevel(log.TRACE)
	} else if params.verbose {
		log.SetLevel(log.DEBUG)
	}
	log.Debug("Params %s", params)

	var file *os.File
	if params.File() == "" || params.File() == "-" {
		log.Trace("Reading stdin")
		file = os.Stdin
	} else {
		var err error
		file, err = os.OpenFile(params.File(), os.O_RDONLY, os.ModePerm)
		if err != nil {
			log.Fatal("Can't open file %s for read", params.File())
		}
		defer file.Close()
	}
	reader := bufio.NewReader(file)

	if params.destination != "" {
		if err := os.Chdir(params.destination); err != nil {
			log.Fatal("Main: Can't change dir to: ", params.destination)
		}
	}
	orders := orders.ReadOrders(params.destination)

	var splitter splitter.Splitter
	if (params.dumptype == DUMPTYPE_PGSQL) {
		log.Trace("Creatinc pg dmpsplitter")
		splitter, err = pg_dumpsplit.NewSplitter(orders, params.chunkSize, log.Default())
	} else if (params.dumptype == DUMPTYPE_MYSQL) {
		log.Trace("Creatinc mysql dmpsplitter")
		splitter, err = mysql_dumpsplit.NewSplitter(orders, params.chunkSize, log.Default())
	} else {
		log.Fatal("Unsupported dumptype: ", params.dumptype)
	}
	if err != nil {
		log.Fatal("Can't initialize datasplitter type: %s, error: %s", params.dumptype, err.Error())
	}
	defer splitter.Close()

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			log.Debug("EOF")
			break
		} else if err != nil {
			log.Fatal("Can't read iput file: %s", err.Error())
		}
		line = strings.TrimRight(line, "\n\r")
		err = splitter.HandleLine(line)
		if err != nil {
			log.Fatal("Can't handle line: [%s] error: %s", line, err.Error())
		}
	}
	splitter.Flush()
	if err := splitter.Error(); err != nil {
		splitter.Close()
		log.Fatal("Error reading input file: ", err)
	}

	if !orders.IsEmpty() {
		err := orders.WriteOrders()
		if err != nil {
			log.Fatal("Error writing orders: %s", err.Error())
		}
	}
}
