package dbtricks

import (
	"os"
	"io/ioutil"
	"fmt"
	"encoding/json"
	"sort"
	"bytes"
)
type orders struct {
	orders map[string]int32
}

type Orders interface {
	GetTableOrder(table string) int32
	GetSchemeTableOrder(scheme string, table string) int32
	getMap() map[string]int32
	writeOrders() string
}

const ORDERS_INCREMENT int32 = 36 * 8

func (i *orders) GetTableOrder(table string) int32 {
	order, got := i.orders[table]
	if got {
		return order
	}

	keys := []string{}
	for k := range i.orders {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	last := int32(0)
	for _, k := range keys {
		if k > table {
			i.orders[table] = (last + i.orders[k]) / 2
			return i.orders[table]
		}
		last = i.orders[k]
	}
	i.orders[table] = last + ORDERS_INCREMENT
	return i.orders[table]
}

func (i *orders) GetSchemeTableOrder(scheme string, table string) int32 {
	return i.GetTableOrder(scheme + "." + table)
}

func (i *orders) getMap() map[string]int32 {
	return i.orders
}

func emptyOrders() *orders {
	return &orders{
		orders: map[string]int32{},
	}
}

const ORDERS_FILE_NAME = ".orders"

func ReadOrders() Orders {
	jsontext, err := ioutil.ReadFile(ORDERS_FILE_NAME)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Can't read ", ORDERS_FILE_NAME, " file, will use empty orders")
		return emptyOrders()
	}

	return readOrders(jsontext)
}

func readOrders(jsontext []byte) Orders {
	_orders := emptyOrders()

	err := json.Unmarshal(jsontext, &_orders.orders)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Can't parse orders json ", string(jsontext), " :", err.Error())
		_orders.orders = map[string]int32{}
	}

	return _orders
}

func (i *orders) writeOrders() string {
	jsonstring, err := json.Marshal(i.orders)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Can't serialize json: ", i.orders, " :", err.Error())
		return "{}"
	}
	var out = bytes.Buffer{}
	err = json.Indent(&out, jsonstring, "", "\t")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Can't reformat json: ", string(jsonstring), " :", err.Error())
		return string(jsonstring)
	}
	return string(out.Bytes())
}