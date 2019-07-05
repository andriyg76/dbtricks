package orders

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOrdersRead(t *testing.T)  {
	_orders := readOrders([]byte("/dev/null"), "")

	assert.Empty(t, _orders.getMap())

	_orders = readOrders([]byte(`{"one": 1}`), "")

	assert.NotEmpty(t, _orders.getMap())
	assert.Equal(t, len(_orders.getMap()), 1)
}

func TestOrdersFirst(t *testing.T) {
	_orders := emptyOrders("")

	order := _orders.GetTableOrder("first")
	assert.Equal(t, order, tablesIncrement)
}

func TestBeforeFirst(t *testing.T) {
	_orders := emptyOrders("").addTable("last", 100)

	order := _orders.GetTableOrder("first")
	assert.Equal(t, order, int(50))
}

func TestAppendLast(t *testing.T) {
	_orders := emptyOrders("").addTable("first", 50)

	order := _orders.GetTableOrder("last")
	assert.Equal(t, order, 50 +tablesIncrement)
}

func TestInsertBetween(t *testing.T) {
	_orders := emptyOrders("").
		addTable("item3", 100).
		addTable("item1", 50)

	order := _orders.GetTableOrder("item2")
	assert.Equal(t, order, int(75))
}

func TestFileName(t *testing.T) {
	table := table{tableName:"name", tableOrder:0x10}

	assert.Equal(t, "000g_name", table.FileName(0))

	assert.Equal(t, "000g_name_000001", table.FileName(1))
}

func TestWrite(t *testing.T) {
	_orders := emptyOrders("")

	jsontext := _orders.writeOrders()

	assert.Equal(t, jsontext, []byte("{}"))

	_orders.GetTableOrder("table2")

	jsontext = _orders.writeOrders()

	assert.Equal(t, jsontext, []byte("{\n\t\"table2\": 288\n}"))

	_orders.GetTableOrder("table1")

	jsontext = _orders.writeOrders()

	assert.Equal(t, jsontext, []byte("{\n\t\"table1\": 144,\n\t\"table2\": 288\n}"))
}