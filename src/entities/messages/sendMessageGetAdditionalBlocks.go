package messages

import "encoding/json"

type getAdditionalBlocks struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type sendMessageGetAdditionalBlocks struct {
	Name      string              `json:"name"`
	RequestID string              `json:"request_id"`
	LocalTime int                 `json:"local_time"`
	Msg       getAdditionalBlocks `json:"msg"`
}

func NewSendMessageGetAdditionalBlocks() sendMessageGetAdditionalBlocks {
	return sendMessageGetAdditionalBlocks{
		Name:      "sendMessage",
		RequestID: GetRequestId(),
		LocalTime: GetLocalTime(),
		Msg: getAdditionalBlocks{
			Name:    "get-additional-blocks",
			Version: "1.0",
		},
	}
}

func (s sendMessageGetAdditionalBlocks) Json() ([]byte, error) {
	j, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return j, nil
}
