package active

import "github.com/xAutoBot/iqoption-sdk/src/entities/messages/responseMessage"

type ActiveStatus struct {
	ActiveId      int
	BinaryIsOpen  bool
	TurboIsOpen   bool
	DigitalIsOpen bool
}

type ActiveInfo struct {
	Digital []responseMessage.Underlying
	Binary  []responseMessage.ActivesData
	Turbo   []responseMessage.ActivesData
}
