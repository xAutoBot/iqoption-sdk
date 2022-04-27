package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/xAutoBot/iqoption-sdk/src/configs"
	"github.com/xAutoBot/iqoption-sdk/src/entities/messages"
	"github.com/xAutoBot/iqoption-sdk/src/entities/messages/responseMessage"
	"github.com/xAutoBot/iqoption-sdk/src/entities/profile"
	"github.com/xAutoBot/iqoption-sdk/src/repositories/iqoptionRepository"
)

const (
	optionTypeBinary int = 1
	optionTypeTurbo  int = 3
)

func GetOptionTypeID(duration int) int {
	if duration <= 5 {
		return optionTypeTurbo
	}
	return optionTypeBinary
}

func OpenOrder(messageToSend, responseOpenOrder chan string, activePriceNow float64, user profile.User, timeSync *int64) (orderID chan string, responseError chan error) {
	go func() {
		activeID := 1
		var duration int = 5
		investiment := 2.0
		direction := "call"
		activePayoutNow := 86

		activePrice, _ := strconv.Atoi(strings.Replace(fmt.Sprintf("%f", activePriceNow), ".", "", -1))

		binaryOptionsOpenOptionBody := messages.BinaryOptionsOpenOptionBody{
			UserBalanceID: user.BalanceId,
			ActiveID:      activeID,
			OptionTypeID:  GetOptionTypeID(duration),
			Direction:     direction,
			Expired:       GetExpirationTime(*timeSync, duration),
			RefundValue:   0,
			Price:         investiment,
			Value:         activePrice,
			ProfitPercent: activePayoutNow,
		}
		binaryOptionsOpenOption := messages.NewSendMessageBinaryOptionsOpenOption(binaryOptionsOpenOptionBody)
		binaryOptionsOpenOptionJson, _ := binaryOptionsOpenOption.Json()
		messageToSend <- string(binaryOptionsOpenOptionJson)
	}()

	orderID = make(chan string)
	responseError = make(chan error)

	go func() {
		for responseOpenOrderJson := range responseOpenOrder {
			log.Println(responseOpenOrderJson)
			responseError <- nil
			orderID <- "evecimar"
			return
		}
	}()
	return
}
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
