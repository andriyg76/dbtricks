package dumpsplit

import (
	"github.com/andriyg76/glogger"
	"github.com/andriyg76/godbtricks/orders"
	"github.com/andriyg76/godbtricks/pg/datasplit"
	"github.com/andriyg76/godbtricks/splitter"
	"github.com/andriyg76/godbtricks/writer"
	"regexp"
)

/**
 * Created by andriy on 04/12/15.
 */

var copyRe, dataCommentRe, constraintCommentRe *regexp.Regexp

func init() {
	copyRe = regexp.MustCompile(`^COPY .*? \((.*?)\) FROM stdin;`)
	dataCommentRe = regexp.MustCompile(
		"^-- Data for Name: (?P<table>.*?); Type: TABLE DATA; Schema: (?P<schema>.*?);")
	constraintCommentRe = regexp.MustCompile("^-- Name: .*; Type: (.*CONSTRAINT|INDEX); Schema: ")
}

func matchToCopy(line string) bool {
	return copyRe.MatchString(line)
}

func matchToDataComment(line string) (match bool, table string, schema string) {
	if !dataCommentRe.MatchString(line) {
		match = false
		table = ""
		schema = ""
		return
	}

	match = true
	matchValues := dataCommentRe.FindStringSubmatch(line)
	table = matchValues[1]
	schema = matchValues[2]
	return
}

func matchToConstraintComment(line string) bool {
	return constraintCommentRe.MatchString(line)
}

const eotLine = "\\."

type pgSplitter struct {
	counter     int
	dumper      writer.Writer
	table       orders.Table
	epilogue    bool
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
	return &pgSplitter{
		counter:   0,
		dumper:    dumper,
		table:     nil,
		chunkSize: chunkSize,
		orders:    orders,
		logger:    logger,
	}, nil
}

func (i *pgSplitter) HandleLine(line string) error {
	if i.err != nil {
		return i.err
	}

	if i.dataHandler != nil {
		if line == eotLine {
			i.err = i.dataHandler.FlushData(i.dumper)
			if i.err != nil {
				return i.err
			}
			i.dataHandler = nil
		} else {
			i.err = i.dataHandler.AddLine(line)
			if i.err != nil {
				return i.err
			}
		}
	} else if i.epilogue || line == "" || line == "--" {
		i.dumper.AddLines(line)
	} else if matchToConstraintComment(line) {
		backup := append(i.dumper.PopLastLine(), line)
		if err := i.dumper.ResetOutput("zzzz_epilogue.sql"); err != nil {
			return err
		}
		i.dumper.AddLines(backup...)
		i.epilogue = true
	} else if match, table, schema := matchToDataComment(line); match {
		i.table = i.orders.GetTable(schema + "." + table)

		backup := append(i.dumper.PopLastLine(), line)
		if err := i.dumper.ResetOutput(i.table.FileName(0) + ".sql"); err != nil {
			return err
		}
		i.dumper.AddLines(backup...)
	} else if matchToCopy(line) {
		i.dataHandler = datasplit.NewDataSplitter(i.chunkSize, line, i.table, i.logger)
	} else {
		i.dumper.AddLines(line)
	}
	return nil
}

func (i *pgSplitter) Flush() error {
	if i.err != nil {
		return i.err
	}
	i.err = i.dumper.Flush()
	return i.err
}

func (i *pgSplitter) Close() {
	i.dumper.Close()
}

func (i *pgSplitter) Error() error {
	return i.err
}
