package iqoptionRepository

import (
	"errors"
	"log"
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
		want     string
	}{
		{activeId: 1, duration: 1, timeStub: "2022/04/29 21:29:10", want: "2022/04/29 21:30:00"},
		{activeId: 1, duration: 1, timeStub: "2022/04/29 21:59:10", want: "2022/04/29 22:00:00"},
		{activeId: 1, duration: 1, timeStub: "2022/04/29 21:29:50", want: "2022/04/29 21:31:00"},
		{activeId: 1, duration: 1, timeStub: "2022/04/29 21:58:50", want: "2022/04/29 22:00:00"},
		{activeId: 1, duration: 1, timeStub: "2022/04/29 23:58:50", want: "2022/04/30 00:00:00"},
		{activeId: 1, duration: 1, timeStub: "2022/04/30 23:58:50", want: "2022/05/01 00:00:00"},
		{activeId: 1, duration: 1, timeStub: "2022/12/31 23:58:50", want: "2023/01/01 00:00:00"},
		{activeId: 1, duration: 2, timeStub: "2022/04/29 21:29:00", want: "2022/04/29 21:31:00"},
		{activeId: 1, duration: 3, timeStub: "2022/04/29 21:29:00", want: "2022/04/29 21:32:00"},
		{activeId: 1, duration: 4, timeStub: "2022/04/29 21:29:00", want: "2022/04/29 21:33:00"},
		{activeId: 2, duration: 5, timeStub: "2022/04/29 21:29:00", want: "2022/04/29 21:34:00"},
		{activeId: 2, duration: 1, timeStub: "2022/04/30 09:35:13", want: "2022/04/30 09:36:00"},
		{activeId: 2, duration: 1, timeStub: "2022/04/30 09:35:50", want: "2022/04/30 09:37:00"},
		{activeId: 2, duration: 2, timeStub: "2022/04/30 09:35:13", want: "2022/04/30 09:37:00"},
		{activeId: 2, duration: 3, timeStub: "2022/04/30 09:35:13", want: "2022/04/30 09:38:00"},
		{activeId: 2, duration: 4, timeStub: "2022/04/30 09:35:13", want: "2022/04/30 09:39:00"},
		{activeId: 2, duration: 5, timeStub: "2022/04/30 09:35:13", want: "2022/04/30 09:40:00"},
		{activeId: 2, duration: 15, timeStub: "2022/04/30 09:35:31", want: "2022/04/30 09:45:00"},
		{activeId: 2, duration: 10, timeStub: "2022/04/30 09:35:31", want: "2022/04/30 09:45:00"},

		{activeId: 2, duration: 1, timeStub: "2022/04/30 17:59:48", want: "2022/04/30 18:01:00"},
		{activeId: 2, duration: 2, timeStub: "2022/04/30 17:59:48", want: "2022/04/30 18:02:00"},
		{activeId: 2, duration: 3, timeStub: "2022/04/30 17:59:48", want: "2022/04/30 18:03:00"},
		{activeId: 2, duration: 4, timeStub: "2022/04/30 17:59:48", want: "2022/04/30 18:04:00"},
		{activeId: 2, duration: 5, timeStub: "2022/04/30 17:59:48", want: "2022/04/30 18:05:00"},
		{activeId: 2, duration: 15, timeStub: "2022/04/30 17:59:48", want: "2022/04/30 18:15:00"},
		{activeId: 2, duration: 30, timeStub: "2022/04/30 17:59:48", want: "2022/04/30 18:30:00"},
		{activeId: 2, duration: 60, timeStub: "2022/04/30 17:59:48", want: "2022/04/30 18:45:00"},

		{activeId: 2, duration: 1, timeStub: "2022/05/01 14:28:04", want: "2022/05/01 14:29:00"},
		{activeId: 2, duration: 2, timeStub: "2022/05/01 14:28:04", want: "2022/05/01 14:30:00"},
		{activeId: 2, duration: 3, timeStub: "2022/05/01 14:28:04", want: "2022/05/01 14:31:00"},
		{activeId: 2, duration: 4, timeStub: "2022/05/01 14:28:04", want: "2022/05/01 14:32:00"},
		{activeId: 2, duration: 5, timeStub: "2022/05/01 14:28:04", want: "2022/05/01 14:33:00"},
		{activeId: 2, duration: 15, timeStub: "2022/05/01 14:28:04", want: "2022/05/01 14:45:00"},
		{activeId: 2, duration: 30, timeStub: "2022/05/01 14:28:04", want: "2022/05/01 15:00:00"},
	}

	for index, teste := range testes {
		connection.time, _ = time.Parse("2006/01/02 15:4:5", teste.timeStub)
		timestamp, _ := connection.GetExpirationTime(int(teste.duration))
		wantInTimeStamp, _ := time.Parse("2006/01/02 15:4:5", teste.want)
		if timestamp != wantInTimeStamp.Unix() {
			t.Errorf("index %v duration is %v. Want %v and received %v", index, teste.duration, teste.want, time.Unix(timestamp, 0).Format("2006/01/02 15:04:05"))
		}

	}
}

