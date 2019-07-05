package orders

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/andriyg76/dbtricks/utils"
	"io/ioutil"
	"os"
	"sort"
	"strings"
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
			utils.IntInBase(i.tableOrder, 36, 4),
			strings.Replace(i.tableName, ".", "_", 0))
	} else {
		return fmt.Sprintf("%v_%v_%v",
			utils.IntInBase(i.tableOrder, 36, 4),
			strings.Replace(i.tableName, ".", "_", 0),
			utils.IntInBase(part, 36, 6),
		)
	}
}

type orders struct {
	orders    map[string]table
	targetDir string
}

const tablesIncrement int = 36 * 8

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
			newTable := table{
				tableName: tableName,
				tableOrder: (last + i.orders[k].tableOrder) / 2,
			}
			i.orders[tableName] = newTable
			return newTable
		}
		last = i.orders[k].tableOrder
	}
	newTable := table{
		tableName: tableName,
		tableOrder: last + tablesIncrement,
	}
	i.orders[tableName] = newTable
	return newTable
}

func (i *orders) GetSchemeTableOrder(scheme string, table string) int {
	return i.GetTableOrder(scheme + "." + table)
}

func (i *orders) getMap() map[string]table {
	return i.orders
}

func emptyOrders(targetDir string) *orders {
	return &orders{
		orders:    map[string]table{},
		targetDir: targetDir,
	}
}

const ordersFileName = ".orders"

func ReadOrders(targetDir string) Orders {
	jsontext, err := ioutil.ReadFile(ordersFileName)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintln(os.Stderr, "Can't read ", ordersFileName, " file, will use empty orders")
			return emptyOrders(targetDir)
		} else {
			panic("Can't read " + ordersFileName + " " + err.Error())
		}
	}

	return readOrders(jsontext, targetDir)
}

func readOrders(jsontext []byte, targetDir string) Orders {
	_orders := emptyOrders(targetDir)

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

var emptyJson = []byte("{}")

func (i *orders) writeOrders() []byte {
	orders := map[string]int{}
	for _, v := range i.orders {
		orders[v.tableName] = v.tableOrder
	}

	jsonstring, err := json.Marshal(orders)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Can't serialize json: ", i.orders, " :", err.Error())
		return emptyJson
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

	err := ioutil.WriteFile(ordersFileName, jsonstring, os.ModePerm)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Can't write ", ordersFileName, " :", err.Error())
		return err
	}
	return nil
}

func (i *orders) IsEmpty() bool {
	return len(i.orders) == 0
}
