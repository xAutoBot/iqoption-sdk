package messages

import "encoding/json"

type sendResults struct {
	SendResults bool `json:"sendResults"`
}

type setOptionsSendResults struct {
	Name      string      `json:"name"`
	RequestID string      `json:"request_id"`
	LocalTime int         `json:"local_time"`
	Msg       sendResults `json:"msg"`
}

func NewSetOptionsSendResults() setOptionsSendResults {
	sendResultsMsg := sendResults{
		SendResults: true,
	}
	return setOptionsSendResults{
		Name:      "setOptions",
		RequestID: GetRequestId(),
		LocalTime: GetLocalTime(),
		Msg:       sendResultsMsg,
	}
}

func (s setOptionsSendResults) Json() ([]byte, error) {
	j, error := json.Marshal(s)
	if error != nil {
		return nil, error
	}

	return j, nil
}
