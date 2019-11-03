package processors

import (
	"fmt"
	"strconv"

	"github.com/maxzurawski/utilities/rabbit/domain"

	"github.com/maxzurawski/temperaturearchive/cache"
	"github.com/maxzurawski/temperaturearchive/dto"
	"github.com/maxzurawski/temperaturearchive/publishers"
	"github.com/maxzurawski/temperaturearchive/service"
)

func MaxProcessor(dto dto.MeasurementDTO) {
	sensor := cache.SensorsCache.GetByUuid(dto.Uuid)

	previousMeasurement, err := service.Service.FindLast(1, dto.Uuid, dto.ProcessId)
	if err != nil {
		publishers.Logger().Warn(dto.ProcessId, dto.Uuid, err.Error())
		return
	}

	var transition domain.NotifierTransitions
	floatFormatted := fmt.Sprintf("%.2f", dto.Value)
	// NOTE: previous value was not higher then max, but current measurement is
	if !isValueHigherThenMax(previousMeasurement[0].Value, sensor.Max) && isValueHigherThenMax(dto.Value, sensor.Max) {
		transition = domain.FirstTransition
		publishers.Notifier().PublishToNotifierMax(dto.ProcessId, dto.Uuid, sensor.Max, floatFormatted, transition)
		return
	}

	// NOTE: previous value was higher then max, but current measurement is not, so temperature is stabilized
	if isValueHigherThenMax(previousMeasurement[0].Value, sensor.Max) && !isValueHigherThenMax(dto.Value, sensor.Max) {
		transition = domain.FinalTransition
		publishers.Notifier().PublishToNotifierMax(dto.ProcessId, dto.Uuid, sensor.Max, floatFormatted, transition)
		return
	}

	// NOTE: no need to bother notifier, value is not higher, was not moving back to normal, nor previous was not higher
	if !isValueHigherThenMax(dto.Value, sensor.Max) {
		return
	}

	// NOTE: process nacta, and eventually publish to notifier
	nacta := sensor.Nacta
	dtos, err := service.Service.FindLast(nacta, dto.Uuid, dto.ProcessId)
	if err != nil {
		publishers.Logger().Warn(dto.ProcessId, dto.Uuid, err.Error())
		return
	}
	allHigher := true
	for _, i := range dtos {
		if higher := isValueHigherThenMax(i.Value, sensor.Max); !higher {
			allHigher = false
			break
		}
	}
	if allHigher {
		transition = domain.ContinuousTransition
		publishers.Notifier().PublishToNotifierMax(dto.ProcessId, dto.Uuid, sensor.Max, floatFormatted, transition)
	}

}

func isValueHigherThenMax(value float64, sensorMax string) bool {
	acceptableMax, _ := strconv.ParseFloat(sensorMax, 64)
	if value > acceptableMax {
		return true
	}
	return false
}
