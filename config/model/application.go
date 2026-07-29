package FlareModel

// Application Data Model
type Application struct {
	Title                   string `yaml:"Title" json:"Title"`
	Footer                  string `yaml:"Footer" json:"Footer"`
	OpenAppNewTab           bool   `yaml:"OpenAppNewTab" json:"OpenAppNewTab"`
	OpenBookmarkNewTab      bool   `yaml:"OpenBookmarkNewTab" json:"OpenBookmarkNewTab"`
	ShowTitle               bool   `yaml:"ShowTitle" json:"ShowTitle"`
	Greetings               string `yaml:"Greetings" json:"Greetings"`
	ShowSearchComponent     bool   `yaml:"ShowSearchComponent" json:"ShowSearchComponent"`
	DisabledSearchAutoFocus bool   `yaml:"DisabledSearchAutoFocus" json:"DisabledSearchAutoFocus"`
	ShowDateTime            bool   `yaml:"ShowDateTime" json:"ShowDateTime"`
	ShowApps                bool   `yaml:"ShowApps" json:"ShowApps"`
	ShowBookmarks           bool   `yaml:"ShowBookmarks" json:"ShowBookmarks"`
	HideSettingsButton      bool   `yaml:"HideSettingButton" json:"HideSettingButton"`
	HideHelpButton          bool   `yaml:"HideHelpButton" json:"HideHelpButton"`
	Theme                   string `yaml:"Theme" json:"Theme"`
	ShowWeather             bool   `yaml:"ShowWeather" json:"ShowWeather"`
	Location                string `yaml:"Location" json:"Location"`
	EnableEncryptedLink     bool   `yaml:"EnableEncryptedLink" json:"EnableEncryptedLink"`
	IconMode                string `yaml:"IconMode" json:"IconMode"`
	KeepLetterCase          bool   `yaml:"KeepLetterCase" json:"KeepLetterCase"`
}
