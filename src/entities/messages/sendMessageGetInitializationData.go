package messages

import "encoding/json"

type sendMessageGetInitializationData struct {
	Name      string                   `json:"name"`
	RequestID string                   `json:"request_id"`
	LocalTime int                      `json:"local_time"`
	Msg       getInitializationDataMsg `json:"msg"`
}
type body struct {
}
type getInitializationDataMsg struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Body    body   `json:"body"`
}

func NewSendMessageGetInitializationData() sendMessageGetInitializationData {
	getInitializationDataMsg := getInitializationDataMsg{
		Name:    "get-initialization-data",
		Version: "3.0",
		Body:    body{},
	}

	return sendMessageGetInitializationData{
		Name:      "sendMessage",
		RequestID: GetRequestId(),
		LocalTime: GetLocalTime(),
		Msg:       getInitializationDataMsg,
	}
}

func (s sendMessageGetInitializationData) Json() ([]byte, error) {
	j, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return j, nil
}
