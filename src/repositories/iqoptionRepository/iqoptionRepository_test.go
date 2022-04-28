package iqoptionRepository

import (
	"errors"
	"log"
	"testing"

	"github.com/xAutoBot/iqoption-sdk/src/configs"
)

var connection *IqOptionRepository
var err error

func init() {
	iqOptionRepository := IqOptionRepository{}
	connection, err = iqOptionRepository.Connect(configs.GetAccountType())
	if err != nil {
		panic(err)
	}

}

func TestConnect(t *testing.T) {
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestGetBalances(t *testing.T) {

	_, err = connection.GetBalances()
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestGetPriceNow(t *testing.T) {
	activeId := 1
	price, responsePriceErro := connection.GetPriceNow(activeId)
	if responsePriceErro != nil {
		t.Errorf(responsePriceErro.Error())
	}
	if price <= 0 {
		t.Errorf(errors.New("price is zero").Error())
	}
}

func TestGetOptionTypeID(t *testing.T) {
	testes := []struct {
		value int
		want  int
	}{
		{value: 1, want: 3},
		{value: 2, want: 3},
		{value: 3, want: 3},
		{value: 4, want: 3},
		{value: 5, want: 3},
		{value: 6, want: 1},
		{value: 7, want: 1},
		{value: 8, want: 1},
		{value: 9, want: 1},
	}

	for _, teste := range testes {
		typeId := connection.GetOptionTypeID(teste.value)
		if typeId != teste.want {
			t.Errorf(errors.New("Incorrrect responsee").Error())
		}
	}

}

func TestGetExpirationTime(t *testing.T) {
	testes := []struct {
		expiration int
		wantError  error
	}{
		{expiration: 1, wantError: nil},
		{expiration: 4, wantError: nil},
		{expiration: 5, wantError: nil},
	}

	for _, teste := range testes {
		_, err := connection.GetExpirationTime(teste.expiration)
		if err != nil {
			t.Errorf("%v", err.Error())
		}

	}
}

func TestOpenOrder(t *testing.T) {
	activeId := 1
	duration := 5
	investiment := 2.00
	direction := "call"
	// priceNow, _ := connection.GetPriceNow(activeId)
	// go connection.OpenOrder(activeId, duration, investiment, direction)
	// go connection.OpenOrder(activeId, duration, investiment, direction)
	// go connection.OpenOrder(activeId, duration, investiment, direction)
	// go connection.OpenOrder(activeId, duration, investiment, direction)
	// go connection.OpenOrder(activeId, duration, investiment, direction)
	order, err := connection.OpenOrder(activeId, duration, investiment, direction)

	// json := `{"name":"sendMessage","request_id":"4a45529c-5bb4-4b82-8c42-dae70a8244d2","local_time":7177585,"msg":{"body": {"profit_percent":0, "refund_value":0,"price": 1, "active_id": 1, "expired": 1651119300, "direction": "put", "option_type_id": 1, "user_balance_id": 21263150}, "name": "binary-options.open-option", "version": "1.0"}}`
	// json := `{"name":"sendMessage","request_id":"4a45529c-5bb4-4b82-8c42-dae70a8244d2","local_time":7177585,"msg":{"name":"binary-options.open-option","version":"1.0","body":{"user_balance_id":21263150,"active_id":1,"option_type_id":1,"direction":"call","expired":1651119300,"refund_value":0,"price":2,"value":0,"profit_percent":0}}}`

	// err := connection.SendMessage([]byte(json))

	if err != nil {
		t.Errorf(err.Error())
	}
	orderJson, _ := order.Json()
	log.Printf("%s ", orderJson)
}
