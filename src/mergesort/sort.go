package mergesort

import (
	"io"
)

type eofReader struct {
}

func (i *eofReader) ReadLine() (error, string) {
	return io.EOF, ""
}

var eof = &eofReader{}

type combinedReaders struct {
	left, right Reader
	one, two    string
	err1, err2  error
}

func (i *combinedReaders) ReadLine() (error, string) {
	if i.err1 != nil && i.err1 != io.EOF {
		return i.err1, ""
	} else if i.err2 != nil && i.err2 != io.EOF {
		return i.err2, ""
	} else if i.err1 == io.EOF && i.err2 == io.EOF {
		return io.EOF, ""
	} else if i.err1 == nil && i.err2 == nil {
		if i.one < i.two {
			i.err1, i.one = i.left.ReadLine()
			return nil, i.one
		} else {
			i.err2, i.two = i.right.ReadLine()
			return nil, i.two
		}
	} else if i.err1 == io.EOF {
		i.err2, i.two = i.right.ReadLine()
		return nil, i.two
	} else if i.err2 == io.EOF {
		i.err1, i.one = i.left.ReadLine()
		return nil, i.one
	}
	panic("unexpected state")
}

func CombineReaders(left, right Reader) Reader {
	i := &combinedReaders{
		left: left,
		right: right,
	}
	i.err1, i.one = i.left.ReadLine()
	i.err2, i.two = i.right.ReadLine()
	return i
}

func MergeSort(readers... Reader) Reader {
	if len(readers) == 0 {
		return eof
	} else if len(readers) == 1 {
		return readers[0]
	}
	middle := len(readers) / 2
	return CombineReaders(MergeSort(readers[:middle]...), MergeSort(readers[middle:]...))
}