package iqoptionRepository

import (
	"log"
	"testing"
)

func TestConnect(t *testing.T) {
	_, err := Connect()

	if err != nil {
		log.Fatal(err.Error())
	}

}
