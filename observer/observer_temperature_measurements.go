package observer

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/xdevices/temperaturearchive/config"
	"github.com/xdevices/temperaturearchive/dto"
	"github.com/xdevices/temperaturearchive/publishers"
	"github.com/xdevices/utilities/rabbit/crosscutting"
	"github.com/xdevices/utilities/rabbit/domain"
)

func TemperatureObserver(processors []func(measurementDTO dto.MeasurementDTO)) {
	observer := config.TemperaturearchiveConfig().RabbitMQManager.InitObserver()
	defer config.TemperaturearchiveConfig().RabbitMQManager.CloseConnection()
	defer observer.Channel.Close()

	observer.DeclareTopicExchange(crosscutting.TopicMeasurements.String())
	observer.BindQueue(observer.Queue, crosscutting.RoutingKeyTemperatureMeasurement.String()+".#", crosscutting.TopicMeasurements.String())
	deliveries := observer.Observe()

	for msg := range deliveries {
		measurementDTO := domain.TemperatureMeasurement{}
		err := json.Unmarshal(msg.Body, &measurementDTO)
		if err != nil {
			publishers.Logger().Error(uuid.New().String(), "", fmt.Sprintf("could not unmarshal expected measurement msg"), err.Error())
			continue
		}
		publishers.Logger().Info(measurementDTO.ProcessId, measurementDTO.SensorUuid, fmt.Sprintf("received measurement msg. [%s]", string(msg.Body)))

		for _, process := range processors {
			go process(mapToSavableDto(measurementDTO))
		}

	}
}

func mapToSavableDto(measurement domain.TemperatureMeasurement) dto.MeasurementDTO {
	measurementDTO := dto.MeasurementDTO{
		ID:         nil,
		ProcessId:  measurement.ProcessId,
		Uuid:       measurement.SensorUuid,
		ReportedAt: measurement.PublishedOn,
		Value:      measurement.Value,
		ReceivedAt: time.Now(),
	}
	return measurementDTO
}
