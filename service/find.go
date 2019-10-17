package service

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/xdevices/temperaturearchive/dto"
	"github.com/xdevices/temperaturearchive/publishers"
)

func (s *service) Find(searchDTO dto.SearchDTO) ([]dto.MeasurementDTO, error) {
	measurements, err := s.mgr.Find(searchDTO)
	if err != nil {
		bytes, _ := json.Marshal(searchDTO)
		publishers.Logger().Error(
			uuid.New().String(),
			"",
			fmt.Sprintf("error during finding of measurements. search dto: [%s]", string(bytes)),
			err.Error())
		return nil, err
	}
	var results []dto.MeasurementDTO
	for _, item := range measurements {
		results = append(results, s.mgr.MapToDto(&item))
	}
	return results, nil
}
