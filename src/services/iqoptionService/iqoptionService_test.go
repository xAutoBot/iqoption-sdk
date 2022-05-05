package iqoptionService

import (
	"log"
	"testing"
)

var iqoptionService IqoptionService
var err error

func init() {
	iqoptionService, err = NewIqoptionService()

	if err != nil {
		panic(err.Error())
	}
}
func TestGetAllActiveInfo(t *testing.T) {
	digitalActives, err := iqoptionService.GetAllActiveInfo()
	if err != nil {
		t.Errorf(err.Error())
	}
	log.Println(digitalActives)
}
