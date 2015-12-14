package mergesort

import (
	"errors"
	"bufio"
	"strings"
	"io"
)

type Reader interface {
	ReadLine() (error, string)
}

// Composition of reader and close, os.File have to match
type DisposableIoReader interface {
	io.Reader
	io.Closer
}

type DisposableReader interface {
	Reader
	Close()
}

type stringAndErr struct {
	string string
	error  error
}

type fileReader struct {
	file    DisposableIoReader
	channel chan stringAndErr
}

func NewAsyncFileReader(file DisposableIoReader) (error, DisposableReader) {
	if file == nil {
		return errors.New("null pointer exception: file"), nil
	}

	fileRrd := bufio.NewReader(file)

	reader := &fileReader{
		file: file,
		channel: make(chan stringAndErr),
	}

	go func() {
		for {
			line, err := fileRrd.ReadString('\n')
			if err == io.EOF {
				reader.Close()
				break
			}
			if err != nil {
				line = ""
			} else {
				line = strings.TrimRight(line, "\n\r")
			}
			reader.channel <- stringAndErr{
				string: line,
				error: err,
			}
		}
	}();

	return nil, reader
}

func (i *fileReader) Close() {
	if i.file != nil {
		i.file.Close()
	}

	close(i.channel)
}

func (i *fileReader) ReadLine() (error, string) {
	lae, ok := <-i.channel
	if !ok {
		return io.EOF, ""
	}
	return lae.error, lae.string
}