package weather

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	FlareData "github.com/soulteary/flare/config/data"
	FlareDefine "github.com/soulteary/flare/config/define"
	FlareModel "github.com/soulteary/flare/config/model"
	FlareAuth "github.com/soulteary/flare/internal/auth"
)

func GetWeatherInfo(location string) (response FlareModel.Weather, desc string, err error) {
	host, key, cfgErr := qweatherConfig()
	if cfgErr != nil {
		return response, cfgErr.Error(), errors.New("获取远程数据失败")
	}
	locID, locErr := resolveLocationID(host, key, location)
	if locErr != nil {
		return response, locErr.Error(), errors.New("获取远程数据失败")
	}

	degree, humidity, icon, text, lastUpdate, fetchErr := fetchQWeatherNow(host, key, locID)
	if fetchErr != nil {
		return response, fetchErr.Error(), errors.New("获取远程数据失败")
	}

	hour, _, _ := time.Now().Clock()
	isDay := hour >= 5 && hour <= 18

	const cacheSeconds = 60 * 10 // 10 minutes

	response.ExternalLastUpdate = lastUpdate
	response.Degree = degree
	response.IsDay = isDay
	response.ConditionCode = mapQWeatherIcon(icon, isDay)
	response.ConditionText = text
	response.Humidity = humidity
	response.UmbrellaHint = tomorrowRainHint(host, key, locID)
	response.Expires = time.Now().Unix() + cacheSeconds

	return response, "接口正常", nil
}

func RegisterRouting(router *gin.Engine) {
	router.GET(FlareDefine.SettingPages.Weather.Path, FlareAuth.AuthRequired, pageHome)
	if !FlareDefine.AppFlags.EnableOfflineMode {
		router.POST(FlareDefine.SettingPages.Weather.Path, FlareAuth.AuthRequired, updateWeatherOptions)
	}
}

func pageHome(c *gin.Context) {
	render(c)
}

func updateWeatherOptions(c *gin.Context) {
	type UpdateBody struct {
		Location     string `form:"location"`
		ShowWeather  bool   `form:"show"`
		QWeatherKey  string `form:"qweather_key"`
		QWeatherHost string `form:"qweather_host"`
	}

	var body UpdateBody
	if c.ShouldBind(&body) != nil {
		c.PureJSON(http.StatusForbidden, "提交数据缺失")
		return
	}

	FlareData.UpdateWeatherAndLocation(body.ShowWeather, body.Location, body.QWeatherKey, body.QWeatherHost)
	render(c)
}

func render(c *gin.Context) {
	location, showWeather := FlareData.GetLocationAndWeatherShow()
	key, host := FlareData.GetQWeatherConfig()
	options := FlareData.GetAllSettingsOptions()

	c.HTML(
		http.StatusOK,
		"settings.html",
		gin.H{
			"DebugMode":         FlareDefine.AppFlags.DebugMode,
			"PageInlineStyle":   FlareDefine.GetPageInlineStyle(),
			"ShowWeatherModule": !FlareDefine.AppFlags.EnableOfflineMode && showWeather,
			"ShowWeather":       showWeather,

			"PageName":       "Weather",
			"PageAppearance": FlareDefine.GetAppBodyStyle(),
			"SettingPages":   FlareDefine.SettingPages,
			"SettingsURI":    FlareDefine.RegularPages.Settings.Path,
			"OptionTitle":    options.Title,

			"Location":     location,
			"QWeatherKey":  key,
			"QWeatherHost": host,
		},
	)
}
