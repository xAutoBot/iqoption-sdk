package order

import (
	"encoding/json"
	"log"
	"time"
)

type OrderMessageBody struct {
	User_balance_id int     `json:"user_balance_id"`
	Active_id       int     `json:"active_id"`
	Option_type_id  int     `json:"option_type_id"`
	Direction       string  `json:"direction"`
	Expired         int     `json:"expired"`
	Refund_value    int     `json:"refund_value"`
	Price           float32 `json:"price"`
	Value           int     `json:"value"`
	Profit_percent  int     `json:"profit_percent"`
}

func NewOrderMessageBody(userBalanceId int, activeId int, optionTypeId int, direction string, expired int, refundValue int, price float32, value int, profitPercent int) OrderMessageBody {
	return OrderMessageBody{
		User_balance_id: userBalanceId,
		Active_id:       activeId,
		Option_type_id:  optionTypeId,
		Direction:       direction,
		Expired:         expired,
		Refund_value:    refundValue,
		Price:           price,
		Value:           value,
		Profit_percent:  profitPercent,
	}
}

type OrderMessage struct {
	Name    string           `json:"name"`
	Version float32          `json:"version"`
	Body    OrderMessageBody `json:"body"`
}

func NewOrderMessage(orderMessageBody OrderMessageBody, name string, version float32) OrderMessage {
	return OrderMessage{
		Name:    name,
		Version: version,
		Body:    orderMessageBody,
	}
}

type Order struct {
	Name       string       `json:"name"`
	Request_id int          `json:"request_id"`
	Local_time int          `json:"local_time"`
	Msg        OrderMessage `json:"msg"`
}

func NewOrder(orderMessage OrderMessage) Order {
	local_time := int(time.Now().Unix())
	log.Printf(" evecimar")
	return Order{
		Name:       "sendMessage",
		Request_id: local_time,
		Local_time: local_time,
		Msg:        orderMessage,
	}
}

func (o Order) Json() ([]byte, error) {
	log.Print(" evecimar")
	j, err := json.Marshal(o)

	if err != nil {
		log.Fatalf(err.Error())
		return nil, err
	}
	return j, nil
}

func (o Order) getOrderMessage() string {

	return `{"name":"sendMessage","request_id":"161","local_time":53822,"msg":{"name":"binary-options.open-option","version":"1.0","body":{"user_balance_id":21263150,"active_id":2,"option_type_id":3,"direction":"call","expired":1649383200,"refund_value":0,"price":1.0,"value":831285,"profit_percent":85}}}`
}
