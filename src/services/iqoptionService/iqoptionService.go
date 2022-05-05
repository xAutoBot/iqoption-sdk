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

func (i *IqoptionService) GetAllActiveInfo() (active.ActiveInfo, error) {

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
