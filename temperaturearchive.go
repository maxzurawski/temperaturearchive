package main

import (
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/maxzurawski/temperaturearchive/cache"
	"github.com/maxzurawski/temperaturearchive/config"
	"github.com/maxzurawski/temperaturearchive/dbprovider"
	"github.com/maxzurawski/temperaturearchive/dto"
	"github.com/maxzurawski/temperaturearchive/handlers"
	"github.com/maxzurawski/temperaturearchive/observer"
	"github.com/maxzurawski/temperaturearchive/processors"
	"github.com/maxzurawski/temperaturearchive/publishers"
	"github.com/maxzurawski/temperaturearchive/service"
)

func main() {

	go observer.ObserveSensorChanges()
	go observer.TemperatureObserver([]func(dto.MeasurementDTO){processors.NotifierProcessor, processors.ArchiveProcessor})

	e := echo.New()
	e.GET("/", handlers.HandleFind)
	e.GET("/ws", handlers.HandleWebsocket)
	e.Logger.Fatal(e.Start(config.TemperaturearchiveConfig().Address()))
}

func init() {
	manager := config.EurekaManagerInit()
	manager.SendRegistrationOrFail()
	manager.ScheduleHeartBeat(config.TemperaturearchiveConfig().ServiceName(), 10)

	dbprovider.InitDbManager()
	service.Init()

	publishers.InitLogger()
	publishers.InitNotifier()

	_ = cache.InitSensorsCache(uuid.New().String())

}
