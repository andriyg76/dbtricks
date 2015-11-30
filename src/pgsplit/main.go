package main

import (
	"os"
	"bufio"
	_ "fmt"
	"encoding/json"
	"fmt"
	"dbtricks"
)


type DataHandler interface {

}

type Dumper interface {

}

func NewDumper(output_file string) (d Dumper, e error) {
	return
}

func main() {
	orders := dbtricks.ReadOrders()

	fmt.Println("table1", orders.GetTableOrder("table1"))
	fmt.Println("table4", orders.GetTableOrder("table4"))
	fmt.Println("table3", orders.GetTableOrder("table3"))

	panic("Panic")

	order := map[string]int32{
		"one": 1,
	}

	val, err := json.Marshal(order)
	if err != nil {
		panic("Can't marshal map to json " + err.Error())
	}
	fmt.Println(string(val))

	err = json.Unmarshal([]byte(`{"one":2, "two":1}`), &order)
	if err != nil {
		panic("Can't unmarshal json to map" + err.Error())
	}
	fmt.Println(order)


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
