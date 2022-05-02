package responseMessage

type Special struct {
	Enabled bool   `json:"enabled"`
	Title   string `json:"title"`
}
type profit struct {
	Commission int `json:"commission"`
	RefundMin  int `json:"refund_min"`
	RefundMax  int `json:"refund_max"`
}
type option struct {
	Profit    profit             `json:"profit"`
	ExpTime   int                `json:"exp_time"`
	Count     int                `json:"count"`
	Special   map[string]Special `json:"special"`
	StartTime int                `json:"start_time"`
}
type minmax struct {
	Max int `json:"max"`
	Min int `json:"min"`
}
type ActivesData struct {
	Name               string  `json:"name"`
	GroupID            int     `json:"group_id"`
	Image              string  `json:"image"`
	Description        string  `json:"description"`
	Exchange           string  `json:"exchange"`
	MinimalBet         int     `json:"minimal_bet"`
	MaximalBet         int     `json:"maximal_bet"`
	TopTradersEnabled  bool    `json:"top_traders_enabled"`
	ID                 int     `json:"id"`
	Precision          int     `json:"precision"`
	Option             option  `json:"option"`
	Deadtime           int     `json:"deadtime"`
	Sum                int     `json:"sum"`
	Schedule           [][]int `json:"schedule"`
	Enabled            bool    `json:"enabled"`
	Minmax             minmax  `json:"minmax"`
	StartTime          int     `json:"start_time"`
	Provider           string  `json:"provider"`
	IsBuyback          int     `json:"is_buyback"`
	IsSuspended        bool    `json:"is_suspended"`
	BuybackDeadtime    int     `json:"buyback_deadtime"`
	RolloverEnabled    bool    `json:"rollover_enabled"`
	RolloverCommission int     `json:"rollover_commission"`
}

type Turbo struct {
	List    []interface{}          `json:"list"`
	Actives map[string]ActivesData `json:"actives"`
}

type Binary struct {
	List    []interface{}          `json:"list"`
	Actives map[string]ActivesData `json:"actives"`
}

type InitializationData struct {
	Turbo     Turbo             `json:"turbo"`
	Binary    Binary            `json:"binary"`
	Currency  string            `json:"currency"`
	IsBuyback int               `json:"is_buyback"`
	Groups    map[string]string `json:"groups"`
}
type ResponseInitializationData struct {
	RequestID          string             `json:"request_id"`
	Name               string             `json:"name"`
	InitializationData InitializationData `json:"msg"`
	Status             int                `json:"status"`
}
