package processors

import (
	"fmt"
	"strconv"

	"github.com/maxzurawski/temperaturearchive/cache"
	"github.com/maxzurawski/temperaturearchive/publishers"
	"github.com/maxzurawski/temperaturearchive/service"
	"github.com/maxzurawski/utilities/rabbit/domain"

	"github.com/maxzurawski/temperaturearchive/dto"
)

func MinProcessor(dto dto.MeasurementDTO) {
	sensor := cache.SensorsCache.GetByUuid(dto.Uuid)

	previousMeasurement, err := service.Service.FindLast(1, dto.Uuid, dto.ProcessId)
	if err != nil {
		publishers.Logger().Warn(dto.ProcessId, dto.Uuid, err.Error())
		return
	}

	var transition domain.NotifierTransitions
	floatFormatted := fmt.Sprintf("%.2f", dto.Value)
	// NOTE: previous value was not lower then min, but current measurement is
	if !isLowerThenMin(previousMeasurement[0].Value, sensor.Min) && isLowerThenMin(dto.Value, sensor.Min) {
		transition = domain.FirstTransition
		publishers.Notifier().PublishToNotifierMin(dto.ProcessId, dto.Uuid, sensor.Min, floatFormatted, transition)
		return
	}

	// NOTE: previous value was lower then min, but current measurement is not, so temperature is stabilized
	if isLowerThenMin(previousMeasurement[0].Value, sensor.Min) && !isLowerThenMin(dto.Value, sensor.Min) {
		transition = domain.FinalTransition
		publishers.Notifier().PublishToNotifierMin(dto.ProcessId, dto.Uuid, sensor.Min, floatFormatted, transition)
		return
	}

	// NOTE: no need to bother notifier, value is not lower, was not moving back to normal, nor previous was not lower
	if !isLowerThenMin(dto.Value, sensor.Min) {
		return
	}

	// NOTE: process nacta, and eventually publish to notifier
	nacta := sensor.Nacta
	dtos, err := service.Service.FindLast(nacta, dto.Uuid, dto.ProcessId)
	if err != nil {
		publishers.Logger().Warn(dto.ProcessId, dto.Uuid, err.Error())
		return
	}
	allLower := true
	for _, i := range dtos {
		if lower := isLowerThenMin(i.Value, sensor.Min); !lower {
			allLower = false
			break
		}
	}
	if allLower {
		transition = domain.ContinuousTransition
		publishers.Notifier().PublishToNotifierMin(dto.ProcessId, dto.Uuid, sensor.Min, floatFormatted, transition)
	}
}

func isLowerThenMin(value float64, sensorMin string) bool {
	acceptableMin, _ := strconv.ParseFloat(sensorMin, 64) // ignore error -> this value should be already verified
	if value < acceptableMin {
		return true
	}
	return false
}
