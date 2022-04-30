package iqoptionRepository

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/xAutoBot/iqoption-sdk/src/configs"
	"github.com/xAutoBot/iqoption-sdk/src/entities/messages"
	"github.com/xAutoBot/iqoption-sdk/src/entities/messages/responseMessage"
	"github.com/xAutoBot/iqoption-sdk/src/entities/profile"
	"github.com/xAutoBot/iqoption-sdk/src/repositories/authenticationRepository"
)

const (
	optionTypeBinary int = 1
	optionTypeTurbo  int = 3
)

type IqOptionRepository struct {
	websocketConnection        *websocket.Conn
	profile                    profile.User
	timeSyc                    int64
	authenticatedChan          chan []byte
	balancesChan               chan []byte
	candleGeneratedChan        chan []byte
	dataOpenedOrderChan        [50][]byte
	dataOpenedDigitalOrderChan [50][]byte
	messageToSendChan          chan messageToSendStruct
	messageToSendChanStatus    chan messageToSendStruct
	time                       time.Time
}

type messageToSendStruct struct {
	RequestID string
	Body      []byte
	Error     error
}

var websocketHost = flag.String("addr", configs.IqoptionWebSocketHost, "http service address")

func (i *IqOptionRepository) authenticate(ssid string) error {
	i.authenticatedChan = make(chan []byte)
	defer close(i.authenticatedChan)

	authenticateMessage := messages.NewAuthenticateMsg(ssid)
	ssidBody := messages.NewAuthenticate(authenticateMessage)
	ssidJson, _ := ssidBody.Json()
	i.SendMessage(ssidJson)

	var authenticatedMessage responseMessage.Authenticated
	json.Unmarshal(<-i.authenticatedChan, &authenticatedMessage)

	if !authenticatedMessage.Msg {
		return errors.New("authentication failed")
	}

	return nil
}

func (i *IqOptionRepository) Connect(accountType string) (*IqOptionRepository, error) {

	var err error

	ssidChan, errChan := authenticationRepository.GetSSID()

	websocketUrl := url.URL{Scheme: "wss", Host: *websocketHost, Path: "/echo/websocket"}
	log.Printf("connecting to %s", websocketUrl.String())
	i.startSendMessageloop()
	i.websocketConnection, _, err = websocket.DefaultDialer.Dial(websocketUrl.String(), nil)

	if err != nil {
		return nil, err
	}
	ssid := <-ssidChan
	if err = <-errChan; err != nil {
		i.websocketConnection.Close()
		return nil, err
	}

	i.startReadResponseMessage()

	err = i.authenticate(ssid)
	if err != nil {
		i.websocketConnection.Close()
		return nil, err
	}

	balances, err := i.GetBalances()
	if err != nil {
		i.websocketConnection.Close()
		return nil, err
	}
	i.profile.Balances = balances
	i.profile.ChangeBalance(accountType)

	return i, nil
}

func (i *IqOptionRepository) GetBalances() ([]responseMessage.BalancesMsg, error) {

	i.balancesChan = make(chan []byte)
	defer close(i.balancesChan)

	messageGetBalances := messages.NewSendMessageGetBalances()
	messageGetBalancesJson, _ := messageGetBalances.Json()
	i.SendMessage(messageGetBalancesJson)

	var balanceMessage responseMessage.Balannces
	json.Unmarshal(<-i.balancesChan, &balanceMessage)

	return balanceMessage.Msg, nil
}

func (i IqOptionRepository) TimeSync() int64 {
	return i.timeSyc
}

