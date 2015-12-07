package datasplit
import "writer"

type DataSplitter interface {
	FlushData(dumper writer.Writer) error
	AddLine(line string) error
}

func NewDataSplitter(chunk_size int, copy_line string, table_name string, table_order int) DataSplitter {
	return nil
}
