package authenticationRepository

import (
	"log"
	"testing"
)

func TestGetSSID(t *testing.T) {

	_, err := GetSSID()

	if err != nil {
		log.Fatalf(err.Error())
	}
}
