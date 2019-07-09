package main

import (
	"bufio"
	"fmt"
	log "github.com/andriyg76/glogger"
	mysqldumpsplit "github.com/andriyg76/godbtricks/mysql/dumpsplit"
	"github.com/andriyg76/godbtricks/orders"
	"github.com/andriyg76/godbtricks/params"
	pgdumpsplit "github.com/andriyg76/godbtricks/pg/dumpsplit"
	"github.com/andriyg76/godbtricks/splitter"
	"io"
	"os"
	"strings"
)

//noinspection GoNilness
func main() {
	err, args := params.ParseParams(os.Args)
	if err == params.OK {
		os.Exit(0)
	} else if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error parsing params ", err)
		os.Exit(2)
	}

	if args.Trace {
		log.SetLevel(log.TRACE)
	} else if args.Verbose {
		log.SetLevel(log.DEBUG)
	}
	log.Debug("Args %s", args)

	var file *os.File
	if args.File() == "" || args.File() == "-" {
		log.Trace("Reading stdin")
		file = os.Stdin
	} else {
		var err error
		file, err = os.OpenFile(args.File(), os.O_RDONLY, os.ModePerm)
		if err != nil {
			log.Fatal("Can't open file %s for read", args.File())
			return
		}
		defer func() { _ = file.Close() }()
	}
	reader := bufio.NewReader(file)

	if args.Destination != "" {
		if err := os.Chdir(args.Destination); err != nil {
			log.Fatal("Main: Can't change dir to: ", args.Destination)
			return
		}
	}
	ordersStructure := orders.ReadOrders(args.Destination)

	var dumpSplitter splitter.Splitter
	if args.DumpType == params.Pgsql {
		log.Trace("Creating pg dumpsplitter")
		dumpSplitter, err = pgdumpsplit.NewSplitter(ordersStructure, args.ChunkSize, log.Default())
	} else if args.DumpType == params.Mysql {
		log.Trace("Creating mysql dumpsplitter")
		dumpSplitter, err = mysqldumpsplit.NewSplitter(ordersStructure, args.ChunkSize, log.Default())
	} else {
		log.Fatal("Unsupported dumptype: ", args.DumpType)
		return
	}
	if err != nil {
		log.Fatal("Can't initialize dumpsplitter type: %s, error: %s", args.DumpType, err.Error())
		return
	}
	defer dumpSplitter.Close()

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			log.Debug("EOF")
			break
		} else if err != nil {
			log.Fatal("Can't read iput file: %s", err.Error())
			return
		}
		line = strings.TrimRight(line, "\n\r")
		err = dumpSplitter.HandleLine(line)
		if err != nil {
			log.Fatal("Can't handle line: [%s] error: %s", line, err.Error())
			return
		}
	}
	if err := dumpSplitter.Flush(); err != nil {
		log.Fatal("Error flushing datasplitter parts: ", err)
		return
	}
	if err := dumpSplitter.Error(); err != nil {
		dumpSplitter.Close()
		log.Fatal("Error reading input file: ", err)
		return
	}

	if !ordersStructure.IsEmpty() {
		err := ordersStructure.WriteOrders()
		if err != nil {
			log.Fatal("Error writing orders: %s", err.Error())
			return
		}
	}
}
