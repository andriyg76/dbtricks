package dbtricks

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestOrdersRead(t *testing.T)  {
	_orders := readOrders([]byte("/dev/null"))

	assert.Empty(t, _orders.getMap())

	_orders = readOrders([]byte(`{"one": 1}`))

	assert.NotEmpty(t, _orders.getMap())
	assert.Equal(t, len(_orders.getMap()), 1)
}

func TestOrdersFirst(t *testing.T) {
	_orders := emptyOrders()

	order := _orders.GetTableOrder("first")
	assert.Equal(t, order, ORDERS_INCREMENT)
}

func TestBeforeFirst(t *testing.T) {
	_orders := emptyOrders()
	_orders.orders["last"] = 100

	order := _orders.GetTableOrder("first")
	assert.Equal(t, order, int32(50))
}

func TestAppendLast(t *testing.T) {
	_orders := emptyOrders()
	_orders.orders["first"] = 50

	order := _orders.GetTableOrder("last")
	assert.Equal(t, order, 50 + ORDERS_INCREMENT)
}

func TestInsertBetween(t *testing.T) {
	_orders := emptyOrders()
	_orders.orders = map[string]int32 {
		"item3": 100,
		"item1": 50,
	}

	order := _orders.GetTableOrder("item2")
	assert.Equal(t, order, int32(75))
}

func TestWrite(t *testing.T) {
	_orders := emptyOrders()

	jsontext := _orders.writeOrders()

	assert.Equal(t, jsontext, []byte("{}"))

	_orders.GetTableOrder("table2")

	jsontext = _orders.writeOrders()

	assert.Equal(t, jsontext, []byte("{\n\t\"table2\": 288\n}"))

	_orders.GetTableOrder("table1")

	jsontext = _orders.writeOrders()

	assert.Equal(t, jsontext, []byte("{\n\t\"table1\": 144,\n\t\"table2\": 288\n}"))
}