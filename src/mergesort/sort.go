package mergesort
import (
	"os"
	"errors"
	"bufio"
	"strings"
	"io"
)

type Reader interface {
	ReadLine() (string, error)
}

type DisposableReader interface {
	Reader
	Close()
}

type fileReader struct {
	file    *os.File
	channel chan struct {string; error}
}

func NewAsyncFileReader(file *os.File) (error, DisposableReader) {
	if file == nil {
		return errors.New("null pointer exception: " + file)
	}

	fileRrd := bufio.NewReader(file)

	reader := &fileReader{
		file: file,
		channel: make(chan struct {string; error}),
	}

	go func() {
		line, err := fileRrd.ReadString('\n')
		if err == io.EOF {
			reader.Close()
		}
		if err != nil {
			reader.channel < struct{nil; err}
		} else {
			line := strings.TrimRight(line, "\n\r")
			reader.channel < struct{line; err}
		}
	}();

	return reader
}

func (i *fileReader) Close() {
	if i.file != nil {
		i.file.Close()
	}
}
