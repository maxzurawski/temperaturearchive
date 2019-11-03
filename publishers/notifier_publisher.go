package publishers

import (
	"encoding/json"
	"time"

	"github.com/maxzurawski/utilities/rabbit/domain"

	"github.com/maxzurawski/temperaturearchive/config"
	"github.com/maxzurawski/utilities/rabbit/crosscutting"
	"github.com/maxzurawski/utilities/rabbit/publishing"
	"github.com/maxzurawski/utilities/stringutils"
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
	value string,
	transition domain.NotifierTransitions) {

	now := time.Now()

	msg := domain.NotifierMsg{
		ProcessId:   processId,
		SensorUuid:  sensorUuid,
		Service:     config.TemperaturearchiveConfig().ServiceName(),
		Max:         max,
		Value:       stringutils.ToMultiString(value),
		Transition:  transition,
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
	value string,
	transition domain.NotifierTransitions) {

	now := time.Now()

	msg := domain.NotifierMsg{
		ProcessId:   processId,
		SensorUuid:  sensorUuid,
		Service:     config.TemperaturearchiveConfig().ServiceName(),
		Max:         max,
		Value:       stringutils.ToMultiString(value),
		Transition:  transition,
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
