package responseMessage

import "encoding/json"

type Balannces struct {
	RequestID string        `json:"request_id"`
	Name      string        `json:"name"`
	Msg       []BalancesMsg `json:"msg"`
	Status    int           `json:"status"`
}
type BalancesMsg struct {
	ID                int         `json:"id"`
	UserID            int         `json:"user_id"`
	Type              int         `json:"type"`
	Amount            float64     `json:"amount"`
	EnrolledAmount    float64     `json:"enrolled_amount"`
	EnrolledSumAmount float64     `json:"enrolled_sum_amount"`
	HoldAmount        int         `json:"hold_amount"`
	OrdersAmount      int         `json:"orders_amount"`
	Currency          string      `json:"currency"`
	TournamentID      interface{} `json:"tournament_id"`
	TournamentName    interface{} `json:"tournament_name"`
	IsFiat            bool        `json:"is_fiat"`
	IsMarginal        bool        `json:"is_marginal"`
	HasDeposits       bool        `json:"has_deposits"`
	AuthAmount        int         `json:"auth_amount"`
	Equivalent        int         `json:"equivalent"`
}

func (b BalancesMsg) Json() ([]byte, error) {
	j, err := json.Marshal(b)
	if err != nil {
		return nil, err
	}
	return j, nil
}
