package service

import (
	"github.com/maxzurawski/temperaturearchive/dto"
)

func (s *service) Find(searchDTO dto.SearchDTO) ([]dto.MeasurementDTO, error) {
	measurements, err := s.mgr.Find(searchDTO)
	if err != nil {
		return nil, err
	}
	var results []dto.MeasurementDTO
	for _, item := range measurements {
		results = append(results, s.mgr.MapToDto(&item))
	}
	return results, nil
}
