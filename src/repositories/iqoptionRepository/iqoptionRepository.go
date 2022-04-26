package iqoptionRepository

import (
	"flag"
	"log"
	"net/url"

	"github.com/evecimar/iqoptionapi/src/configs"
	"github.com/evecimar/iqoptionapi/src/entities/messages"
	"github.com/evecimar/iqoptionapi/src/repositories/authenticationRepository"
	"github.com/gorilla/websocket"
)

var websocketHost = flag.String("addr", configs.IqoptionWebSocketHost, "http service address")

func Connect() (*websocket.Conn, error) {
	ssid, err := authenticationRepository.GetSSID()

	if err != nil {
		return nil, err
	}

	websocketUrl := url.URL{Scheme: "wss", Host: *websocketHost, Path: "/echo/websocket"}
	log.Printf("connecting to %s", websocketUrl.String())

	websocketConnection, _, err := websocket.DefaultDialer.Dial(websocketUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	authenticateMessage := messages.NewAuthenticateMsg(ssid)
	ssidBody, _ := messages.NewAuthenticate(authenticateMessage).Json()
	log.Println("sending => ", string(ssidBody))
	websocketConnection.WriteMessage(websocket.TextMessage, ssidBody)

	return websocketConnection, nil
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

func OnMessage(websocketConnection *websocket.Conn, messageChan chan string) {
	for {
		_, message, err := websocketConnection.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		messageStr := string(message)
		log.Println("Received =>", messageStr)
		messageChan <- messageStr
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
