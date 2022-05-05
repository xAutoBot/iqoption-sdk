package iqoptionService

import (
	"log"
	"testing"
)

func TestGetAllActiveInfo(t *testing.T) {
	digitalActives, err := connection.GetAllActiveInfo()
	if err != nil {
		t.Errorf(err.Error())
	}
	log.Println(digitalActives)
}