func (i *IqOptionRepository) startReadResponseMessage() {

	go func() {
		defer i.websocketConnection.Close()

		for {
			_, receivedMessageJson, _ := i.websocketConnection.ReadMessage()
			var receivedMessage messages.Message
			json.Unmarshal(receivedMessageJson, &receivedMessage)
			switch receivedMessage.Name {
			case "heartbeat":
			case "timeSync":
				var responseTimeSync responseMessage.TimeSync
				json.Unmarshal(receivedMessageJson, &responseTimeSync)
				i.timeSyc = responseTimeSync.Msg
			case "balances":
				i.balancesChan <- receivedMessageJson
			case "authenticated":
				i.authenticatedChan <- receivedMessageJson
			case "candle-generated":
				i.candleGeneratedChan <- receivedMessageJson
			case "option":
				for index := 0; index < len(i.dataOpenedOrderChan); index++ {
					if i.dataOpenedOrderChan[index] == nil {
						i.dataOpenedOrderChan[index] = receivedMessageJson
						break
					}
				}
			case "digital-option-placed":
				for index := 0; index < len(i.dataOpenedDigitalOrderChan); index++ {
					if i.dataOpenedDigitalOrderChan[index] == nil {
						i.dataOpenedDigitalOrderChan[index] = receivedMessageJson
						break
					}
				}
			}

		}
	}()
}

func (i *IqOptionRepository) GetPriceNow(activeID int) (responsePrice float64, responseError error) {

	i.candleGeneratedChan = make(chan []byte)

	candleSize := 5
	sendMessageStartCandleGenerate, _ := messages.NewSendMessageStartCandleGenerate(activeID, candleSize).Json()
	i.SendMessage(sendMessageStartCandleGenerate)

	responseGeneratedCandle := <-i.candleGeneratedChan
	var responnseCandleGenerated responseMessage.ResponnseCandleGenerated
	json.Unmarshal([]byte(responseGeneratedCandle), &responnseCandleGenerated)

	for index := 0; index < 100; index++ {

		if responnseCandleGenerated.MicroserviceName == "quotes" && responnseCandleGenerated.Msg.ActiveID == activeID {
			responsePrice = responnseCandleGenerated.Msg.Close
			responseError = nil
			sendMessageStopCandleGenerate, _ := messages.NewSendMessageStopCandleGenerate(activeID, candleSize).Json()
			i.SendMessage(sendMessageStopCandleGenerate)

			return
		}
		time.Sleep(time.Millisecond)
	}
	return
}

func (i *IqOptionRepository) GetOptionTypeID(duration int) int {
	if duration <= 5 {
		return optionTypeTurbo
	}
	return optionTypeBinary
}

//Return the timestamp of the sum timenow with duration
func (i IqOptionRepository) GetExpirationTime(duration int) (int64, error) {

	if duration > 5 {
		return 0, errors.New("max time duration is 5 minutes yet")
	}

	timeZone := time.UTC
	timeNow := i.Time(timeZone)
	year := timeNow.Year()
	month := timeNow.Month()
	day := timeNow.Day()
	hour := timeNow.Hour()
	minute := timeNow.Minute()
	second := timeNow.Second()

	switch {
	case duration == 1:
		if second >= 40 {
			minute += 2
		} else {
			minute++
		}
	case duration > 1:
		minute += duration
	}
	second = 0

	lastDayOfMonth := getLastDayOfMonth(year, int(month), timeZone)
	if minute > 59 {
		hour++
		minute = minute - 60
	}
	if hour > 23 {
		day++
		hour = hour - 24
	}
	if day > lastDayOfMonth {
		month++
		day = day - lastDayOfMonth
	}
	if month > 12 {
		year++
		month = month - 12
	}
	second = 0
	expirationTime, _ := time.Parse("2006/January/2 15:4:5", fmt.Sprintf("%v/%v/%v %v:%v:%v", year, month, day, hour, minute, second))

	return expirationTime.Unix(), nil
}

func (i IqOptionRepository) Time(timeZone *time.Location) time.Time {
	if i.time.IsZero() {
		time.Local = timeZone
		return time.Now()
	}
	return i.time
}

