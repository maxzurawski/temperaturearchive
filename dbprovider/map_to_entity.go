package dbprovider

import (
	"github.com/maxzurawski/temperaturearchive/dto"
	"github.com/maxzurawski/temperaturearchive/model"
)

func (mgr *manager) MapToEntity(dto dto.MeasurementDTO) (measurement *model.Measurement) {
	measurement = &model.Measurement{
		ID:         dto.ID,
		Uuid:       &dto.Uuid,
		ProcessId:  &dto.ProcessId,
		Value:      &dto.Value,
		ReportedAt: &dto.ReportedAt,
		ReceivedAt: &dto.ReceivedAt,
	}
	return
}
