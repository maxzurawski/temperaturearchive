package observer

import (
	"encoding/json"
	"fmt"

	"github.com/xdevices/temperaturearchive/cache"

	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"github.com/xdevices/temperaturearchive/config"
	"github.com/xdevices/temperaturearchive/publishers"
	"github.com/xdevices/utilities/rabbit/crosscutting"
	"github.com/xdevices/utilities/rabbit/domain"
)

func ObserveSensorChanges() {
	observer := config.TemperaturearchiveConfig().InitObserver()
	defer observer.Channel.Close()
	observer.DeclareTopicExchange(crosscutting.TopicConfigurationChanged.String())
	observer.BindQueue(observer.Queue, crosscutting.RoutingKeySensors.String()+".#", crosscutting.TopicConfigurationChanged.String())
	deliveries := observer.Observe()

	for msg := range deliveries {
		confMsg := domain.ConfigurationChanged{}
		err := json.Unmarshal(msg.Body, &confMsg)
		if err != nil {
			publishers.Logger().Error(uuid.New().String(), "", "could not update sensors cache", err.Error())
			continue
		}
		log.Info(fmt.Sprintf("Received: [%s]\n", string(msg.Body)))
		log.Info(fmt.Sprintf("Routing key: [%s]", msg.RoutingKey))
		err = cache.InitSensorsCache(confMsg.ProcessId)
		if err == nil {
			publishers.Logger().Info(confMsg.ProcessId, "", "successfully updated sensors cache")
		}
	}
}
