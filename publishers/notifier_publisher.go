package publishers

import (
	"encoding/json"
	"time"

	"github.com/xdevices/utilities/rabbit/domain"

	"github.com/xdevices/temperaturearchive/config"
	"github.com/xdevices/utilities/rabbit/crosscutting"
	"github.com/xdevices/utilities/rabbit/publishing"
	"github.com/xdevices/utilities/stringutils"
)

type notifier struct {
	*publishing.Publisher
}

var notifierPublisher *publishing.Publisher
var notifierPublisherInstance *notifier

func InitNotifier() {
	if notifierPublisher == nil && config.TemperaturearchiveConfig().ConnectToRabbit() {
		notifierPublisher = config.TemperaturearchiveConfig().InitPublisher()
		notifierPublisher.DeclareTopicExchange(crosscutting.TopicNotifier.String())
	}
}

func Notifier() *notifier {
	if notifierPublisherInstance == nil {
		notifierPublisherInstance = new(notifier)
		notifierPublisherInstance.Publisher = notifierPublisher
	}
	return notifierPublisherInstance
}

func (n *notifier) PublishToNotifierMax(
	processId,
	sensorUuid,
	max,
	value string) {

	now := time.Now()

	msg := domain.NotifierMsg{
		ProcessId:   processId,
		SensorUuid:  sensorUuid,
		Service:     config.TemperaturearchiveConfig().ServiceName(),
		Max:         max,
		Value:       stringutils.ToMultiString(value),
		PublishedOn: now,
	}

	bytes, _ := json.Marshal(msg)
	n.Reset()
	n.Publish(
		crosscutting.TopicNotifier.String(),
		crosscutting.RoutingKeyNotifierTemperatureMax.String(),
		string(bytes),
	)

	n.Reset()
	n.PublishExtInfo(msg)
}

func (n *notifier) PublishToNotifierMin(
	processId,
	sensorUuid,
	max,
	value string) {

	now := time.Now()

	msg := domain.NotifierMsg{
		ProcessId:   processId,
		SensorUuid:  sensorUuid,
		Service:     config.TemperaturearchiveConfig().ServiceName(),
		Max:         max,
		Value:       stringutils.ToMultiString(value),
		PublishedOn: now,
	}

	bytes, _ := json.Marshal(msg)
	n.Reset()
	n.Publish(
		crosscutting.TopicNotifier.String(),
		crosscutting.RoutingKeyNotifierTemperatureMin.String(),
		string(bytes),
	)

	n.Reset()
	n.PublishExtInfo(msg)
}
