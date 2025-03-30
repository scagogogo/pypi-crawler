package model

// Package 表示PyPI中的一个包的完整信息
//
// 这个结构体是对PyPI JSON API返回的包信息的映射，包含了包的基本信息、各版本信息及下载链接等
//
// 使用示例:
//
//	repo := repository.NewRepository()
//	pkg, err := repo.GetPackage(context.Background(), "requests")
//	if err != nil {
//	    log.Fatalf("获取包信息失败: %v", err)
//	}
//	// 访问包的基本信息
//	fmt.Printf("包名: %s\n", pkg.Information.Name)
//	fmt.Printf("最新版本: %s\n", pkg.Information.Version)
//
//	// 访问包的历史版本
//	for _, version := range pkg.Releases.VersionOrders {
//	    fmt.Printf("版本: %s\n", version)
//	}
type Package struct {
	// Information 包的相关基本信息
	// 包含作者、描述、许可证、依赖等信息
	Information *Information `json:"info"`

	// LastSerial 包的最后序列号
	// 用于标识包更新的序列
	LastSerial int `json:"last_serial"`

	// Releases 这个包发布的所有版本信息
	// 以版本号为键，对应版本的发布信息为值的映射
	Releases *Releases `json:"releases"`

	// Urls 最新版本发布的包的下载URL列表
	// 包含不同格式(如wheel、源码包等)的下载链接
	Urls ReleaseUrls `json:"urls"`

	// Vulnerabilities 包的安全漏洞信息
	// 如果包存在已知安全漏洞，这里会包含相关信息
	// 目前API返回格式尚不明确
	Vulnerabilities []any `json:"vulnerabilities"`
}
