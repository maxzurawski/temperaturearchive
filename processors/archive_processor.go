package processors

import (
	"encoding/json"
	"fmt"

	"github.com/maxzurawski/temperaturearchive/dto"
	"github.com/maxzurawski/temperaturearchive/publishers"
	"github.com/maxzurawski/temperaturearchive/service"
)

func ArchiveProcessor(dto dto.MeasurementDTO) {
	measurementDTO, err := service.Service.SaveMeasurement(dto)
	if err != nil {
		publishers.Logger().Error(
			dto.ProcessId,
			dto.Uuid,
			"could not save measurement",
			err.Error())
		return
	}
	bytes, _ := json.Marshal(measurementDTO)
	publishers.Logger().Info(dto.ProcessId, dto.Uuid,
		fmt.Sprintf("measurement [%s] stored successfully", string(bytes)))
}
