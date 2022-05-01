package cmd

import (
	"os"
	"os/signal"

	"github.com/xAutoBot/iqoption-sdk/src/configs"
	"github.com/xAutoBot/iqoption-sdk/src/repositories/iqoptionRepository"
)

func Start() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	iqOptionRepository := iqoptionRepository.IqOptionRepository{}
	connection, err := iqOptionRepository.Connect(configs.GetAccountType())
	if err != nil {
		panic(err)
	}

	activeId := 1
	duration := 5
	investiment := 2.00
	direction := "call"

	go connection.OpenOrder(activeId, duration, investiment, direction)

	<-interrupt
}
