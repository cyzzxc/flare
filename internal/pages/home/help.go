package home

import (
	"html/template"

	FlareDefine "github.com/soulteary/flare/config/define"
	FlareModel "github.com/soulteary/flare/config/model"
	FlareMDI "github.com/soulteary/flare/internal/resources/mdi"
)

func GenerateHelpTemplate() template.HTML {
	apps := []FlareModel.Bookmark{}
	apps = append(apps, []FlareModel.Bookmark{
		{
			Name: "程序首页",
			URL:  FlareDefine.RegularPages.Home.Path,
			Icon: "home-line",
			Desc: "",
		},
		{
			Name: "LLM 指南",
			URL:  "/llm.txt",
			Icon: "robot-line",
			Desc: "给 Agent 的 API 说明",
		},
		{
			Name: "程序设置",
			URL:  FlareDefine.RegularPages.Settings.Path,
			Icon: "settings-3-line",
			Desc: "",
		},
	}...)

	if FlareDefine.AppFlags.EnableGuide {
		apps = append(apps, FlareModel.Bookmark{
			Name: "向导页面",
			URL:  FlareDefine.RegularPages.Guide.Path,
			Icon: "compass-3-line",
			Desc: "",
		})
	}

	if FlareDefine.AppFlags.EnableEditor {
		apps = append(apps, FlareModel.Bookmark{
			Name: "内容编辑",
			URL:  FlareDefine.RegularPages.Editor.Path,
			Icon: "edit-circle-line",
			Desc: "",
		})
	}

	apps = append(apps, []FlareModel.Bookmark{
		{
			Name: "图标挑选",
			URL:  FlareDefine.RegularPages.Icons.Path,
			Icon: "heart-line",
			Desc: "",
		},
		{
			Name: "主题设置",
			URL:  FlareDefine.SettingPages.Theme.Path,
			Icon: "star-line",
			Desc: "",
		},
		{
			Name: "天气设置",
			URL:  FlareDefine.SettingPages.Weather.Path,
			Icon: "leaf-line",
			Desc: "",
		},
		{
			Name: "搜索设置",
			URL:  FlareDefine.SettingPages.Search.Path,
			Icon: "flashlight-line",
			Desc: "",
		},
		{
			Name: "界面设置",
			URL:  FlareDefine.SettingPages.Appearance.Path,
			Icon: "palette-line",
			Desc: "",
		},
		{
			Name: "程序版本",
			URL:  FlareDefine.SettingPages.Others.Path,
			Icon: "information-line",
			Desc: "",
		},
		{
			Name: "问题反馈",
			URL:  "https://github.com/soulteary/docker-flare/issues",
			Icon: "bug-line",
			Desc: "GitHub Issues",
		},
	}...)

	tpl := ""

	for _, app := range apps {

		desc := ""
		if app.Desc == "" {
			desc = app.URL
		} else {
			desc = app.Desc
		}

		tpl = tpl + `
			<div class="app-container" data-id="` + app.Icon + `">
			<a href="` + app.URL + `" class="app-item" title="` + app.Name + `">
			  <div class="app-icon">` + FlareMDI.GetIconByName(app.Icon) + `</div>
			  <div class="app-text">
				<p class="app-title">` + app.Name + `</p>
				<p class="app-desc">` + desc + `</p>
			  </div>
			</a>
			</div>
			`
	}
	return template.HTML(tpl)
}
