package model

// Package 表示pypi中的一个包的信息
type Package struct {

	// 包的相关信息
	Information *Information `json:"info"`

	//
	LastSerial int `json:"last_serial"`

	// 这个包发布的版本都有哪些
	Releases *Releases `json:"releases"`

	// 最新版本发布的包
	Urls ReleaseUrls `json:"urls"`

	// TODO 漏洞信息，还不知道长啥样
	Vulnerabilities []any `json:"vulnerabilities"`
}
