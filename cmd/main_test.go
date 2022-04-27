package cmd

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/xAutoBot/iqoption-sdk/src/entities/messages"
	"github.com/xAutoBot/iqoption-sdk/src/entities/messages/responseMessage"
	"github.com/xAutoBot/iqoption-sdk/src/repositories/iqoptionRepository"
)

func TestGetExpirationTime(t *testing.T) {
	testes := []struct {
		timeSyn  int64
		duration int
		want     int64
	}{
		{timeSyn: 1650670200, duration: 5, want: 1650670500},
		{timeSyn: 1650661200, duration: 2, want: 1650661320},
	}

	for _, test := range testes {
		if GetExpirationTime(test.timeSyn, test.duration) != test.want {
			t.Errorf("GetExpirationTime(%v, %v) is not equal a expected value  %v", test.timeSyn, test.duration, test.want)
		}
	}
}

func TestGetOptionTypeID(t *testing.T) {
	tests := []struct {
		duration int
		want     int
	}{
		{duration: 1, want: 3},
		{duration: 2, want: 3},
		{duration: 5, want: 3},
		{duration: 10, want: 1},
	}

	for _, test := range tests {
		if GetOptionTypeID(test.duration) != test.want {
			t.Errorf("GetOptionTypeID(%v) is not equal a expected value  %v", test.duration, test.want)
		}
	}

}

func TestGetPriceNow(t *testing.T) {

	websocketConnection, err := iqoptionRepository.Connect()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer websocketConnection.Close()

	receivedMessage := make(chan string)
	messageToSend := make(chan string)
	go iqoptionRepository.OnMessage(websocketConnection, receivedMessage)
	go iqoptionRepository.SendMessage(websocketConnection, messageToSend)

	responseGeneratedCandle := make(chan string)
	responsePrice := make(chan float64)
	responseError := make(chan error)

	for {
		select {
		case getPriceNowError := <-responseError:
			if getPriceNowError != nil {
				t.Errorf("GetPriceNow() apresent error")
			}
			return
		case price := <-responsePrice:
			log.Printf("response price is %v", price)
			return

		case receivedMessageJson := <-receivedMessage:
			var receivedMessage messages.Message
			json.Unmarshal([]byte(receivedMessageJson), &receivedMessage)
			switch receivedMessage.Name {
			case "authenticated":
				activeID := 76
				responsePrice, responseError = GetPriceNow(messageToSend, responseGeneratedCandle, activeID)

			case "candle-generated":
				responseGeneratedCandle <- receivedMessageJson
			}
		}
	}
}
func TestGetBalances(t *testing.T) {
	websocketConnection, err := iqoptionRepository.Connect()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer websocketConnection.Close()

	receivedMessage := make(chan string)
	messageToSend := make(chan string)
	go iqoptionRepository.OnMessage(websocketConnection, receivedMessage)
	go iqoptionRepository.SendMessage(websocketConnection, messageToSend)

	responseGetBalances := make(chan string)
	responseBalances := make(chan []responseMessage.BalancesMsg)
	getBalancesResponseError := make(chan error)

	for {
		select {
		case getBalancesResponseError := <-getBalancesResponseError:
			t.Errorf("getBalances() apresent error => %v", getBalancesResponseError)
		case <-responseBalances:
			return

		case receivedMessageJson := <-receivedMessage:
			var receivedMessage messages.Message
			json.Unmarshal([]byte(receivedMessageJson), &receivedMessage)
			switch receivedMessage.Name {
			case "authenticated":
				var authenticatedMessage responseMessage.Authenticated
				json.Unmarshal([]byte(receivedMessageJson), &authenticatedMessage)
				if authenticatedMessage.Msg == true {
					responseBalances, getBalancesResponseError = getBalances(messageToSend, responseGetBalances)
				}

			case "balances":
				responseGetBalances <- receivedMessageJson
			}
		}
	}
}

func TestOpenOrder(t *testing.T) {
	websocketConnection, err := iqoptionRepository.Connect()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer websocketConnection.Close()

	receivedMessage := make(chan string)
	messageToSend := make(chan string)
	go iqoptionRepository.OnMessage(websocketConnection, receivedMessage)
	go iqoptionRepository.SendMessage(websocketConnection, messageToSend)

	for {
		select {
		case receivedMessageJson := <-receivedMessage:
			log.Println(receivedMessageJson)
		}
	}
}
