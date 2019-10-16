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
	proxyService string
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

	if proxyService, err := os.LookupEnv("PROXY_SERVICE"); !err {
		tm.proxyService = "http://localhost:8000/api"
	} else {
		tm.proxyService = proxyService
	}

	if tm.ConnectToRabbit() {
		tm.RabbitMQManager.InitConnection(tm.RabbitURL())
	}
}

func (tm *configManager) DBPath() string {
	return tm.dbPath
}

func (tm *configManager) ProxyService() string {
	return tm.proxyService
}
