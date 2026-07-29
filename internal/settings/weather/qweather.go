package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	FlareData "github.com/soulteary/flare/config/data"
)

const (
	envQWeatherKey  = "FLARE_QWEATHER_KEY"
	envQWeatherHost = "FLARE_QWEATHER_HOST"
)

var qweatherClient = &http.Client{Timeout: 5 * time.Second}

type qweatherNowResponse struct {
	Code       string `json:"code"`
	UpdateTime string `json:"updateTime"`
	Now        struct {
		ObsTime   string `json:"obsTime"`
		Temp      string `json:"temp"`
		Icon      string `json:"icon"`
		Text      string `json:"text"`
		Humidity  string `json:"humidity"`
		FeelsLike string `json:"feelsLike"`
	} `json:"now"`
}

type qweatherGeoResponse struct {
	Code     string `json:"code"`
	Location []struct {
		Name string `json:"name"`
		ID   string `json:"id"`
		Adm1 string `json:"adm1"`
		Adm2 string `json:"adm2"`
	} `json:"location"`
}

func qweatherConfig() (host, key string, err error) {
	// config.yml first, then env
	key, host = FlareData.GetQWeatherConfig()
	key = strings.TrimSpace(key)
	host = strings.TrimSpace(host)
	if key == "" {
		key = strings.TrimSpace(os.Getenv(envQWeatherKey))
	}
	if host == "" {
		host = strings.TrimSpace(os.Getenv(envQWeatherHost))
	}
	host = strings.TrimPrefix(host, "https://")
	host = strings.TrimPrefix(host, "http://")
	host = strings.TrimSuffix(host, "/")
	if key == "" || host == "" {
		return "", "", errors.New("未配置和风天气 Key / API Host（设置页或 FLARE_QWEATHER_KEY / FLARE_QWEATHER_HOST）")
	}
	return host, key, nil
}

func qweatherGet(apiURL, key string, target interface{}) error {
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-QW-Api-Key", key)

	res, err := qweatherClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", res.StatusCode)
	}
	return json.NewDecoder(res.Body).Decode(target)
}

func resolveLocationID(host, key, location string) (string, error) {
	location = strings.TrimSpace(location)
	if location == "" {
		return "", errors.New("地区为空")
	}
	// LocationID or lon,lat can pass through
	if isLocationID(location) || strings.Contains(location, ",") {
		return location, nil
	}

	u := fmt.Sprintf("https://%s/geo/v2/city/lookup?location=%s&number=1&lang=zh",
		host, url.QueryEscape(location))
	var geo qweatherGeoResponse
	if err := qweatherGet(u, key, &geo); err != nil {
		return "", err
	}
	if geo.Code != "200" || len(geo.Location) == 0 {
		return "", fmt.Errorf("城市查询失败 code=%s", geo.Code)
	}
	return geo.Location[0].ID, nil
}

func isLocationID(s string) bool {
	if len(s) < 6 || len(s) > 12 {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

type qweatherDailyResponse struct {
	Code  string `json:"code"`
	Daily []struct {
		FxDate    string `json:"fxDate"`
		IconDay   string `json:"iconDay"`
		TextDay   string `json:"textDay"`
		IconNight string `json:"iconNight"`
		TextNight string `json:"textNight"`
		Precip    string `json:"precip"`
	} `json:"daily"`
}

func fetchQWeatherNow(host, key, locID string) (temp, humidity int, icon, text, updateTime string, err error) {
	u := fmt.Sprintf("https://%s/v7/weather/now?location=%s&lang=zh",
		host, url.QueryEscape(locID))
	var body qweatherNowResponse
	if err := qweatherGet(u, key, &body); err != nil {
		return 0, 0, "", "", "", err
	}
	if body.Code != "200" {
		return 0, 0, "", "", "", fmt.Errorf("天气查询失败 code=%s", body.Code)
	}

	temp, err = strconv.Atoi(body.Now.Temp)
	if err != nil {
		return 0, 0, "", "", "", err
	}
	humidity, err = strconv.Atoi(body.Now.Humidity)
	if err != nil {
		return 0, 0, "", "", "", err
	}

	updateTime = body.UpdateTime
	if updateTime == "" {
		updateTime = body.Now.ObsTime
	}
	return temp, humidity, body.Now.Icon, body.Now.Text, updateTime, nil
}

// tomorrowRainHint returns a short tip if tomorrow's forecast looks rainy.
func tomorrowRainHint(host, key, locID string) string {
	u := fmt.Sprintf("https://%s/v7/weather/3d?location=%s&lang=zh",
		host, url.QueryEscape(locID))
	var body qweatherDailyResponse
	if err := qweatherGet(u, key, &body); err != nil || body.Code != "200" {
		return ""
	}
	// daily[0]=today, daily[1]=tomorrow
	if len(body.Daily) < 2 {
		return ""
	}
	d := body.Daily[1]
	if isRainyDay(d.IconDay, d.TextDay, d.IconNight, d.TextNight, d.Precip) {
		return "明天有雨，记得带伞"
	}
	return ""
}

func isRainyDay(iconDay, textDay, iconNight, textNight, precip string) bool {
	if isRainIcon(iconDay) || isRainIcon(iconNight) {
		return true
	}
	if strings.Contains(textDay, "雨") || strings.Contains(textNight, "雨") {
		return true
	}
	if p, err := strconv.ParseFloat(precip, 64); err == nil && p > 0 {
		return true
	}
	return false
}

// 和风 300–399 = 雨/阵雨/雷雨等
func isRainIcon(icon string) bool {
	code, err := strconv.Atoi(icon)
	if err != nil {
		return false
	}
	return code >= 300 && code <= 399
}

// mapQWeatherIcon maps 和风 icon code → funny-china-weather SVG key.
func mapQWeatherIcon(icon string, isDay bool) string {
	code, err := strconv.Atoi(icon)
	if err != nil {
		if isDay {
			return "day"
		}
		return "night"
	}

	switch {
	case code == 100 || code == 150:
		if isDay {
			return "day"
		}
		return "night"
	case code >= 101 && code <= 103, code >= 151 && code <= 153:
		if isDay {
			return "cloudy-day"
		}
		return "cloudy-night"
	case code == 104:
		return "cloudy"
	case code == 300 || code == 301 || code == 350 || code == 351:
		if isDay {
			return "rainy-1"
		}
		return "rainy-night-1"
	case code == 302 || code == 303:
		return "thunder"
	case code == 304:
		return "thunder-rain"
	case code == 305 || code == 309 || code == 399:
		return "rainy-4"
	case code == 306 || code == 314 || code == 315:
		return "rainy-5"
	case code >= 307 && code <= 313, code >= 316 && code <= 318:
		return "rainy-6"
	case code == 400 || code == 407 || code == 456 || code == 457:
		if isDay {
			return "snowy-day"
		}
		return "snowy-night"
	case code == 401 || code == 408:
		return "snowy-4"
	case code == 402 || code == 409:
		return "snowy-5"
	case code == 403 || code == 410 || code == 499:
		return "snowy-6"
	case code >= 404 && code <= 406:
		if isDay {
			return "rain-snow-day"
		}
		return "rain-snow-night"
	case code >= 500 && code <= 502, code == 509 || code == 510 || code == 514 || code == 515:
		return "icon-fog"
	case code >= 503 && code <= 508, code >= 511 && code <= 513:
		return "icon-dust"
	default:
		if isDay {
			return "day"
		}
		return "night"
	}
}
