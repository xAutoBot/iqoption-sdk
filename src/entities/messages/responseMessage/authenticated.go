package responseMessage

import (
	"encoding/json"

	"github.com/evecimar/iqoptionapi/src/entities/messages"
)

type Authenticated struct {
	messages.Message
	Msg             bool   `json:"msg"`
	ClientSessionId string `json:"client_session_id"`
}

func (a Authenticated) Json() ([]byte, error) {
	j, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return j, nil
}
