package config

import (
	"fmt"
	"os"

	"github.com/xdevices/utilities/config"
	"github.com/xdevices/utilities/rabbit"
)

type configManager struct {
	*config.Manager
	dbPath string
	rabbit.RabbitMQManager
}

var instance *configManager

func TemperaturearchiveConfig() *configManager {
	if instance == nil {
		instance = new(configManager)
		instance.Manager = new(config.Manager)
		instance.Init()
	}
	return instance
}

func (tm *configManager) Init() {
	tm.Manager.Init()
	if dbPath, err := os.LookupEnv("DB_PATH"); !err {
		panic(fmt.Sprintf("set DB_PATH and try again"))
	} else {
		tm.dbPath = dbPath
	}
}
