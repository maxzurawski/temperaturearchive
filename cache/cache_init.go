package cache

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/xdevices/temperaturearchive/config"
	"github.com/xdevices/temperaturearchive/publishers"
)

func InitSensorsCache() error {
	proxy := config.TemperaturearchiveConfig().ProxyService()
	url := fmt.Sprintf("%s/api/register/cachesensors/", proxy)
	response, err := http.Get(url)
	if err != nil {
		publishers.Logger().Error(
			uuid.New().String(),
			"",
			fmt.Sprintf("could not obtain sensors"),
			err.Error())
		return err
	}

	var sensors []Sensor
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		publishers.Logger().Error(
			uuid.New().String(),
			"",
			fmt.Sprintf("could not read response body. message: [%s]", string(body)),
			err.Error())
		return err
	}

	if response.StatusCode == http.StatusNoContent {
		publishers.Logger().Info(
			uuid.New().String(),
			"",
			fmt.Sprintf("no sensors registered yet. sensors cache is empty"))
		SensorsCache = &Cache{}
		return nil
	}

	err = json.Unmarshal(body, &sensors)
	if err != nil {
		publishers.Logger().Error(
			uuid.New().String(),
			"",
			fmt.Sprintf("could not decode cached sensors"),
			err.Error())
		return err
	}

	publishers.Logger().Info(
		uuid.New().String(),
		"",
		fmt.Sprintf("resetting sensors cache"))
	SensorsCache = &Cache{}
	SensorsCache.FillCache(sensors)
	return nil
}
