package messages

import "encoding/json"

type getFaqBody struct {
}
type getFaq struct {
	Name    string     `json:"name"`
	Version string     `json:"version"`
	Body    getFaqBody `json:"boby"`
}
type sendMessageGetFaq struct {
	Name      string `json:"name"`
	RequestID string `json:"request_id"`
	LocalTime int    `json:"local_time"`
	Msg       getFaq `json:"msg"`
}

func NewSendMessageGetFaq() sendMessageGetFaq {
	body := getFaqBody{}
	getfaq := getFaq{
		Name:    "get-faq",
		Version: "1.0",
		Body:    body,
	}
	return sendMessageGetFaq{
		Name:      "sendMessage",
		RequestID: GetRequestId(),
		LocalTime: GetLocalTime(),
		Msg:       getfaq,
	}
}

func (s sendMessageGetFaq) Json() ([]byte, error) {
	j, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return j, nil
}
