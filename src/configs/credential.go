package configs

import "os"

var (
	iqOptionEmail string = os.Getenv("IQOPTION_USERNAME")
	password      string = os.Getenv("IQOPTION_PASSWORD")
	accountType   string = os.Getenv("ACCOUNT_TYPE")
)

func GetIqOptionEmail() string {
	return iqOptionEmail
}

func GetIqOptionPassword() string {
	return password
}

func GetAccountType() string {
	if accountType == "" {
		return "pratic"
	}
	return accountType
}