func getLastDayOfMonth(year, month int, timeZone *time.Location) int {
	return time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, timeZone).Day()
}
func (i IqOptionRepository) GetDiginalExpirationTime(duration uint8) (time.Time, error) {
	timeZone := time.UTC
	timeNow := i.Time(timeZone)
	year := timeNow.Year()
	month := int(timeNow.Month())
	day := timeNow.Day()
	hour := timeNow.Hour()
	minute := timeNow.Minute()
	second := 0

	switch {
	case duration == 1:
		if timeNow.Second() >= 40 {
			timeNow = timeNow.Add(time.Minute * 2)
			return timeNow, nil
		}
		timeNow = timeNow.Add(time.Minute * 1)
		return timeNow, nil

	case duration == 5:
		switch {
		case timeNow.Minute() >= 55:
			hour++
			minute = 00
		case timeNow.Minute() >= 50:
			minute = 55
		case timeNow.Minute() >= 45:
			minute = 50
		case timeNow.Minute() >= 40:
			minute = 45
		case timeNow.Minute() >= 35:
			minute = 40
		case timeNow.Minute() >= 30:
			minute = 35
		case timeNow.Minute() >= 25:
			minute = 30
		case timeNow.Minute() >= 20:
			minute = 25
		case timeNow.Minute() >= 15:
			minute = 20
		case timeNow.Minute() >= 10:
			minute = 15
		case timeNow.Minute() >= 5:
			minute = 10
		default:
			minute = 05
		}
	case duration == 15:
		switch {
		case timeNow.Minute() >= 45:
			hour++
			minute = 00
		case timeNow.Minute() >= 30:
			minute = 45
		case timeNow.Minute() >= 15:
			minute = 30
		default:
			minute = 15
		}
	default:
		return timeNow, errors.New("invalid duration for digital option")
	}

	lastDayOfMonth := getLastDayOfMonth(year, month, timeZone)
	if hour >= 24 {
		day++
		hour = hour - 24
	}
	if day > lastDayOfMonth {
		month++
		day = day - lastDayOfMonth
	}
	if month > 12 {
		year++
		month = month - 12
	}
	timeNow, err := time.Parse("2006/1/2 15:4:5", fmt.Sprintf("%v/%v/%v %v:%v:%v", year, month, day, hour, minute, second))

	return timeNow, err
}

func (i *IqOptionRepository) startSendMessageloop() {

	i.messageToSendChan = make(chan messageToSendStruct, 10)
	i.messageToSendChanStatus = make(chan messageToSendStruct, 10)
	go func() {
		for {
			messageToSendStruct := <-i.messageToSendChan

			err := i.websocketConnection.WriteMessage(websocket.TextMessage, messageToSendStruct.Body)
			messageToSendStruct.Error = nil
			if err != nil {
				messageToSendStruct.Error = err
			}
			i.messageToSendChanStatus <- messageToSendStruct
		}
	}()
}
func (i *IqOptionRepository) SendMessage(message []byte) error {
	requestId := uuid.NewString()
	i.messageToSendChan <- messageToSendStruct{
		Body:      message,
		RequestID: requestId,
		Error:     nil,
	}

	for {
		response := <-i.messageToSendChanStatus
		if response.RequestID == requestId {
			return response.Error
		}
		i.messageToSendChanStatus <- response
	}
}

func (i IqOptionRepository) CloseConnection() error {
	log.Print("closing conection")
	err := i.websocketConnection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		return err
	}
	return nil
}

