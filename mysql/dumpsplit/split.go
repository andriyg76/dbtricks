package dumpsplit

import (
	"github.com/andriyg76/dbtricks/mysql/datasplit"
	"github.com/andriyg76/dbtricks/orders"
	"github.com/andriyg76/dbtricks/splitter"
	"github.com/andriyg76/dbtricks/writer"
	"github.com/andriyg76/glogger"
	"regexp"
)

var tableStructureRe, insertIntoRe *regexp.Regexp

func init() {
	tableStructureRe = regexp.MustCompile("^-- Table structure for table `(?P<table>.*?)`")
	insertIntoRe = regexp.MustCompile(
		"^(?P<insert_into>INSERT INTO .* VALUES) \\((?P<data>.*?)\\);$")
}

func matchTableData(line string) (match bool, insertInto, data string) {
	if !insertIntoRe.MatchString(line) {
		return false, "", ""
	}

	match = true
	matchValues := insertIntoRe.FindStringSubmatch(line)
	insertInto = matchValues[1]
	data = matchValues[2]
	return
}

func matchTableStructure(line string) (match bool, table string) {
	if !tableStructureRe.MatchString(line) {
		match = false
		table = ""
		return
	}

	match = true
	matchValues := tableStructureRe.FindStringSubmatch(line)
	table = matchValues[1]
	return
}

type mysqlSplitter struct {
	counter     int
	dumper      writer.Writer
	table       orders.Table
	dataHandler datasplit.DataSplitter
	chunkSize   int
	orders      orders.Orders
	err         error
	logger      glogger.Logger
}

func NewSplitter(orders orders.Orders, chunkSize int, logger glogger.Logger) (splitter.Splitter, error) {
	dumper, err := writer.NewWriter("0000_prologue.sql", logger)
	if err != nil {
		return nil, err
	}
	return &mysqlSplitter{
		counter:   0,
		dumper:    dumper,
		table:     nil,
		chunkSize: chunkSize,
		orders:    orders,
		logger:    logger,
	}, nil
}

func (i *mysqlSplitter) HandleLine(line string) error {
	if i.err != nil {
		i.logger.Trace("Returning error: %s", i.err)
		return i.err
	}

	if match, tableName := matchTableStructure(line); match {
		i.table = i.orders.GetTable(tableName)
		backup := i.dumper.PopLastLines(2)
		if err := i.dumper.Flush(); err != nil {
			i.err = err
			return err
		}
		i.dumper.ResetOutput(i.table.FileName(0) + ".sql")
		i.dumper.AddLines(backup...)
		i.dumper.AddLines(line)
	} else if match, insertInto, data := matchTableData(line); match {
		if i.dataHandler == nil {
			startLine := insertInto
			i.dataHandler = datasplit.NewDataSplitter(i.chunkSize, startLine, i.table, i.logger)
		}
		i.dataHandler.AddLine(data)
	} else if i.dataHandler != nil && line == "\n" {
		i.logger.Trace("Skipping empty line")
	} else if i.dataHandler != nil {
		i.logger.Trace("Closing data handler")
		if err := i.dataHandler.FlushData(i.dumper); err != nil {
			i.err = err
			return err
		}
		i.dataHandler = nil
	} else {
		i.logger.Trace("Append line to dumper: %s", line)
		i.dumper.AddLines(line)
	}
	return nil
}

func (i *mysqlSplitter) Flush() error {
	if i.err != nil {
		return i.err
	}
	i.err = i.dumper.Flush()
	return i.err
}

func (i *mysqlSplitter) Close() {
	i.dumper.Close()
}

func (i *mysqlSplitter) Error() error {
	return i.err
}
