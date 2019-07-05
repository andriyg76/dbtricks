package datasplit

import (
	"fmt"
	"github.com/andriyg76/dbtricks/orders"
	"github.com/andriyg76/dbtricks/writer"
	"github.com/andriyg76/glogger"
	"github.com/andriyg76/mergesort"
	"io"
	"io/ioutil"
	"os"
	"sort"
)

type DataSplitter interface {
	FlushData(dumper writer.Writer) error
	AddLine(line string) error
}

func NewDataSplitter(chunkSize int, columnsList string, table orders.Table, logger glogger.Logger) DataSplitter {
	logger.Debug("Start dumping data of table: ", table, " columns: ", columnsList)
	return &dataSplitter{
		chunkSize: int64(chunkSize),
		copyLine:  columnsList,
		table:     table,
		logger:    logger,
	}
}

type dataSplitter struct {
	chunkSize   int64
	copyLine    string
	table       orders.Table
	buffer      buffer
	currentSize int
	tempFiles   []string
	logger      glogger.Logger
}

func (i *dataSplitter) AddLine(line string) error {
	// data_split.DataSplitter interface
	i.buffer = append(i.buffer, line)
	i.currentSize += len(line) + 1

	if i.currentSize > i.currentSize * 1024 {
		sort.Sort(i.buffer)
		return flushBufferToTemp(&i.tempFiles, &i.buffer, &i.currentSize)
	}
	return nil
}

// data_split.DataSplitter interface
func (i *dataSplitter) FlushData(writer writer.Writer) error {
	sort.Sort(i.buffer)
	part := 1

	writer.AddLines(i.copyLine)
	readers := []mergesort.Reader{mergesort.NewArrayReader(i.buffer)}
	for _, row := range i.tempFiles {
		file, err := os.Open(row)
		if err != nil {
			return err
		}
		err, reader := mergesort.NewAsyncFileReader(file, i.logger.TraceLogger())
		if err != nil {
			return err
		}
		readers = append(readers, reader)
		defer func(i mergesort.DisposableReader){i.Close()} (reader)
	}

	sorted := mergesort.MergeSort(lessByFirstOrNextValue, i.logger.TraceLogger(), readers...)
	for {
		err, row := sorted.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		if writer.DataSize() > i.chunkSize * 1024 {
			writer.AddLines("\\.")
			writer.ResetOutput(i.table.FileName(part) + ".sql")
			writer.AddLines(i.copyLine)
			part++
		}

		writer.AddLines(row)
	}

	writer.AddLines("\\.")

	return nil
}

func flushBufferToTemp(tempFiles *[]string, buffer *buffer, currentBuffer *int) error {
	tempfile, err := ioutil.TempFile(os.TempDir(), "data.XXXX")
	if err != nil {
		return err
	}
	defer tempfile.Close()
	for _, line := range *buffer {
		_, err = fmt.Println(line)
		if err != nil {
			return err
		}
	}
	*tempFiles = append(*tempFiles, tempfile.Name())
	*buffer = nil
	*currentBuffer = 0
	return nil
}
