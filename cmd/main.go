package cmd

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"

	"github.com/xAutoBot/iqoption-sdk/src/configs"
	"github.com/xAutoBot/iqoption-sdk/src/entities/messages"
	"github.com/xAutoBot/iqoption-sdk/src/entities/messages/responseMessage"
	"github.com/xAutoBot/iqoption-sdk/src/entities/profile"
	"github.com/xAutoBot/iqoption-sdk/src/repositories/iqoptionRepository"
)

func Start() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	websocketConnection, err := iqoptionRepository.Connect()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer websocketConnection.Close()

	receivedMessage := make(chan string)
	messageToSend := make(chan string)
	responseGeneratedCandle := make(chan string)
	responseBalances := make(chan []responseMessage.BalancesMsg)
	responseGetBalances := make(chan string)
	getBalancesResponseError := make(chan error)
	//getPriceNowResponseError := make(chan error)

	var timeSync int64
	println(timeSync)
	profile := profile.User{}

	go iqoptionRepository.OnMessage(websocketConnection, receivedMessage)
	go iqoptionRepository.SendMessage(websocketConnection, messageToSend)

	for {
		select {
		case responseBalances := <-responseBalances:
			profile.Balances = responseBalances
			profile.ChangeBalance(configs.GetAccountType())

		case <-getBalancesResponseError:
			responseBalances, getBalancesResponseError = getBalances(messageToSend, responseGetBalances)

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
			case "timeSync":
				var responseTimeSync responseMessage.TimeSync
				json.Unmarshal([]byte(receivedMessageJson), &responseTimeSync)
				timeSync = responseTimeSync.Msg
			case "candle-generated":
				responseGeneratedCandle <- receivedMessageJson
			case "balances":
				responseGetBalances <- receivedMessageJson

			default:
				log.Printf(receivedMessageJson)
			}

		case <-interrupt:
			err = iqoptionRepository.CloseConnection(websocketConnection)
			if err != nil {
				log.Println("Error on close connection:", err)
				return
			}
			return
		}
	}
}
