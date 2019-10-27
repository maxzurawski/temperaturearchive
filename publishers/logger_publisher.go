package publishers

import (
	"github.com/labstack/gommon/log"
	"github.com/maxzurawski/temperaturearchive/config"
	"github.com/maxzurawski/utilities/rabbit/crosscutting"
	"github.com/maxzurawski/utilities/rabbit/publishing"
)

type logger struct {
	*publishing.Publisher
}

var publisher *publishing.Publisher
var loggerInstance *logger

func InitLogger() {
	if publisher == nil && config.TemperaturearchiveConfig().ConnectToRabbit() {
		publisher = config.TemperaturearchiveConfig().InitPublisher()
		publisher.DeclareTopicExchange(crosscutting.TopicLogs.String())
	}
}

func Logger() *logger {
	if loggerInstance == nil {
		loggerInstance = new(logger)
		loggerInstance.Publisher = publisher
	}
	return loggerInstance
}

func (l *logger) Info(processId, sensorUuid, msg string) {

	if !config.TemperaturearchiveConfig().ConnectToRabbit() {
		log.Info("connection to rabbit disabled")
		return
	}

	l.Reset()
	l.PublishInfo(processId,
		sensorUuid,
		config.TemperaturearchiveConfig().ServiceName(),
		msg,
	)
}

func (l *logger) Warn(processId, sensorUuid, msg string) {
	if !config.TemperaturearchiveConfig().ConnectToRabbit() {
		log.Info("connection to rabbit disabled")
		return
	}

	l.Reset()
	l.PublishWarn(processId,
		sensorUuid,
		config.TemperaturearchiveConfig().ServiceName(),
		msg,
	)
}

func (l *logger) Error(processId, sensorUuid, msg, details string) {
	if !config.TemperaturearchiveConfig().ConnectToRabbit() {
		log.Info("connection to rabbit disabled")
		return
	}

	l.Reset()
	l.PublishError(processId,
		sensorUuid,
		config.TemperaturearchiveConfig().ServiceName(),
		msg,
		details,
	)
}
