package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/maxzurawski/utilities/resterror"

	"github.com/maxzurawski/temperaturearchive/service"

	"github.com/maxzurawski/temperaturearchive/publishers"

	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/maxzurawski/temperaturearchive/dto"
	"github.com/maxzurawski/utilities/stringutils"
)

func HandleFind(c echo.Context) error {
	searchDto := dto.SearchDTO{}

	lastLimited := c.QueryParam("lastLimited")
	if !stringutils.IsZero(lastLimited) {
		value, err := strconv.Atoi(lastLimited)
		if err != nil {
			value = 10
		}
		searchDto.LastLimited = &value
	}

	sensorUuid := c.QueryParam("uuid")
	if !stringutils.IsZero(sensorUuid) {
		searchDto.Uuid = sensorUuid
	}

	processId := c.QueryParam("processId")
	if !stringutils.IsZero(processId) {
		searchDto.ProcessId = processId
	}

	valueFrom := c.QueryParam("valueFrom")
	if !stringutils.IsZero(valueFrom) {
		value, err := strconv.ParseFloat(valueFrom, 64)
		if err != nil {
			publishers.Logger().Warn(
				uuid.New().String(),
				sensorUuid,
				fmt.Sprintf("could not parse [valueFrom] to float. Error: [%s]", err.Error()))
		} else {
			searchDto.ValueFrom = &value
		}
	}

	valueTo := c.QueryParam("valueTo")
	if !stringutils.IsZero(valueTo) {
		value, err := strconv.ParseFloat(valueTo, 64)
		if err != nil {
			publishers.Logger().Warn(
				uuid.New().String(),
				sensorUuid,
				fmt.Sprintf("could not parse [valueTo] to float. Error: [%s]", err.Error()))
		} else {
			searchDto.ValueTo = &value
		}
	}

	reportedAtFrom := c.QueryParam("reportedAtFrom")
	if !stringutils.IsZero(reportedAtFrom) {
		parse, err := time.Parse(time.RFC3339, reportedAtFrom)
		if err != nil {
			publishers.Logger().Warn(
				uuid.New().String(),
				sensorUuid,
				fmt.Sprintf("could not parse [reportedAtFrom] into time. value to parse: [%s], "+
					"Error: [%s]", reportedAtFrom, err.Error()))
		} else {
			searchDto.ReportedAtFrom = &parse
		}
	}

	reportedAtTo := c.QueryParam("reportedAtTo")
	if !stringutils.IsZero(reportedAtTo) {
		parse, err := time.Parse(time.RFC3339, reportedAtTo)
		if err != nil {
			publishers.Logger().Warn(
				uuid.New().String(),
				sensorUuid,
				fmt.Sprintf("could not parse [reportedAtTo] into time. value to parse: [%s],"+
					"Error: [%s]", reportedAtTo, err.Error()))
		} else {
			searchDto.ReportedAtTo = &parse
		}
	}

	orderDesc := c.QueryParam("orderDesc")
	desc := false
	searchDto.OrderDesc = &desc
	if !stringutils.IsZero(orderDesc) {
		parseBool, err := strconv.ParseBool(orderDesc)

		if err != nil {
			parseBool = true
		}

		searchDto.OrderDesc = &parseBool

	}

	result, err := service.Service.Find(searchDto)
	if err != nil {
		bytes, _ := json.Marshal(searchDto)
		publishers.Logger().Error(
			uuid.New().String(),
			"",
			fmt.Sprintf("problem during searching of measurements for searchdto: [%s]", string(bytes)),
			err.Error())
		return c.JSON(http.StatusInternalServerError, resterror.New(err.Error()))
	}

	if len(result) == 0 {
		return c.JSON(http.StatusNoContent, result)
	}

	mapResult := make(map[string][]dto.MeasurementDTO)

	for _, item := range result {
		if dtos, ok := mapResult[item.Uuid]; !ok {
			dtos = []dto.MeasurementDTO{}
			dtos = append(dtos, item)
			mapResult[item.Uuid] = dtos
		} else {
			dtos = append(dtos, item)
			mapResult[item.Uuid] = dtos
		}

	}

	return c.JSON(http.StatusOK, mapResult)
}
