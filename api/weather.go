package api

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/LollipopKit/custed-server/consts"
	"github.com/LollipopKit/custed-server/logger"
	"github.com/LollipopKit/custed-server/model"
	"github.com/labstack/echo/v4"
)

var (
	data model.Weather
)

func init() {
	go func() {
		for {
			resp, err := http.Get(consts.WeatherAPIUrl)
			if err != nil {
				logger.E("api.weather get weather err: %s", err.Error())
				time.Sleep(time.Minute)
				continue
			}
			res, err := ioutil.ReadAll(resp.Body)

			err = json.Unmarshal(res, &data)
			if err != nil {
				logger.E("api.weather json unmarshal: %s", err.Error())
			}

			resp.Body.Close()

			time.Sleep(time.Minute * 20)
		}
	}()
}

func WeatherApi(c echo.Context) error {
	return okWithData(c, data.Data)
}