func TestOpenOrder(t *testing.T) {
	activeId := 76
	duration := 60
	investiment := 2.00
	direction := "call"

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
		{activeId: 1, duration: 1, direction: "call", timeStub: "2022/04/29 20:59:10", want: "do1A20220429D210000T1MCSPT"},
		{activeId: 1, duration: 1, direction: "put", timeStub: "2022/04/29 20:31:44", want: "do1A20220429D203300T1MPSPT"},
		{activeId: 1, duration: 5, direction: "call", timeStub: "2022/04/29 20:31:44", want: "do1A20220429D203500T5MCSPT"},
		{activeId: 1, duration: 5, direction: "put", timeStub: "2022/04/29 20:31:44", want: "do1A20220429D203500T5MPSPT"},
		{activeId: 1, duration: 15, direction: "put", timeStub: "2022/04/29 20:31:44", want: "do1A20220429D204500T15MPSPT"},
		{activeId: 1, duration: 15, direction: "call", timeStub: "2022/04/29 20:00:44", want: "do1A20220429D201500T15MCSPT"},
		{activeId: 1, duration: 5, direction: "call", timeStub: "2022/04/29 23:57:44", want: "do1A20220430D000000T5MCSPT"},
		{activeId: 1, duration: 5, direction: "call", timeStub: "2022/04/29 23:57:44", want: "do1A20220430D000000T5MCSPT"},
		{activeId: 1, duration: 5, direction: "call", timeStub: "2022/12/31 23:57:44", want: "do1A20230101D000000T5MCSPT"},
		{activeId: 1, duration: 15, direction: "call", timeStub: "2022/12/31 23:57:44", want: "do1A20230101D000000T15MCSPT"},
		{activeId: 1, duration: 15, direction: "call", timeStub: "2022/11/30 23:57:44", want: "do1A20221201D000000T15MCSPT"},
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

func TestOpenDigitalOrder(t *testing.T) {
	activeId := 1
	duration := 15
	investiment := 2.00
	direction := "call"
	instrumentId, _ := connection.GetDigitalInstrumentID(uint8(activeId), uint8(duration), direction)
	instrumentIndex, _ := connection.GetInstrumentIndex(activeId)

	orderID, err := connection.OpenDigitalOrder(instrumentId, activeId, instrumentIndex, investiment, direction)

	if err != nil {
		t.Errorf(err.Error())
	}
	orderStr, _ := orderID.Json()
	t.Logf("%s", orderStr)
}

func TestGetInstrumentIndex(t *testing.T) {

	activeID := 85
	instrumentIndex, _ := connection.GetInstrumentIndex(activeID)
	log.Println(instrumentIndex)
	for {
		time.Sleep(time.Second)
	}
}

func TestGetAllActiveDigitalInfo(t *testing.T) {
	digitalActives, err := connection.GetAllActiveDigitalInfo()
	if err != nil {
		t.Errorf(err.Error())
	}
	log.Println(digitalActives)
}

func TestGetAllActiveBinaryInfo(t *testing.T) {
	digitalActives, err := connection.GetAllActiveBinaryInfo()
	if err != nil {
		t.Errorf(err.Error())
	}
	log.Println(digitalActives)
}

func TestGetAllActiveInfo(t *testing.T) {
	digitalActives, err := connection.GetAllActiveInfo()
	if err != nil {
		t.Errorf(err.Error())
	}
	log.Println(digitalActives)
}
