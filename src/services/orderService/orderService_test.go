package orderService

var orderService OrderService
var err error

func init() {
	orderService, err = NewOrderService()
	if err != nil {
		panic(err.Error())
	}
}
