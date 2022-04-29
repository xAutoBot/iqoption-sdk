package iqoptionRepository

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/url"
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
	websocketConnection     *websocket.Conn
	profile                 profile.User
	timeSyc                 int64
	authenticatedChan       chan []byte
	balancesChan            chan []byte
	candleGeneratedChan     chan []byte
	dataOpenedOrderChan     [50][]byte
	messageToSendChan       chan messageToSendStruct
	messageToSendChanStatus chan messageToSendStruct
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

	time.Local = time.UTC
	timeNow := time.Now()
	year := timeNow.Year()
	month := timeNow.Month()
	day := timeNow.Day()
	hour := timeNow.Hour()
	minute := timeNow.Minute() + duration
	second := 0

	if minute > 59 {
		hour++
		minute = minute - 60
	}
	if hour > 23 {
		day++
		hour = hour - 24
	}

	expirationTime, _ := time.Parse("2006/January/02 15:4:5", fmt.Sprintf("%v/%v/%v %v:%v:%v", year, month, day, hour, minute, second))

	return expirationTime.Unix(), nil
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

// {"name":"sendMessage","request_id":"59ce3396-1470-4709-bd3d-26793238537f", "msg":{"name":"binary-options.open-option","version":"1.0","body":{"user_balance_id":21263150,"active_id":1,"option_type_id":3,"direction":"call","expired":1651106945,"price":2}}}
// {"name":"sendMessage","request_id":"59ce3396-1470-4709-bd3d-26793238537f", "msg":{"name":"binary-options.open-option","version":"1.0","body":{"user_balance_id":21263150,"active_id":1,"option_type_id":3,"direction":"call","expired":1651105800,"price":1}}}
