package processors

import (
	"fmt"
	"strconv"

	"github.com/maxzurawski/temperaturearchive/cache"
	"github.com/maxzurawski/temperaturearchive/dto"
	"github.com/maxzurawski/temperaturearchive/publishers"
	"github.com/maxzurawski/utilities/stringutils"
)

func NotifierProcessor(dto dto.MeasurementDTO) {
	sensor := cache.SensorsCache.GetByUuid(dto.Uuid)

	if !stringutils.IsZero(sensor.Max) {
		acceptableMax, _ := strconv.ParseFloat(sensor.Max, 64) // ignore error -> this value should be already verified
		if dto.Value > acceptableMax {
			floatFormatted := fmt.Sprintf("%.2f", dto.Value)
			publishers.Notifier().PublishToNotifierMax(
				dto.ProcessId,
				dto.Uuid,
				sensor.Max,
				floatFormatted)
		}
	}

	if !stringutils.IsZero(sensor.Min) {
		acceptableMin, _ := strconv.ParseFloat(sensor.Min, 64) // ignore error -> this value should be already verified
		if dto.Value < acceptableMin {
			floatFormatted := fmt.Sprintf("%.2f", dto.Value)
			publishers.Notifier().PublishToNotifierMin(
				dto.ProcessId,
				dto.Uuid,
				sensor.Min,
				floatFormatted)
		}
	}
}
