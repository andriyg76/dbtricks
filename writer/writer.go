package writer

import (
	"os"
	"errors"
	"fmt"
	"github.com/andriyg76/glogger"
)

type Writer interface {
	AddLines(line... string)
	PopLastLine() []string
	PopLastLines(count int) []string
	ResetOutput(output_file string) error
	Flush() error
	Close()
	DataSize() int64
}

func NewWriter(output_file string, logger glogger.Logger) (Writer, error) {
	dumper := &dumper{
		lines: []string{},
		logger: logger,
	}
	err := dumper.setOutput(output_file)
	if err != nil {
		return nil, err
	}
	return dumper, nil
}

type dumper struct {
	lines       []string
	logger      glogger.Logger
	output_file *os.File
	err         error
}

func (i *dumper) PopLastLine() []string {
	return i.PopLastLines(1)
}

func (i *dumper) PopLastLines(count int) []string {
	if i.lines != nil {
		if count <= 0 {
			count = 1
		}
		if count > len(i.lines) {
			count = len(i.lines)
		}
		from := len(i.lines) - count
		ret := i.lines[from:]
		i.lines = i.lines[:from]
		return ret
	}
	return nil
}

func (i* dumper) setOutput(file_name string) error {
	if i.err != nil {
		return i.err
	}

	if i.output_file != nil {
		return errors.New("setOutput: Dumper have already output_file defined")
	}

	i.output_file, i.err = os.OpenFile(file_name, os.O_CREATE | os.O_TRUNC | os.O_RDWR, os.ModePerm)
	if i.err != nil {
		return i.err
	}

	return nil
}

func (i* dumper) Flush() error {
	if i.err != nil {
		return i.err
	}

	if i.lines != nil {
		if i.output_file == nil {
			return errors.New("Flush: Dumper output file is not set")
		}
		for _, line := range i.lines {
			_, i.err = fmt.Fprintln(i.output_file, line)
			if i.err != nil {
				return i.err
			}
		}
		i.lines = nil
	}
	return i.err
}

func (i* dumper) Close() {
	if i.output_file != nil {
		i.output_file.Close()
		i.output_file = nil
	}
	i.lines = nil
	i.err = nil
}

func (i* dumper) AddLines(lines... string)  {
	i.lines = append(i.lines, lines...)
}

func (i* dumper) ResetOutput(output_file string) error {
	if i.err != nil {
		return i.err
	}
	if err := i.Flush(); err != nil {
		return i.err
	}
	i.Close()
	i.logger.Debug("Resetting output to %s", output_file)
	if err := i.setOutput(output_file); err != nil {
		return err
	}
	return nil
}

func (i* dumper) DataSize() int64 {
	if i.err != nil {
		return -1
	}
	if i.lines == nil {
		return 0
	}
	var result int64
	for _, line := range i.lines {
		result += int64(len(line) + 1)
	}
	return result
}