package dumpsplit

import (
	"github.com/andriyg76/dbtricks/orders"
	"github.com/andriyg76/dbtricks/writer"
	"github.com/andriyg76/dbtricks/mysql/datasplit"
	"regexp"
	"github.com/andriyg76/dbtricks/splitter"
	"github.com/andriyg76/glogger"
)

var table_structure_re, insert_into_re *regexp.Regexp

func init() {
	table_structure_re = regexp.MustCompile("^-- Table structure for table `(?P<table>.*?)`")
	insert_into_re = regexp.MustCompile(
		"^(?P<insert_into>INSERT INTO .* VALUES) \\((?P<data>.*?)\\);$")
}

func match_table_data(line string) (match bool, insert_into, data string) {
	if !insert_into_re.MatchString(line) {
		return false, "", ""
	}

	match = true
	match_values := insert_into_re.FindStringSubmatch(line)
	insert_into = match_values[1]
	data = match_values[2]
	return
}

func match_table_structure(line string) (match bool, table string) {
	if !table_structure_re.MatchString(line) {
		match = false
		table = ""
		return
	}

	match = true
	match_values := table_structure_re.FindStringSubmatch(line)
	table = match_values[1]
	return
}

type mysqlSplitter struct {
	counter      int
	dumper       writer.Writer
	table        orders.Table
	data_handler datasplit.DataSplitter
	chunk_size   int
	orders       orders.Orders
	err          error
	logger       glogger.Logger
}

func NewSplitter(orders orders.Orders, chunk_size int, logger glogger.Logger) (splitter.Splitter, error) {
	dumper, err := writer.NewWriter("0000_prologue.sql", logger)
	if err != nil {
		return nil, err
	}
	return &mysqlSplitter{
		counter:    0,
		dumper:     dumper,
		table:      nil,
		chunk_size: chunk_size,
		orders:     orders,
		logger:     logger,
	}, nil
}

func (i *mysqlSplitter) HandleLine(line string) error {
	if i.err != nil {
		i.logger.Trace("Returning error: %s", i.err)
		return i.err
	}

	if match, table_name := match_table_structure(line); match {
		i.table = i.orders.GetTable(table_name)
		backup := i.dumper.PopLastLines(2)
		if err := i.dumper.Flush(); err != nil {
			i.err = err
			return err
		}
		i.dumper.ResetOutput(i.table.FileName(0) + ".sql")
		i.dumper.AddLines(backup...)
		i.dumper.AddLines(line)
	} else if match, insert_into, data := match_table_data(line); match {
		if i.data_handler == nil {
			start_line := insert_into
			i.data_handler = datasplit.NewDataSplitter(i.chunk_size, start_line, i.table, i.logger)
		}
		i.data_handler.AddLine(data)
	} else if i.data_handler != nil && line == "\n" {
		i.logger.Trace("Skipping empty line")
	} else if i.data_handler != nil {
		i.logger.Trace("Closing data handler")
		if err := i.data_handler.FlushData(i.dumper); err != nil {
			i.err = err
			return err
		}
		i.data_handler = nil
	} else {
		i.logger.Trace("Append line to dumper: %s", line)
		i.dumper.AddLines(line)
	}
	return nil
}

func (i *mysqlSplitter) Flush() error {
	return nil
}

func (i *mysqlSplitter) Close() {

}

func (i *mysqlSplitter) Error() error {
	return i.err
}
