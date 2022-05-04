package authenticationRepository

import (
	"log"
	"testing"
)

func TestGetSSID(t *testing.T) {

	_, errChan := GetSSID()
	err := <-errChan
	if err != nil {
		log.Fatalf(err.Error())
	}
}
