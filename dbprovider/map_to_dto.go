package dbprovider

import (
	"github.com/xdevices/temperaturearchive/dto"
	"github.com/xdevices/temperaturearchive/model"
)

func (mgr *manager) MapToDto(measurement *model.Measurement) dto.MeasurementDTO {
	if measurement == nil {
		return dto.MeasurementDTO{}
	}

	return dto.MeasurementDTO{
		ProcessId:  *measurement.ProcessId,
		Uuid:       *measurement.Uuid,
		ID:         measurement.ID,
		ReceivedAt: *measurement.ReceivedAt,
		ReportedAt: *measurement.ReportedAt,
		Value:      *measurement.Value,
	}
}
