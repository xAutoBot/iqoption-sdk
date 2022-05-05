package iqoptionRepository

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"math"
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
	clientPriceGeneratedChan   [50][]byte
	initializationDataChan     [50][]byte
	underlyingDataChan         [50][]byte
}

type messageToSendStruct struct {
	RequestID string
	Body      []byte
	Error     error
}

var websocketHost = flag.String("addr", configs.IqoptionWebSocketHost, "http service address")

func NewIqOptionRepository() (*IqOptionRepository, error) {
	iqoption := IqOptionRepository{}
	iqoptionRepository, err := iqoption.connect(configs.GetAccountType())

	return iqoptionRepository, err
}

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

func (i *IqOptionRepository) connect(accountType string) (*IqOptionRepository, error) {

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
			log.Printf("%s", receivedMessageJson)
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
				log.Printf("%s", receivedMessageJson)
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
			case "client-price-generated":
				for index := 0; index < len(i.clientPriceGeneratedChan); index++ {
					if i.clientPriceGeneratedChan[index] == nil {
						i.clientPriceGeneratedChan[index] = receivedMessageJson
						break
					}
				}
			case "initialization-data":
				for index := 0; index < len(i.initializationDataChan); index++ {
					if i.initializationDataChan[index] == nil {
						i.initializationDataChan[index] = receivedMessageJson
						break
					}
				}
			case "underlying-list":
				for index := 0; index < len(i.underlyingDataChan); index++ {
					if i.underlyingDataChan[index] == nil {
						i.underlyingDataChan[index] = receivedMessageJson
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

	for index := 0; index < 100; index++ {

		responseGeneratedCandle := <-i.candleGeneratedChan
		var responnseCandleGenerated responseMessage.ResponnseCandleGenerated
		json.Unmarshal([]byte(responseGeneratedCandle), &responnseCandleGenerated)

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

	timeZone := time.UTC
	timeNow := i.Time(timeZone)

	expirations := []time.Time{}

	if duration <= 5 {
		expirationTime := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), timeNow.Hour(), timeNow.Minute(), 0, 0, timeZone)

		if expirationTime.Add(time.Minute).Unix()-timeNow.Unix() > 30 {
			expirationTime = expirationTime.Add(time.Minute)
		} else {
			expirationTime = expirationTime.Add(time.Minute * 2)
		}
		for range [5]int{} {
			expirations = append(expirations, expirationTime)
			expirationTime = expirationTime.Add(time.Minute)

		}
	}
	timeNow = i.Time(timeZone)
	expirationTime := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), timeNow.Hour(), timeNow.Minute(), 0, 0, timeZone)
	if duration > 5 {
		for index := 0; index < 50; index++ {
			minute := expirationTime.Minute()
			if minute%15 == 0 {
				expirations = append(expirations, expirationTime)
				index++
			}
			expirationTime = expirationTime.Add(time.Minute)
		}
	}
	secodsTocloses := []int64{}
	for _, expiration := range expirations {
		remaning := (expiration.Unix() - i.Time(timeZone).Unix())
		secodsToclose := math.Abs(float64(remaning - int64(60)*int64(duration)))
		secodsTocloses = append(secodsTocloses, int64(secodsToclose))
	}

	minimuValue := secodsTocloses[0]
	minimuValueIndex := 0
	for index, secondToClose := range secodsTocloses {
		if secondToClose < minimuValue {
			minimuValue = secondToClose
			minimuValueIndex = index
		}
	}

	return expirations[minimuValueIndex].Unix(), nil
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
			log.Printf("sent -> %s", messageToSendStruct.Body)
			err := i.websocketConnection.WriteMessage(websocket.TextMessage, messageToSendStruct.Body)
			messageToSendStruct.Error = nil
			if err != nil {
				log.Fatalf(messageToSendStruct.Error.Error())
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
					return openOrderResultDataMsg, errors.New("error on open order in iqoption")
				}
				if openOrderData.Status == 4103 {
					return openOrderResultDataMsg, errors.New("error on open order in iqoption. probably the active is closed")
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

func (i *IqOptionRepository) GetInstrumentIndex(activeID int) (indexInstrument int, err error) {

	sendMessageStartPriceSplitterClientPriceGenerated, _ := messages.NewSendMessageStartPriceSplitterClientPriceGenerated(activeID).Json()
	i.SendMessage(sendMessageStartPriceSplitterClientPriceGenerated)

	for {
		for index := 0; index < len(i.clientPriceGeneratedChan); index++ {
			if i.clientPriceGeneratedChan[index] == nil {
				continue
			}
			var responseClientPriceGenerated responseMessage.ClientPriceGenerated
			json.Unmarshal([]byte(i.clientPriceGeneratedChan[index]), &responseClientPriceGenerated)

			if responseClientPriceGenerated.MicroserviceName == "price-splitter" && responseClientPriceGenerated.Msg.AssetID == activeID {
				i.clientPriceGeneratedChan[index] = nil
				indexInstrument = responseClientPriceGenerated.Msg.InstrumentIndex
				err = nil
				sendMessageStopClientPriceGenerated, _ := messages.NewSendMessageStopPriceSplitterClientPriceGenerated(activeID).Json()
				i.SendMessage(sendMessageStopClientPriceGenerated)
				return
			}
			time.Sleep(time.Millisecond)
		}
	}
	return
}

func (i *IqOptionRepository) OpenDigitalOrder(instrumentId string, activeId, instrumentIndex int, investiment float64, direction string) (openedDigitalOrderData responseMessage.OpenDigitalOrderDataMsg, responseError error) {

	var openDigitalOrderResultDataMsg responseMessage.OpenDigitalOrderDataMsg

	binaryOptionsOpenDigitalBody := messages.BinaryOptionsOpenDigitalBody{
		UserBalanceID:   i.profile.BalanceId,
		InstrumentID:    instrumentId,
		Amount:          fmt.Sprintf("%v", investiment),
		InstrumentIndex: instrumentIndex,
		AssetID:         activeId,
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
				switch openDigitalOrderData.Status {
				case 2000:
					return openDigitalOrderData.Msg, nil
				case 5000:
					return openDigitalOrderResultDataMsg, errors.New("error on open order in iqoption. Probably asset is closed")
				default:
					return openDigitalOrderResultDataMsg, errors.New("error on open order in iqoption")
				}

			}
		}
		time.Sleep(time.Second)
	}
}
func (i *IqOptionRepository) GetAllActiveDigitalInfo() (responseMessage.UnderlyingData, error) {
	sendMessageGetAllDigitalInfo := messages.NewSendMessageDigitalOptionInstrumentsGetUnderlyingList()
	sendMessageGetAllDigitalInfoJson, _ := sendMessageGetAllDigitalInfo.Json()
	i.SendMessage(sendMessageGetAllDigitalInfoJson)

	for {

		for index := 0; index < len(i.underlyingDataChan); index++ {
			if i.underlyingDataChan[index] == nil {
				continue
			}
			var underlyingList responseMessage.ResponseUnderlyingList
			json.Unmarshal(i.underlyingDataChan[index], &underlyingList)

			if underlyingList.RequestID == sendMessageGetAllDigitalInfo.RequestID {
				i.underlyingDataChan[index] = nil
				if underlyingList.Status == 2000 {
					return underlyingList.UnderlyingData, nil
				}
				return responseMessage.UnderlyingData{}, errors.New("error on get active digital info ")
			}
		}
		time.Sleep(time.Second)
	}
}

func (i *IqOptionRepository) GetAllActiveBinaryInfo() (responseMessage.InitializationData, error) {
	sendGetAllActiveInfo := messages.NewSendMessageGetInitializationData()
	sendGetAllActiveInfoJson, _ := sendGetAllActiveInfo.Json()
	i.SendMessage(sendGetAllActiveInfoJson)

	for {

		for index := 0; index < len(i.initializationDataChan); index++ {
			if i.initializationDataChan[index] == nil {
				continue
			}
			var initializationData responseMessage.ResponseInitializationData
			json.Unmarshal(i.initializationDataChan[index], &initializationData)

			if initializationData.RequestID == sendGetAllActiveInfo.RequestID {
				i.initializationDataChan[index] = nil

				return initializationData.InitializationData, nil
			}
		}
		time.Sleep(time.Second)
	}
}
