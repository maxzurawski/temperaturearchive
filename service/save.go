package service

import (
	"encoding/json"
	"fmt"

	"github.com/maxzurawski/temperaturearchive/publishers"

	"github.com/maxzurawski/temperaturearchive/dto"
)

func (s *service) SaveMeasurement(dto dto.MeasurementDTO) (*dto.MeasurementDTO, error) {
	measurement, err := s.mgr.Save(dto)
	if err != nil {
		bytes, _ := json.Marshal(dto)
		publishers.Logger().Error(
			dto.ProcessId,
			dto.Uuid,
			fmt.Sprintf("could not save measurement: [%s]", string(bytes)),
			err.Error())
		return nil, err
	}
	dto = s.mgr.MapToDto(measurement)
	return &dto, nil
}