func (i *IqOptionRepository) OpenOrder(activeId, duration int, investiment float64, direction string) (openedOrderData responseMessage.OpenOrderDataMsg, responseError error) {

	var openOrderResultDataMsg responseMessage.OpenOrderDataMsg

	expirationTime, _ := i.GetExpirationTime(duration)
	binaryOptionsOpenOptionBody := messages.BinaryOptionsOpenOptionBody{
		UserBalanceID: i.profile.BalanceId,
		ActiveID:      activeId,
		OptionTypeID:  i.GetOptionTypeID(duration),
		Direction:     direction,
		Expired:       expirationTime,
		RefundValue:   0,
		Price:         investiment,
		Value:         0,
		ProfitPercent: 0,
	}

	binaryOptionsOpenOption := messages.NewSendMessageBinaryOptionsOpenOption(binaryOptionsOpenOptionBody)
	binaryOptionsOpenOptionJson, _ := binaryOptionsOpenOption.Json()
	i.SendMessage(binaryOptionsOpenOptionJson)

	for {

		for index := 0; index < len(i.dataOpenedOrderChan); index++ {
			if i.dataOpenedOrderChan[index] == nil {
				continue
			}
			var openOrderData responseMessage.OpenOrderData
			json.Unmarshal(i.dataOpenedOrderChan[index], &openOrderData)

			if openOrderData.RequestID == binaryOptionsOpenOption.RequestID {
				i.dataOpenedOrderChan[index] = nil
				if openOrderData.Status != 2000 {
					return openOrderResultDataMsg, errors.New("Error on open order in iqoption")
				}
				return openOrderData.Msg, nil
			}
		}
		time.Sleep(time.Second)
	}
}

func (i IqOptionRepository) GetDigitalInstrumentID(activeId, duration uint8, direction string) (string, error) {

	action := "P"
	switch strings.ToUpper(direction) {
	case "CALL":
		action = "C"
	case "PUT":
		action = "P"
	default:
		return "", errors.New("invalid direction. You should to use PUT or CALL")
	}

	expirationDateTime, err := i.GetDiginalExpirationTime(duration)
	if err != nil {
		return "", err
	}
	dateNow := expirationDateTime.Format("20060102")
	timeNow := expirationDateTime.Format("150400")
	dateTimeFormated := fmt.Sprintf("%vD%v", dateNow, timeNow)

	return fmt.Sprintf("do%vA%vT%vM%vSPT", activeId, dateTimeFormated, duration, action), nil
}
func (i *IqOptionRepository) OpenDigitalOrder(activeId, duration int, investiment float64, direction string) (openedDigitalOrderData responseMessage.OpenDigitalOrderDataMsg, responseError error) {
	var openDigitalOrderResultDataMsg responseMessage.OpenDigitalOrderDataMsg

	instrumentID := "teste"
	indestiment := "2.5"
	instrumentIndex := 1231231
	assetId := 12345
	// expirationTime, _ := i.GetExpirationTime(duration)

	binaryOptionsOpenDigitalBody := messages.BinaryOptionsOpenDigitalBody{
		UserBalanceID:   i.profile.BalanceId,
		InstrumentID:    instrumentID,
		Amount:          indestiment,
		InstrumentIndex: instrumentIndex,
		AssetID:         assetId,
	}

	binaryOptionsOpenDigital := messages.NewSendMessageBinaryOptionsOpenDigital(binaryOptionsOpenDigitalBody)
	binaryOptionsOpenDigitalJson, _ := binaryOptionsOpenDigital.Json()
	i.SendMessage(binaryOptionsOpenDigitalJson)

	for {

		for index := 0; index < len(i.dataOpenedDigitalOrderChan); index++ {
			if i.dataOpenedDigitalOrderChan[index] == nil {
				continue
			}
			var openDigitalOrderData responseMessage.OpenDigitalOrderResult
			json.Unmarshal(i.dataOpenedDigitalOrderChan[index], &openDigitalOrderData)

			if openDigitalOrderData.RequestID == binaryOptionsOpenDigital.RequestID {
				i.dataOpenedOrderChan[index] = nil
				if openDigitalOrderData.Status != 2000 {
					return openDigitalOrderResultDataMsg, errors.New("Error on open order in iqoption")
				}
				return openDigitalOrderData.Msg, nil
			}
		}
		time.Sleep(time.Second)
	}
}
