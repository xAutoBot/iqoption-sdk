package orderService

import (
	"github.com/xAutoBot/iqoption-sdk/src/repositories/iqoptionRepository"
)

type OrderService struct {
	IqoptionRepository *iqoptionRepository.IqOptionRepository
}

func NewOrderService() (OrderService, error) {
	var orderService OrderService
	iqoptionRepository, err := iqoptionRepository.NewIqOptionRepository()

	if err != nil {
		return orderService, err
	}
	orderService.IqoptionRepository = iqoptionRepository

	return orderService, nil
}
