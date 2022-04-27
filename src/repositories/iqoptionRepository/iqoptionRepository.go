package iqoptionRepository

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/xAutoBot/iqoption-sdk/src/configs"
	"github.com/xAutoBot/iqoption-sdk/src/entities/messages"
	"github.com/xAutoBot/iqoption-sdk/src/entities/messages/responseMessage"
	"github.com/xAutoBot/iqoption-sdk/src/entities/profile"
	"github.com/xAutoBot/iqoption-sdk/src/repositories/authenticationRepository"
)

type IqOptionRepository struct {
	websocketConnection *websocket.Conn
	profile             profile.User
	timeSyc             int64
	authenticatedChan   chan []byte
	balancesChan        chan []byte
	candleGeneratedChan chan []byte
}

var websocketHost = flag.String("addr", configs.IqoptionWebSocketHost, "http service address")

func (i *IqOptionRepository) authenticate(ssid string) error {
	i.authenticatedChan = make(chan []byte)
	defer close(i.authenticatedChan)

	authenticateMessage := messages.NewAuthenticateMsg(ssid)
	ssidBody := messages.NewAuthenticate(authenticateMessage)
	ssidJson, _ := ssidBody.Json()
	i.websocketConnection.WriteMessage(websocket.TextMessage, ssidJson)

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
	i.websocketConnection.WriteMessage(websocket.TextMessage, messageGetBalancesJson)

	var balanceMessage responseMessage.Balannces
	json.Unmarshal(<-i.balancesChan, &balanceMessage)

	return balanceMessage.Msg, nil
}

func (i IqOptionRepository) TimeSync() int64 {
	return i.timeSyc
}

func (i *IqOptionRepository) startReadResponseMessage() {

	go func() {

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
				println(i.timeSyc)
			case "balances":
				i.balancesChan <- receivedMessageJson
			case "authenticated":
				i.authenticatedChan <- receivedMessageJson
			case "candle-generated":
				i.candleGeneratedChan <- receivedMessageJson
			}
		}
	}()
}

// func (i *IqOptionRepository) GetPriceNow(activeID int) (responsePrice chan float64, responseError chan error) {

// 	i.candleGeneratedChan = make(chan []byte)
// 	defer close(i.candleGeneratedChan)

// 	candleSize := 5
// 	sendMessageStartCandleGenerate, _ := messages.NewSendMessageStartCandleGenerate(activeID, candleSize).Json()
// 	i.websocketConnection.WriteMessage(websocket.TextMessage, sendMessageStartCandleGenerate)

// 	responseGeneratedCandle := <-i.candleGeneratedChan
// 	var responnseCandleGenerated responseMessage.ResponnseCandleGenerated
// 	json.Unmarshal([]byte(candleJson), &responnseCandleGenerated)
// 	if responnseCandleGenerated.MicroserviceName == "quotes" && responnseCandleGenerated.Msg.ActiveID == activeID {
// 		responsePrice <- responnseCandleGenerated.Msg.Close
// 		responseError <- nil
// 		sendMessageStopCandleGenerate, _ := messages.NewSendMessageStopCandleGenerate(activeID, candleSize).Json()
// 		messageToSend <- string(sendMessageStopCandleGenerate)

// 		return

// 	}

// 	return
// }

//Return the timestamp of the sum timeSyn with duration
func (i *IqOptionRepository) GetExpirationTime(timeSyc int64, duration int) int64 {
	timeNow := fmt.Sprintf("%d", time.Unix(timeSyc, 0).Add(time.Minute*time.Duration(duration)).Unix())
	timeNowRunes := string([]rune(timeNow)[0:10])
	timeNowInt64, _ := strconv.ParseInt(timeNowRunes, 10, 64)
	return timeNowInt64
}

func SendMessage(websocketConnection *websocket.Conn, messageChan chan string) {

	for message := range messageChan {
		log.Println("Sent =>", message)
		err := websocketConnection.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println("write:", err)
			return
		}
	}

}

func CloseConnection(websocketConnection *websocket.Conn) error {
	log.Print("closing conection")
	err := websocketConnection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		return err
	}
	return nil
}

func ChangeBalance(websocketConnection *websocket.Conn, balance string) {

	messageChan := make(chan string)
	messageChan <- balance
	go SendMessage(websocketConnection, messageChan)
}
