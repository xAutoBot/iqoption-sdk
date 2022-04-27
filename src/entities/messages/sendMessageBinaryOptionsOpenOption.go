package messages

import "encoding/json"

type BinaryOptionsOpenOptionBody struct {
	UserBalanceID int     `json:"user_balance_id"`
	ActiveID      int     `json:"active_id"`
	OptionTypeID  int     `json:"option_type_id"`
	Direction     string  `json:"direction"`
	Expired       int64   `json:"expired"`
	RefundValue   int     `json:"refund_value"`
	Price         float64 `json:"price"`
	Value         int     `json:"value"`
	ProfitPercent int     `json:"profit_percent"`
}

type binaryOptionsOpenOptionMsg struct {
	Name    string                      `json:"name"`
	Version string                      `json:"version"`
	Body    BinaryOptionsOpenOptionBody `json:"body"`
}

type sendMessageBinaryOptionsOpenOption struct {
	Name      string                     `json:"name"`
	RequestID string                     `json:"request_id"`
	LocalTime int                        `json:"local_time"`
	Msg       binaryOptionsOpenOptionMsg `json:"msg"`
}

func NewBinaryOptionsOpenOptionBody(userBalanceId int, activeId int, optionTypeId int, direction string, expired int64, refundValue int, price float64, value int, profitPercent int) BinaryOptionsOpenOptionBody {

	return BinaryOptionsOpenOptionBody{
		UserBalanceID: userBalanceId,
		ActiveID:      activeId,
		OptionTypeID:  optionTypeId,
		Direction:     direction,
		Expired:       expired,
		RefundValue:   refundValue,
		Price:         price,
		Value:         value,
		ProfitPercent: profitPercent,
	}
}

func NewSendMessageBinaryOptionsOpenOption(b BinaryOptionsOpenOptionBody) sendMessageBinaryOptionsOpenOption {

	binaryOptionsOpenOptionMsg := binaryOptionsOpenOptionMsg{
		Name:    "binary-options.open-option",
		Version: "1.0",
		Body:    b,
	}
	return sendMessageBinaryOptionsOpenOption{
		Name:      "sendMessage",
		RequestID: GetRequestId(),
		LocalTime: GetLocalTime(),
		Msg:       binaryOptionsOpenOptionMsg,
	}
}

func (s sendMessageBinaryOptionsOpenOption) Json() ([]byte, error) {
	j, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return j, nil
}

// {"name":"sendMessage","request_id":"1503","local_time":1370323,"msg":{"name":"binary-options.open-option","version":"1.0","body":{
// 	"user_balance_id":21263150,"active_id":4,"option_type_id":1,"direction":"call","expired":1650372300,"refund_value":0,"price":1.0,"value":138410745,"profit_percent":81}}}
