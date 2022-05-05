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

func TestActiveIsOpen(t *testing.T) {
	activeId := 1

	isOpen, err := iqoptionService.ActiveIsOpen(activeId)
	if err != nil {
		t.Errorf(err.Error())
	}
	log.Println(isOpen)
}
