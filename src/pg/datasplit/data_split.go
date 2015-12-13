package datasplit

import (
	"writer"
	"io/ioutil"
	"os"
	"fmt"
	"log"
	"bufio"
	"io"
	"orders"
	"sort"
)

type DataSplitter interface {
	FlushData(dumper writer.Writer) error
	AddLine(line string) error
}

func NewDataSplitter(chunk_size int, copy_line string, table orders.Table, ) DataSplitter {
	log.Println("Start dumping data of table: ", table, " columns: ", copy_line)
	return &dataSplitter{
		chunkSize: int64(chunk_size),
		copyLine: copy_line,
		table: table,
	}
}

type dataSplitter struct {
	chunkSize   int64
	copyLine    string
	table       orders.Table
	buffer      buffer
	currentSize int
	tempFiles   []string
}

func (i *dataSplitter) AddLine(line string) error { // data_split.DataSplitter interface
	i.buffer = append(i.buffer, line, )
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
	for _, row := range i.tempFiles {
		file, err := os.Open(row)
		if err != nil {
			return err
		}
		defer file.Close()

		reader := bufio.NewReader(file)
		for {
			if writer.DataSize() > i.chunkSize * 1024 {
				writer.AddLines("\\.")
				writer.ResetOutput(i.table.FileName(part) + ".sql")
				writer.AddLines(i.copyLine)
				part++
			}

			line, err := reader.ReadString('\n')
			if err == io.EOF {
				break
			} else if err != nil {
				return err
			}

			writer.AddLines(line)
		}
	}
	i.tempFiles = nil
	if i.buffer != nil {
		for _, line := range i.buffer {
			if writer.DataSize() > i.chunkSize * 1024 {
				writer.AddLines("\\.")
				writer.ResetOutput(i.table.FileName(part) + ".sql")
				writer.AddLines(i.copyLine)
				part++
			}

			writer.AddLines(line)
		}
		i.buffer = nil
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
	*tempFiles = append(*tempFiles, tempfile.Name(), )
	*buffer = nil
	*currentBuffer = 0
	return nil
}
