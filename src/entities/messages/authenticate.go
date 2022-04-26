package messages

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"
)

type authenticateMsg struct {
	Ssid            string `json:"ssid"`
	Protocol        int    `json:"protocol"`
	SessionId       string `json:"session_id"`
	ClientSessionId string `json:"client_session_id"`
}

type authenticate struct {
	Name      string          `json:"name"`
	RequestId string          `json:"request_id"`
	LocalTime int             `json:"local_time"`
	Msg       authenticateMsg `json:"msg"`
}

func (a authenticate) Json() ([]byte, error) {
	j, err := json.Marshal(a)
	if err != nil {
		log.Fatalf(err.Error())
		return nil, err
	}
	return j, nil
}

func NewAuthenticateMsg(ssid string) authenticateMsg {
	return authenticateMsg{
		Ssid:            ssid,
		Protocol:        3,
		SessionId:       "",
		ClientSessionId: "",
	}
}

func getRandonInt(minRand int, maxRand int) int {
	rand.Seed(time.Now().UnixNano())

	return rand.Intn(maxRand-minRand+1) + minRand
}

func (a authenticate) GetRequestId() string {
	minRand := 1000000
	maxRand := 9999999

	return fmt.Sprintf("%v_%v", getRandonInt(minRand, maxRand), getRandonInt(minRand, maxRand))
}

func GetLocalTime() int {
	minRand := 1000000
	maxRand := 9999999

	return getRandonInt(minRand, maxRand)
}

func NewAuthenticate(a authenticateMsg) authenticate {
	return authenticate{
		Name:      "authenticate",
		LocalTime: GetLocalTime(),
		RequestId: (authenticate).GetRequestId(authenticate{}),
		Msg:       a,
	}
}
