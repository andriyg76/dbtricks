package datasplit

import (
	"fmt"
	"github.com/andriyg76/glogger"
	"github.com/andriyg76/godbtricks/orders"
	"github.com/andriyg76/godbtricks/writer"
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

func NewDataSplitter(chunkSize int, insertLine string, table orders.Table, logger glogger.Logger) DataSplitter {
	logger.Debug("Start dumping data of table: %s columns: %s ", table.TableName(), insertLine)
	return &dataSplitter{
		chunkSize:  int64(chunkSize),
		insertLine: insertLine,
		table:      table,
		logger:     logger,
	}
}

type dataSplitter struct {
	chunkSize   int64
	insertLine  string
	table       orders.Table
	buffer      buffer
	currentSize int
	logger      glogger.Logger
	tempFiles   []string
}

func (i *dataSplitter) AddLine(line string) error { // data_split.DataSplitter interface
	i.buffer = append(i.buffer, line)
	i.currentSize += len(line) + 1

	if i.currentSize > i.currentSize*1024 {
		sort.Sort(i.buffer)
		return flushBufferToTemp(&i.tempFiles, &i.buffer, &i.currentSize)
	}
	return nil
}

// data_split.DataSplitter interface
func (i *dataSplitter) FlushData(writer writer.Writer) error {
	sort.Sort(i.buffer)
	part := 1

	readers := []mergesort.Reader{mergesort.NewArrayReader(i.buffer)}
	var files []*os.File
	defer func() {
		for _, row := range files {
			_ = row.Close()
		}
	}()
	for _, row := range i.tempFiles {
		file, err := os.Open(row)
		if err != nil {
			return err
		}
		files = append(files, file)
		err, reader := mergesort.NewAsyncFileReader(file, i.logger.TraceLogger())
		if err != nil {
			return err
		}
		readers = append(readers, reader)
	}

	sorted := mergesort.MergeSort(lessByFirstOrNextValue, i.logger.TraceLogger(), readers...)
	endChunk := false
	endBatch := true
	var chunkSize int64 = 0
	var batchSize int64 = 0
	for {
		err, row := sorted.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		if endChunk {
			if err = writer.ResetOutput(i.table.FileName(part) + ".sql"); err != nil {
				return err
			}
			part++
			chunkSize = 0
		}

		if endChunk || endBatch {
			writer.AddLines(i.insertLine)
			chunkSize += int64(len(i.insertLine))
			batchSize = int64(len(i.insertLine))
		}

		endChunk = false
		endBatch = false

		if chunkSize+int64(len(row))+4 >= i.chunkSize {
			endChunk = true
		}

		if batchSize+int64(len(row))+4 >= 5000 {
			endBatch = true
		}

		var delimeter string
		if endChunk || endBatch {
			delimeter = ";"
		} else {
			delimeter = ","
		}
		writer.AddLines(fmt.Sprintf("(%s)%s", row, delimeter))
		batchSize += int64(len(row) + 4)
		chunkSize += int64(len(row) + 4)

		if endBatch || endChunk {
			if err := writer.Flush(); err != nil {
				return err
			}
		}

		if endChunk {
			part++
		}
	}

	if lines := writer.PopLastLine(); len(lines) > 0 {
		line := lines[0][0:len(lines[0])-1] + ";"
		writer.AddLines(line)
	}

	if err := writer.Flush(); err != nil {
		return err
	}

	return nil
}

func flushBufferToTemp(tempFiles *[]string, buffer *buffer, currentBuffer *int) error {
	tempfile, err := ioutil.TempFile(os.TempDir(), "data.XXXX")
	if err != nil {
		return err
	}
	defer func() { _ = tempfile.Close() }()
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
