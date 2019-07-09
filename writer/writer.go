package writer

import (
	"errors"
	"fmt"
	"github.com/andriyg76/glogger"
	"os"
)

type Writer interface {
	AddLines(line ...string)
	PopLastLine() []string
	PopLastLines(count int) []string
	ResetOutput(outputFile string) error
	Flush() error
	Close()
	DataSize() int64
}

func NewWriter(outputFile string, logger glogger.Logger) (Writer, error) {
	dumper := &dumper{
		lines:  []string{},
		logger: logger,
	}
	err := dumper.setOutput(outputFile)
	if err != nil {
		return nil, err
	}
	return dumper, nil
}

type dumper struct {
	lines      []string
	logger     glogger.Logger
	outputFile *os.File
	err        error
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

func (i *dumper) setOutput(fileName string) error {
	if i.err != nil {
		return i.err
	}

	if i.outputFile != nil {
		return errors.New("setOutput: Dumper have already output_file defined")
	}

	i.outputFile, i.err = os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if i.err != nil {
		return i.err
	}

	return nil
}

func (i *dumper) Flush() error {
	if i.err != nil {
		return i.err
	}

	if i.lines != nil {
		if i.outputFile == nil {
			return errors.New("flush: dumper output file is not set")
		}
		for _, line := range i.lines {
			_, i.err = fmt.Fprintln(i.outputFile, line)
			if i.err != nil {
				return i.err
			}
		}
		i.lines = nil
	}
	return i.err
}

func (i *dumper) Close() {
	if i.outputFile != nil {
		_ = i.outputFile.Close()
		i.outputFile = nil
	}
	i.lines = nil
	i.err = nil
}

func (i *dumper) AddLines(lines ...string) {
	i.lines = append(i.lines, lines...)
}

func (i *dumper) ResetOutput(outputFile string) error {
	if i.err != nil {
		return i.err
	}
	if err := i.Flush(); err != nil {
		return i.err
	}
	i.Close()
	i.logger.Debug("Resetting output to %s", outputFile)
	if err := i.setOutput(outputFile); err != nil {
		return err
	}
	return nil
}

func (i *dumper) DataSize() int64 {
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
