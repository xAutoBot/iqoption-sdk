package iqoptionRepository

import (
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
	if err != nil {
		t.Errorf(err.Error())
	}
	_, err = connection.GetBalances()
	if err != nil {
		t.Errorf(err.Error())
	}
}
