package config

import (
	"github.com/xdevices/utilities/discovery"
)

type EurekaManager struct {
	discovery.Manager
}

func EurekaManagerInit() *EurekaManager {
	manager := EurekaManager{
		Manager: discovery.Manager{
			RegistrationTicket: TemperaturearchiveConfig().RegistrationTicket(),
			EurekaService:      TemperaturearchiveConfig().EurekaService(),
		},
	}
	return &manager
}
