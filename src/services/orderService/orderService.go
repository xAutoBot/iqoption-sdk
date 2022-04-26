package orderService

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/xAutoBot/iqoption-sdk/src/entities/messages"
	"github.com/xAutoBot/iqoption-sdk/src/entities/messages/responseMessage"
)

const (
	optionTypeBinary int = 1
	optionTypeTurbo  int = 3
	demoAccountType  int = 4
	realAccountType  int = 1
)

func GetExpirationTime(timeSyc int, duration int) uint64 {
	timeNow := time.Unix(int64(timeSyc), 0).Add(time.Minute * time.Duration(duration))
	return uint64(timeNow.Unix())
}

func GetOptionTypeID(duration int) int {
	if duration <= 5 {
		return optionTypeTurbo
	}
	return optionTypeBinary
}

var timeSync int

func GetPriceNow(messageToSend, responseGetCandleChan chan string, activeID int) (responsePrice chan float64, responseError chan error) {

	totalCandle := 1
	endTime := int(time.Now().Add(time.Minute * 5).Unix())
	qtdCandle := 1
	onlyClosed := false
	splitNormalization := true

	sendMessageGetCandle, _ := messages.NewSendMessageGetCandle(activeID, totalCandle, endTime, qtdCandle, onlyClosed, splitNormalization).Json()
	messageToSend <- string(sendMessageGetCandle)

	responsePrice = make(chan float64)
	responseError = make(chan error)
	go func() {
		for candleJson := range responseGetCandleChan {
			var candleResponse responseMessage.CandleResponse
			json.Unmarshal([]byte(candleJson), &candleResponse)
			if candleResponse.Status == 200 {
				responsePrice <- candleResponse.Msg.Candles[0].Close
				responseError <- nil
			}
			responseError <- errors.New("Cant get prince now")
			responsePrice <- 0.00
		}
	}()

	return
}

func OpenOrder() (string, error) {

	for timeSync == 0 {
		time.Sleep(time.Second)
	}

	duration := 5
	activeID := 76
	investiment := 2.0
	direction := "call"
	activePayoutNow := 80

	activePriceNow, err := GetPriceNow(messageToSend, responseGetCandleChan, activeID)
	if err != nil {
		log.Fatal(err)
	}
	activePrice, _ := strconv.Atoi(strings.Replace(fmt.Sprintf("%f", activePriceNow), ".", "", -1))

	binaryOptionsOpenOptionBody := messages.BinaryOptionsOpenOptionBody{
		UserBalanceID: profile.BalanceId,
		ActiveID:      activeID,
		OptionTypeID:  GetOptionTypeID(duration),
		Direction:     direction,
		Expired:       int(GetExpirationTime(timeSync, duration)),
		RefundValue:   0,
		Price:         investiment,
		Value:         activePrice,
		ProfitPercent: activePayoutNow,
	}
	binaryOptionsOpenOption, _ := messages.NewSendMessageBinaryOptionsOpenOption(binaryOptionsOpenOptionBody).Json()
	messageToSend <- string(binaryOptionsOpenOption)
}
