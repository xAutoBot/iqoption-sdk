package responseMessage

import "encoding/json"

type OpenOrderResult struct {
	Name      string             `json:"name"`
	RequestID string             `json:"request_id"`
	Msg       OpenOrderResultMsg `json:"msg"`
}
type OpenOrderResultMsg struct {
	Success bool `json:"success"`
}

type OpenOrderData struct {
	RequestID string           `json:"request_id"`
	Name      string           `json:"name"`
	Msg       OpenOrderDataMsg `json:"msg"`
	Status    int              `json:"status"`
}
type OpenOrderDataMsg struct {
	UserID             int64       `json:"user_id"`
	ID                 int64       `json:"id"`
	RefundValue        int         `json:"refund_value"`
	Price              float64     `json:"price"`
	Exp                int64       `json:"exp"`
	Created            int64       `json:"created"`
	CreatedMillisecond int64       `json:"created_millisecond"`
	TimeRate           int64       `json:"time_rate"`
	Type               string      `json:"type"`
	Act                int         `json:"act"`
	Direction          string      `json:"direction"`
	ExpValue           int64       `json:"exp_value"`
	Value              float64     `json:"value"`
	ProfitIncome       int         `json:"profit_income"`
	ProfitReturn       int         `json:"profit_return"`
	RobotID            interface{} `json:"robot_id"`
	ClientPlatformID   int         `json:"client_platform_id"`
}

func (o OpenOrderDataMsg) Json() ([]byte, error) {
	j, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}
	return j, nil
}
func (o OpenOrderData) Json() []byte {
	j, err := json.Marshal(o)
	if err != nil {
		return nil
	}
	return j
}
