package dbprovider

import (
	"github.com/xdevices/temperaturearchive/dto"
	"github.com/xdevices/temperaturearchive/model"
	"github.com/xdevices/utilities/stringutils"
)

func (mgr *manager) Find(searchDTO dto.SearchDTO) ([]model.Measurement, error) {
	db := mgr.db

	var measurements []model.Measurement
	if searchDTO.LastLimited != nil {

		// NOTE: grouping last n-measurements for each of the sensors (defined by uuid)
		// if ? - lastLimited attribute is 5, and if there are 2 sensors
		// then we get 10 last results, 5 last measurements per each of the sensor
		err := db.Raw("select * "+
			"FROM "+
			"( "+
			"select *, row_number() over (partition by uuid order by reported_at desc) as n "+
			"from measurement "+
			") "+
			"where n < (? + 1)", searchDTO.LastLimited).Scan(&measurements).Error
		return measurements, err
	}

	if !stringutils.IsZero(searchDTO.Uuid) {
		db = db.Where("uuid LIKE ?", "%"+searchDTO.Uuid+"%")
	}

	if !stringutils.IsZero(searchDTO.ProcessId) {
		db = db.Where("process_id LIKE ?", "%"+searchDTO.ProcessId+"%")
	}

	if searchDTO.ValueFrom != nil {
		db = db.Where("value >= ?", searchDTO.ValueFrom)
	}

	if searchDTO.ValueTo != nil {
		db = db.Where("value <= ?", searchDTO.ValueTo)
	}

	if searchDTO.ReportedAtFrom != nil {
		db = db.Where("reported_at >= ?", searchDTO.ReportedAtFrom)
	}

	if searchDTO.ReportedAtTo != nil {
		db = db.Where("reported_at <= ?", searchDTO.ReportedAtTo)
	}

	var err error
	if *searchDTO.OrderDesc {
		err = db.Order("reported_at desc").Find(&measurements).Error
	} else {
		err = db.Order("reported_at asc").Find(&measurements).Error
	}

	return measurements, err
}
