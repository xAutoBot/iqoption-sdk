package iqoptionRepository

import (
	"errors"
	"testing"
	"time"

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
		activeId uint8
		duration uint8
		timeStub string
		want     int64
	}{
		{activeId: 1, duration: 1, timeStub: "2022/04/29 21:29:50", want: 1231231},
	}

	for _, teste := range testes {
		connection.time, _ = time.Parse("2006/01/02 15:4:5", teste.timeStub)
		_, err := connection.GetExpirationTime(int(teste.duration))
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

	go connection.OpenOrder(activeId, duration, investiment, direction)
	go connection.OpenOrder(activeId, duration, investiment, direction)
	go connection.OpenOrder(activeId, duration, investiment, direction)
	go connection.OpenOrder(activeId, duration, investiment, direction)
	go connection.OpenOrder(activeId, duration, investiment, direction)
	_, err := connection.OpenOrder(activeId, duration, investiment, direction)

	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestGetDigitalInstrumentID(t *testing.T) {
	testes := []struct {
		activeId  uint8
		duration  uint8
		direction string
		timeStub  string
		want      string
	}{
		{activeId: 1, duration: 1, direction: "call", timeStub: "2022/04/29 20:31:25", want: "do1A20220429D203200T1MCSPT"},
		{activeId: 1, duration: 1, direction: "call", timeStub: "2022/04/29 20:31:44", want: "do1A20220429D203300T1MCSPT"},
		{activeId: 1, duration: 1, direction: "put", timeStub: "2022/04/29 20:31:44", want: "do1A20220429D203300T1MPSPT"},
		{activeId: 1, duration: 5, direction: "call", timeStub: "2022/04/29 20:31:44", want: "do1A20220429D203500T5MCSPT"},
		{activeId: 1, duration: 5, direction: "put", timeStub: "2022/04/29 20:31:44", want: "do1A20220429D203500T5MPSPT"},
		{activeId: 1, duration: 15, direction: "put", timeStub: "2022/04/29 20:31:44", want: "do1A20220429D204500T15MPSPT"},
		{activeId: 1, duration: 15, direction: "call", timeStub: "2022/04/29 20:00:44", want: "do1A20220429D201500T15MCSPT"},
		{activeId: 1, duration: 5, direction: "call", timeStub: "2022/04/29 23:57:44", want: "do1A20220430D000000T5MCSPT"},
		{activeId: 1, duration: 5, direction: "call", timeStub: "2022/04/29 23:57:44", want: "do1A20220430D000000T5MCSPT"},
		{activeId: 1, duration: 5, direction: "call", timeStub: "2022/12/31 23:57:44", want: "do1A20230101D000000T5MCSPT"},
		{activeId: 1, duration: 15, direction: "call", timeStub: "2022/12/31 23:57:44", want: "do1A20230101D000000T15MCSPT"},
	}

	for _, teste := range testes {
		connection.time, _ = time.Parse("2006/01/02 15:4:5", teste.timeStub)
		expiration, err := connection.GetDigitalInstrumentID(teste.activeId, teste.duration, teste.direction)
		if err != nil {
			t.Errorf(err.Error())
		}
		if expiration != teste.want {
			t.Errorf("duration is %d. I want %s but response is %s", teste.duration, teste.want, expiration)
		}
	}
}
