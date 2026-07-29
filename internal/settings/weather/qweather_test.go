package weather

import (
	"os"
	"testing"
)

func TestMapQWeatherIcon(t *testing.T) {
	cases := []struct {
		icon string
		day  bool
		want string
	}{
		{"100", true, "day"},
		{"150", false, "night"},
		{"101", true, "cloudy-day"},
		{"151", false, "cloudy-night"},
		{"104", true, "cloudy"},
		{"300", true, "rainy-1"},
		{"350", false, "rainy-night-1"},
		{"302", true, "thunder"},
		{"305", true, "rainy-4"},
		{"306", true, "rainy-5"},
		{"307", true, "rainy-6"},
		{"400", true, "snowy-day"},
		{"401", true, "snowy-4"},
		{"404", false, "rain-snow-night"},
		{"500", true, "icon-fog"},
		{"503", true, "icon-dust"},
		{"bad", true, "day"},
		{"bad", false, "night"},
	}
	for _, c := range cases {
		got := mapQWeatherIcon(c.icon, c.day)
		if got != c.want {
			t.Fatalf("mapQWeatherIcon(%q, %v)=%q want %q", c.icon, c.day, got, c.want)
		}
	}
}

func TestIsLocationID(t *testing.T) {
	if !isLocationID("101010100") {
		t.Fatal("expected LocationID")
	}
	if isLocationID("北京") || isLocationID("116.41,39.92") {
		t.Fatal("should reject non-ID")
	}
}

func TestIsRainyDay(t *testing.T) {
	if !isRainyDay("305", "小雨", "150", "晴", "0.0") {
		t.Fatal("iconDay rain")
	}
	if !isRainyDay("100", "晴", "350", "阵雨", "0.0") {
		t.Fatal("iconNight rain")
	}
	if !isRainyDay("100", "晴转雨", "150", "晴", "0.0") {
		t.Fatal("text rain")
	}
	if !isRainyDay("100", "晴", "150", "晴", "1.2") {
		t.Fatal("precip rain")
	}
	if isRainyDay("100", "晴", "150", "晴", "0.0") {
		t.Fatal("should be dry")
	}
}

func TestFetchQWeatherNowLive(t *testing.T) {
	if os.Getenv("FLARE_QWEATHER_KEY") == "" || os.Getenv("FLARE_QWEATHER_HOST") == "" {
		t.Skip("no qweather credentials")
	}
	host, key, err := qweatherConfig()
	if err != nil {
		t.Fatal(err)
	}
	locID, err := resolveLocationID(host, key, "北京")
	if err != nil {
		t.Fatal(err)
	}
	temp, hum, icon, text, _, err := fetchQWeatherNow(host, key, locID)
	if err != nil {
		t.Fatal(err)
	}
	if temp == 0 && hum == 0 {
		t.Fatal("empty weather")
	}
	hint := tomorrowRainHint(host, key, locID)
	t.Logf("temp=%d humidity=%d icon=%s text=%s mapped=%s hint=%q", temp, hum, icon, text, mapQWeatherIcon(icon, true), hint)
}
