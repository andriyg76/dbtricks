package orders

import (
	"os"
	"io/ioutil"
	"fmt"
	"encoding/json"
	"sort"
	"bytes"
)
type orders struct {
	orders map[string]int
	target_dir string
}

type Orders interface {
	GetTableOrder(table string) int
	GetSchemeTableOrder(scheme string, table string) int
	getMap() map[string]int
	writeOrders() []byte
	WriteOrders() error
	IsEmpty() bool
}

const tables_increment int = 36 * 8

func (i *orders) GetTableOrder(table string) int {
	order, got := i.orders[table]
	if got {
		return order
	}

	keys := []string{}
	for k := range i.orders {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	last := int(0)
	for _, k := range keys {
		if k > table {
			i.orders[table] = (last + i.orders[k]) / 2
			return i.orders[table]
		}
		last = i.orders[k]
	}
	i.orders[table] = last + tables_increment
	return i.orders[table]
}

func (i *orders) GetSchemeTableOrder(scheme string, table string) int {
	return i.GetTableOrder(scheme + "." + table)
}

func (i *orders) getMap() map[string]int {
	return i.orders
}

func emptyOrders(target_dir string) *orders {
	return &orders{
		orders: map[string]int{},
		target_dir: target_dir,
	}
}

const orders_file_name = ".orders"

func ReadOrders(target_dir string) Orders {
	jsontext, err := ioutil.ReadFile(orders_file_name)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintln(os.Stderr, "Can't read ", orders_file_name, " file, will use empty orders")
			return emptyOrders(target_dir)
		} else {
			panic("Can't read " + orders_file_name + " " + err.Error())
		}
	}

	return readOrders(jsontext, target_dir)
}

func readOrders(jsontext []byte, target_dir string) Orders {
	_orders := emptyOrders(target_dir)

	err := json.Unmarshal(jsontext, &_orders.orders)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Can't parse orders json ", string(jsontext), " :", err.Error())
		_orders.orders = map[string]int{}
	}

	return _orders
}

var empty_json = []byte("{}")

func (i *orders) writeOrders() []byte {
	jsonstring, err := json.Marshal(i.orders)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Can't serialize json: ", i.orders, " :", err.Error())
		return empty_json
	}
	var out = bytes.Buffer{}
	err = json.Indent(&out, jsonstring, "", "\t")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Can't reformat json: ", string(jsonstring), " :", err.Error())
		return jsonstring
	}
	return out.Bytes()
}

func (i *orders) WriteOrders() error {
	jsonstring := i.writeOrders()

	err := ioutil.WriteFile(orders_file_name, jsonstring, os.ModePerm)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Can't write ", orders_file_name, " :", err.Error())
		return err
	}
	return nil
}

func (i *orders) IsEmpty() bool  {
	return len(i.orders) == 0
}