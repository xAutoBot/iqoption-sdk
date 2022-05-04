package orderService

import "testing"

var orderService OrderService
var err error

func init() {
	orderService, err = NewOrderService()
	if err != nil {
		panic(err.Error())
	}
}

func TestActiveIsOpen(t *testing.T) {
	activeId := 1
	activeIsOpen, err := orderService.ActiveIsOpen(activeId)
	if err != nil {
		t.Error(err.Error())
	}
	t.Logf("%v", activeIsOpen)
}

// func TestOpenBinaryOrder(t *testing.T) {
// 	response, _ := openBinaryOrder()

// 	fmt.Println(response)

// }
