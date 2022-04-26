package messages

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Message struct {
	Name      string `json:"name"`
	RequestId string `json:"request_id"`
	LocalTime int    `json:"local_time"`
}

func GetRequestId() string {

	return uuid.NewString()
}

func (a Message) Json() ([]byte, error) {
	j, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return j, nil
}
