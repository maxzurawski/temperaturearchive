package dbprovider

import (
	"errors"
	"fmt"

	"github.com/xdevices/temperaturearchive/dto"
	"github.com/xdevices/temperaturearchive/model"
	"github.com/xdevices/utilities/stringutils"
)

func (mgr *manager) Save(dto dto.MeasurementDTO) (*model.Measurement, error) {
	if !stringutils.IsZero(dto.ID) {
		return nil, errors.New(fmt.Sprintf("given dto has already an Id. cannot be saved. [id: %d]", dto.ID))
	}

	if stringutils.IsZero(dto.Uuid) {
		return nil, errors.New(fmt.Sprintf("given dto has no uuid. cannot be saved."))
	}

	measurement := mgr.MapToEntity(dto)
	err := mgr.db.Create(measurement).Error
	return measurement, err
}
