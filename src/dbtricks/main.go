package main

import (
	"os"
	"bufio"
	_ "fmt"
)


type DataHandler interface {

}

type Dumper interface {

}

func NewDumper(output_file string) (d Dumper, e error) {
	return
}

func main() {
	panic("Panic")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		_, _ = NewDumper("0000_prologue.sql")

		var _ string
		var _ bool
		var _ DataHandler
	}
	if scanner.Err() != nil {
		panic(scanner.Err())
	}
	return
}
