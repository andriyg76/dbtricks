package main

import (
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

	err := orders.WriteOrders()
	if err != nil {
		panic("Error writing orders: " + err.Error())
	}
}
