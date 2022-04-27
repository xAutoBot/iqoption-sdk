package iqoptionRepository

import (
	"errors"
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
