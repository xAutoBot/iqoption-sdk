package orderService

import (
	"github.com/xAutoBot/iqoption-sdk/src/configs"
	"github.com/xAutoBot/iqoption-sdk/src/entities/active"
	"github.com/xAutoBot/iqoption-sdk/src/repositories/iqoptionRepository"
)

type OrderService struct {
	IqoptionRepository *iqoptionRepository.IqOptionRepository
}

func NewOrderService() (OrderService, error) {
	var orderService OrderService
	iqoption := iqoptionRepository.IqOptionRepository{}
	iqoptionRepository, err := iqoption.Connect(configs.GetAccountType())

	if err != nil {
		return orderService, err
	}
	orderService.IqoptionRepository = iqoptionRepository

	return orderService, nil
}

func (o OrderService) ActiveIsOpen(actionId int) (active.ActiveStatus, error) {
	activeStatus := active.ActiveStatus{
		ActiveId:      actionId,
		BinaryIsOpen:  false,
		DigitalIsOpen: false,
		TurboIsOpen:   false,
	}

	activesInfo, err := o.IqoptionRepository.GetAllActiveInfo()

	if err != nil {
		return activeStatus, err
	}

	for _, activeDada := range activesInfo.Binary.Actives {
		if activeDada.ID == actionId && activeDada.Enabled {
			if activeDada.IsSuspended {
				activeStatus.BinaryIsOpen = false
			}
			activeStatus.BinaryIsOpen = true
			break
		}
	}
	for _, activeDada := range activesInfo.Turbo.Actives {
		if activeDada.ID == actionId && activeDada.Enabled {
			if activeDada.IsSuspended {
				activeStatus.TurboIsOpen = false
			}
			activeStatus.TurboIsOpen = true
			break
		}
	}

	return activeStatus, nil
}
