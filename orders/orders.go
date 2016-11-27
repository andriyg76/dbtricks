package orders

import (
	"os"
	"io/ioutil"
	"fmt"
	"encoding/json"
	"sort"
	"bytes"
	"strings"
	"github.com/andriyg76/godbtricks/dbtricks"
)

type Table interface {
	TableName() string
	TableOrder() int
	FileName(part int) string
}

type Orders interface {
	GetTableOrder(table string) int
	GetTable(table string) Table
	GetSchemeTableOrder(scheme string, table string) int
	WriteOrders() error
	IsEmpty() bool
	getMap() map[string]table
}

type table struct {
	tableName  string
	tableOrder int
}

func (i table) TableName() string {
	return i.tableName
}

func (i table) TableOrder() int {
	return i.tableOrder
}

func (i table) FileName(part int) string {
	if part == 0 {
		return fmt.Sprintf("%v_%v",
			dbtricks.IntInBase(i.tableOrder, 36, 4),
			strings.Replace(i.tableName, ".", "_", 0))
	} else {
		return fmt.Sprintf("%v_%v_%v",
			dbtricks.IntInBase(i.tableOrder, 36, 4),
			strings.Replace(i.tableName, ".", "_", 0),
			dbtricks.IntInBase(part, 36, 6),
		)
	}
}

type orders struct {
	orders    map[string]table
	targetDir string
}

const tables_increment int = 36 * 8

func (i orders) addTable(tableName string, tableOrder int) orders {
	var _orders orders = i
	_orders.orders[tableName] = table{
		tableName: tableName,
		tableOrder: tableOrder,
	}
	return _orders
}

func (i *orders) GetTableOrder(tableName string) int {
	return i.GetTable(tableName).TableOrder()
}

func (i *orders) GetTable(tableName string) Table  {
	order, got := i.orders[tableName]
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
		if k > tableName {
			new_table := table{
				tableName: tableName,
				tableOrder: (last + i.orders[k].tableOrder) / 2,
			}
			i.orders[tableName] = new_table
			return new_table
		}
		last = i.orders[k].tableOrder
	}
	new_table := table{
		tableName: tableName,
		tableOrder: last + tables_increment,
	}
	i.orders[tableName] = new_table
	return new_table
}

func (i *orders) GetSchemeTableOrder(scheme string, table string) int {
	return i.GetTableOrder(scheme + "." + table)
}

func (i *orders) getMap() map[string]table {
	return i.orders
}

func emptyOrders(target_dir string) *orders {
	return &orders{
		orders: map[string]table{},
		targetDir: target_dir,
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

	orders := map[string]int{}
	err := json.Unmarshal(jsontext, &orders)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Can't parse orders json ", string(jsontext), " :", err.Error())
	}

	for k, v := range orders {
		_orders.orders[k] = table{
			tableName: k,
			tableOrder: v,
		}
	}

	return _orders
}

var empty_json = []byte("{}")

func (i *orders) writeOrders() []byte {
	orders := map[string]int{}
	for _, v := range i.orders {
		orders[v.tableName] = v.tableOrder
	}

	jsonstring, err := json.Marshal(orders)
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

func (i *orders) IsEmpty() bool {
	return len(i.orders) == 0
}
