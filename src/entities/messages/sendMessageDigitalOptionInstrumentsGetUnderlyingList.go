package messages

import "encoding/json"

type sendMessageDigitalOptionInstrumentsGetUnderlyingList struct {
	Name      string                                                  `json:"name"`
	RequestID string                                                  `json:"request_id"`
	LocalTime int                                                     `json:"local_time"`
	Msg       sendMessageDigitalOptionInstrumentsGetUnderlyingListMgs `json:"msg"`
}
type sendMessageDigitalOptionInstrumentsGetUnderlyingListBody struct {
	FilterSuspended bool `json:"filter_suspended"`
}
type sendMessageDigitalOptionInstrumentsGetUnderlyingListMgs struct {
	Name    string                                                   `json:"name"`
	Version string                                                   `json:"version"`
	Body    sendMessageDigitalOptionInstrumentsGetUnderlyingListBody `json:"body"`
}

func NewSendMessageDigitalOptionInstrumentsGetUnderlyingList() sendMessageDigitalOptionInstrumentsGetUnderlyingList {
	getInitializationDataMsg := sendMessageDigitalOptionInstrumentsGetUnderlyingListMgs{
		Name:    "digital-option-instruments.get-underlying-list",
		Version: "1.0",
		Body: sendMessageDigitalOptionInstrumentsGetUnderlyingListBody{
			FilterSuspended: true,
		},
	}

	return sendMessageDigitalOptionInstrumentsGetUnderlyingList{
		Name:      "sendMessage",
		RequestID: GetRequestId(),
		LocalTime: GetLocalTime(),
		Msg:       getInitializationDataMsg,
	}
}

func (s sendMessageDigitalOptionInstrumentsGetUnderlyingList) Json() ([]byte, error) {
	j, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return j, nil
}
