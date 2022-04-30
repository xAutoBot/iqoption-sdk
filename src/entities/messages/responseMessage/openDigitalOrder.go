package responseMessage

import "encoding/json"

type OpenDigitalOrderResult struct {
	RequestID string                  `json:"request_id"`
	Name      string                  `json:"name"`
	Msg       OpenDigitalOrderDataMsg `json:"msg"`
	Status    int                     `json:"status"`
}
type OpenDigitalOrderDataMsg struct {
	ID int64 `json:"id"`
}

func (o OpenDigitalOrderDataMsg) Json() ([]byte, error) {
	j, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}
	return j, nil
}
