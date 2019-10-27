package service

import (
	"github.com/maxzurawski/temperaturearchive/dbprovider"
	"github.com/maxzurawski/temperaturearchive/dto"
)

type service struct {
	mgr dbprovider.DBManager
}

type TemperatureService interface {
	SaveMeasurement(dto dto.MeasurementDTO) (*dto.MeasurementDTO, error)
	Find(searchDTO dto.SearchDTO) ([]dto.MeasurementDTO, error)
}

var Service TemperatureService

func Init() {
	s := service{}
	s.mgr = dbprovider.Mgr
	Service = &s
}
