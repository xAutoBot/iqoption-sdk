package profile

import "github.com/evecimar/iqoptionapi/src/entities/messages/responseMessage"

const (
	DemoAccountType int = 4
	RealAccountType int = 1
)

type User struct {
	Balance       string                        `json:"balace"`
	BalanceId     int                           `json:"balance_id"`
	Balances      []responseMessage.BalancesMsg `json:"balances"`
	Currency      string                        `json:"currency"`
	MinimumAmount float32                       `json:"minimum_amount"`
	BalanceType   int                           `json:"balance_type"`
	CurrencyChar  string                        `json:"currency_char"`
	TimeZone      int                           `json:"time_zone"`
	Amount        float64                       `json:"amount"`
}

func (u *User) ChangeBalance(balanceType string) *User {
	newBalanceType := DemoAccountType
	if balanceType == "REAL" {
		newBalanceType = RealAccountType
	}

	for _, balance := range u.Balances {
		if balance.Type == newBalanceType {
			u.Amount = balance.Amount
			u.BalanceId = balance.ID
			u.BalanceType = newBalanceType
			u.Currency = balance.Currency
		}
	}

	return u
}
func (u User) BalanceID() int {
	if u.BalanceId == 0 {
		u.ChangeBalance("PRATIC")
	}
	return u.BalanceId
}
