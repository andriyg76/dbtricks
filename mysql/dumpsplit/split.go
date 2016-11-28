package dumpsplit

import (
	"github.com/andriyg76/dbtricks/orders"
	"github.com/andriyg76/dbtricks/writer"
	"github.com/andriyg76/dbtricks/mysql/datasplit"
	"regexp"
	"github.com/andriyg76/dbtricks/splitter"
)



var copy_re, data_comment_re, constraint_comment_re *regexp.Regexp

func init() {
	copy_re = regexp.MustCompile(`^COPY .*? \((.*?)\) FROM stdin;`)
	data_comment_re = regexp.MustCompile(
		"^-- Data for Name: (?P<table>.*?); Type: TABLE DATA; Schema: (?P<schema>.*?);")
	constraint_comment_re = regexp.MustCompile("^-- Name: .*; Type: (.*CONSTRAINT|INDEX); Schema: ")
}

func match_to_copy(line string) bool {
	return copy_re.MatchString(line)
}

func match_to_data_comment(line string) (match bool, table string, schema string) {
	if !data_comment_re.MatchString(line) {
		match = false
		table = ""
		schema = ""
		return
	}

	match = true
	match_values := data_comment_re.FindStringSubmatch(line)
	table = match_values[1]
	schema = match_values[2]
	return
}

func match_to_constraint_comment(line string) bool {
	return constraint_comment_re.MatchString(line)
}

const eot_line = "\\."

type mysqlSplitter struct {
	counter      int
	dumper       writer.Writer
	table        orders.Table
	epilogue     bool
	data_handler datasplit.DataSplitter
	chunk_size   int
	orders       orders.Orders
	err          error
}

func NewSplitter(orders orders.Orders, chunk_size int) (splitter.Splitter, error) {
	dumper, err := writer.NewWriter("0000_prologue.sql")
	if err != nil {
		return nil, err
	}
	return &mysqlSplitter{
		counter:    0,
		dumper:     dumper,
		table:      nil,
		chunk_size: chunk_size,
		orders:     orders,
	}, nil
}

func (i *mysqlSplitter) HandleLine(line string) error {
	if i.err != nil {
		return i.err
	}

	if i.data_handler != nil {
		if line == eot_line {
			i.err = i.data_handler.FlushData(i.dumper)
			if i.err != nil {
				return i.err
			}
			i.data_handler = nil
		} else {
			i.err = i.data_handler.AddLine(line)
			if i.err != nil {
				return i.err
			}
		}
	} else if i.epilogue || line == "" || line == "--" {
		i.dumper.AddLines(line)
	} else if match_to_constraint_comment(line) {
		backup := append(i.dumper.PopLastLine(), line)
		i.dumper.ResetOutput("zzzz_epilogue.sql")
		i.dumper.AddLines(backup...)
		i.epilogue = true
	} else if match, table, schema := match_to_data_comment(line); match {
		i.table = i.orders.GetTable(schema + "." + table)

		backup := append(i.dumper.PopLastLine(), line)
		i.dumper.ResetOutput(i.table.FileName(0) + ".sql")
		i.dumper.AddLines(backup...)
	} else if match_to_copy(line) {
		i.data_handler = datasplit.NewDataSplitter(i.chunk_size, line, i.table)
	} else {
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
