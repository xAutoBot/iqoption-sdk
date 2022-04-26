package cmd

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/evecimar/iqoptionapi/src/configs"
	"github.com/evecimar/iqoptionapi/src/entities/messages"
	"github.com/evecimar/iqoptionapi/src/entities/messages/responseMessage"
	"github.com/evecimar/iqoptionapi/src/entities/profile"
	"github.com/evecimar/iqoptionapi/src/repositories/iqoptionRepository"
)

const (
	optionTypeBinary int = 1
	optionTypeTurbo  int = 3
)

//Return the timestamp of the sum timeSyn with duration
func GetExpirationTime(timeSyc int64, duration int) int64 {
	timeNow := time.Unix(timeSyc, 0).Add(time.Minute * time.Duration(duration))

	return timeNow.Unix()
}

func GetOptionTypeID(duration int) int {
	if duration <= 5 {
		return optionTypeTurbo
	}
	return optionTypeBinary
}

func GetPriceNow(messageToSend, responseGeneratedCandle chan string, activeID int) (responsePrice chan float64, responseError chan error) {

	candleSize := 5
	go func() {
		sendMessageStartCandleGenerate, _ := messages.NewSendMessageStartCandleGenerate(activeID, candleSize).Json()
		messageToSend <- string(sendMessageStartCandleGenerate)
	}()

	responsePrice = make(chan float64)
	responseError = make(chan error)

	go func() {
		for candleJson := range responseGeneratedCandle {
			var responnseCandleGenerated responseMessage.ResponnseCandleGenerated
			json.Unmarshal([]byte(candleJson), &responnseCandleGenerated)
			if responnseCandleGenerated.MicroserviceName == "quotes" && responnseCandleGenerated.Msg.ActiveID == activeID {
				responsePrice <- responnseCandleGenerated.Msg.Close
				responseError <- nil
				sendMessageStopCandleGenerate, _ := messages.NewSendMessageStopCandleGenerate(activeID, candleSize).Json()
				messageToSend <- string(sendMessageStopCandleGenerate)

				return
			}
			responseError <- errors.New("Cant get prince now")
			responsePrice <- 0.00
		}

	}()

	return
}

func getBalances(messageToSend, responseGetBalances chan string) (balances chan []responseMessage.BalancesMsg, responseError chan error) {

	go func() {
		sendMessageGetBalances, _ := messages.NewSendMessageGetBalances().Json()
		messageToSend <- string(sendMessageGetBalances)
	}()

	balances = make(chan []responseMessage.BalancesMsg)
	responseError = make(chan error)

	go func() {
		for receivedMessageJson := range responseGetBalances {
			var balanceMessage responseMessage.Balannces
			json.Unmarshal([]byte(receivedMessageJson), &balanceMessage)
			balances <- balanceMessage.Msg
			responseError <- nil
		}
		balances <- nil
		responseError <- errors.New("Error com get Balance")
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

	// var timeSync int64
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
			// case "timeSync":
			// 	var responseTimeSync responseMessage.TimeSync
			// 	json.Unmarshal([]byte(receivedMessageJson), &responseTimeSync)
			// 	timeSync = responseTimeSync.Msg
			case "candle-generated":
				responseGeneratedCandle <- receivedMessageJson
			case "balances":
				responseGetBalances <- receivedMessageJson
				// balanceID := balanceMessage.Msg[1].ID
				// activeID := 76
				// go func() {
				// 	for timeSync == 0 {
				// 		time.Sleep(time.Second)
				// 	}

				// 	var duration int = 5
				// 	activeID := 76
				// 	investiment := 2.0
				// 	direction := "call"
				// 	activePayoutNow := 80

				// 	activePriceNow, getPriceNowResponseError = GetPriceNow(messageToSend, responseGeneratedCandle, activeID)
				// 	if err != nil {
				// 		log.Fatal(err)
				// 	}
				// 	activePrice, _ := strconv.Atoi(strings.Replace(fmt.Sprintf("%f", <-activePriceNow), ".", "", -1))

				// 	binaryOptionsOpenOptionBody := messages.BinaryOptionsOpenOptionBody{
				// 		UserBalanceID: profile.BalanceId,
				// 		ActiveID:      activeID,
				// 		OptionTypeID:  GetOptionTypeID(duration),
				// 		Direction:     direction,
				// 		Expired:       int(GetExpirationTime(timeSync, duration)),
				// 		RefundValue:   0,
				// 		Price:         investiment,
				// 		Value:         activePrice,
				// 		ProfitPercent: activePayoutNow,
				// 	}
				// 	binaryOptionsOpenOption, _ := messages.NewSendMessageBinaryOptionsOpenOption(binaryOptionsOpenOptionBody).Json()
				// 	messageToSend <- string(binaryOptionsOpenOption)

				// }()

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
