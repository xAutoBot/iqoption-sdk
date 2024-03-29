package iqoptionService

import (
	"sync"

	"github.com/xAutoBot/iqoption-sdk/src/entities/active"
	"github.com/xAutoBot/iqoption-sdk/src/entities/messages/responseMessage"
	"github.com/xAutoBot/iqoption-sdk/src/repositories/iqoptionRepository"
)

type IqoptionService struct {
	IqOptionRepository *iqoptionRepository.IqOptionRepository
}

func NewIqoptionService() (IqoptionService, error) {
	var iqoptionService IqoptionService

	iqOptionRepository, err := iqoptionRepository.NewIqOptionRepository()
	iqoptionService.IqOptionRepository = iqOptionRepository

	return iqoptionService, err
}

func (i IqoptionService) GetAllActiveInfo() (active.ActiveInfo, error) {

	allActiveDigitalInfoChan := make(chan responseMessage.UnderlyingData)
	allActiveBinaryInfoChan := make(chan responseMessage.InitializationData)
	allActiveDigitalInfoErrChan := make(chan error)
	allActiveBinaryInfoErrChan := make(chan error)

	go func() {
		allActiveDigitalInfo, err := i.IqOptionRepository.GetAllActiveDigitalInfo()
		allActiveDigitalInfoErrChan <- err
		allActiveDigitalInfoChan <- allActiveDigitalInfo
	}()
	go func() {
		allActiveBinaryInfo, err := i.IqOptionRepository.GetAllActiveBinaryInfo()
		allActiveBinaryInfoErrChan <- err
		allActiveBinaryInfoChan <- allActiveBinaryInfo
	}()
	allActiveDigitalInfoErr := <-allActiveDigitalInfoErrChan
	allActiveBinaryInfoErr := <-allActiveBinaryInfoErrChan
	allActiveDigitalInfo := <-allActiveDigitalInfoChan
	allActiveBinaryInfo := <-allActiveBinaryInfoChan

	if allActiveDigitalInfoErr != nil {
		return active.ActiveInfo{}, allActiveDigitalInfoErr
	}
	if allActiveBinaryInfoErr != nil {
		return active.ActiveInfo{}, allActiveDigitalInfoErr
	}

	binaryActives := make([]responseMessage.ActivesData, 0)
	turboActives := make([]responseMessage.ActivesData, 0)
	digitalActives := allActiveDigitalInfo.Underlying

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		for _, active := range allActiveBinaryInfo.Binary.Actives {
			binaryActives = append(binaryActives, active)
		}
		wg.Done()
	}()
	go func() {
		for _, active := range allActiveBinaryInfo.Turbo.Actives {
			turboActives = append(turboActives, active)
		}
		wg.Done()
	}()
	wg.Wait()

	activeInfo := active.ActiveInfo{
		Turbo:   turboActives,
		Binary:  binaryActives,
		Digital: digitalActives,
	}

	return activeInfo, nil
}

func (i IqoptionService) ActiveIsOpen(activeId int) (active.ActiveStatus, error) {
	activeStatus := active.ActiveStatus{
		DigitalIsOpen: false,
		BinaryIsOpen:  false,
		TurboIsOpen:   false,
		ActiveId:      activeId,
	}
	allActiveInfo, err := i.GetAllActiveInfo()

	if err != nil {
		return activeStatus, err
	}

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {

		for _, binaryActiveInfo := range allActiveInfo.Binary {
			if binaryActiveInfo.ID == activeId && binaryActiveInfo.Enabled {
				if !binaryActiveInfo.IsSuspended {
					activeStatus.BinaryIsOpen = true
				}
			}
		}
		wg.Done()
	}()
	go func() {

		for _, turboActiveInfo := range allActiveInfo.Turbo {
			if turboActiveInfo.ID == activeId && turboActiveInfo.Enabled {
				if !turboActiveInfo.IsSuspended {
					activeStatus.TurboIsOpen = true
				}
			}
		}
		wg.Done()
	}()
	go func() {

		for _, digitalActiveInfo := range allActiveInfo.Digital {
			if digitalActiveInfo.ActiveID == activeId && digitalActiveInfo.IsEnabled {
				if !digitalActiveInfo.IsSuspended {
					activeStatus.DigitalIsOpen = true
				}
			}
		}
		wg.Done()
	}()

	wg.Wait()

	return activeStatus, nil
}
