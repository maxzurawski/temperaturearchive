package service

import "github.com/maxzurawski/temperaturearchive/dto"

// Finds last amount of entries for specific uuid different then processId
func (s *service) FindLast(amount int, uuid, processId string) ([]dto.MeasurementDTO, error) {
	entities, err := s.mgr.GetLastEntries(amount, uuid, processId)
	if err != nil {
		return nil, err
	}
	var results []dto.MeasurementDTO
	for _, item := range entities {
		results = append(results, s.mgr.MapToDto(&item))
	}
	return results, nil
}
