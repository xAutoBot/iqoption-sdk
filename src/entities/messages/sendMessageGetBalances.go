package messages

import "encoding/json"

type sendMessageGetBalancesMsgBody struct {
	TypesIds               []int `json:"types_ids"`
	TournamentsStatusesIds []int `json:"tournaments_statuses_ids"`
}
type sendMessageGetBalancesMsg struct {
	Name    string                        `json:"name"`
	Version string                        `json:"version"`
	Body    sendMessageGetBalancesMsgBody `json:"body"`
}
type sendMessageGetBalances struct {
	Name      string                    `json:"name"`
	RequestID string                    `json:"request_id"`
	LocalTime int                       `json:"local_time"`
	Msg       sendMessageGetBalancesMsg `json:"msg"`
}

func NewSendMessageGetBalances() sendMessageGetBalances {
	sendMessageGetBalancesMsgBody := sendMessageGetBalancesMsgBody{
		TypesIds:               []int{1, 4, 2},
		TournamentsStatusesIds: []int{2, 3},
	}
	sendMessageGetBalancesMsg := sendMessageGetBalancesMsg{
		Name:    "get-balances",
		Version: "1.0",
		Body:    sendMessageGetBalancesMsgBody,
	}
	return sendMessageGetBalances{
		Name:      "sendMessage",
		RequestID: GetRequestId(),
		LocalTime: GetLocalTime(),
		Msg:       sendMessageGetBalancesMsg,
	}
}

func (s sendMessageGetBalances) Json() ([]byte, error) {
	j, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return j, nil
}

//{"name":"sendMessage","request_id":"378","local_time":245736,"msg":{"name":"get-balances","version":"1.0","body":{"types_ids":[1,4,2],"tournaments_statuses_ids":[3,2]}}}
