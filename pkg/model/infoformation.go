package model

// Information 包含Python包的详细元数据信息
//
// 这个结构体映射了PyPI JSON API返回的包的info字段，包含了包的作者、描述、依赖等信息
//
// 使用示例:
//
//	pkg, err := repo.GetPackage(context.Background(), "requests")
//	if err == nil {
//	    info := pkg.Information
//	    fmt.Printf("包名: %s\n", info.Name)
//	    fmt.Printf("版本: %s\n", info.Version)
//	    fmt.Printf("摘要: %s\n", info.Summary)
//	    fmt.Printf("作者: %s (%s)\n", info.Author, info.AuthorEmail)
//	    fmt.Printf("依赖项数量: %d\n", len(info.RequiresDist))
//	}
type Information struct {
	// Author 包的作者名称
	Author string `json:"author"`

	// AuthorEmail 作者的电子邮箱地址
	AuthorEmail string `json:"author_email"`

	// BugtrackURL 项目的Bug追踪系统URL
	// 可能为nil或不同类型，因此使用interface{}
	BugtrackURL interface{} `json:"bugtrack_url"`

	// Classifiers 包的分类信息列表
	// 遵循PyPI的分类标准，如编程语言、许可证、发展状态等
	Classifiers []string `json:"classifiers"`

	// Description 包的详细描述
	// 通常是markdown或rst格式
	Description string `json:"description"`

	// DescriptionContentType 描述内容的MIME类型
	// 如 text/markdown、text/x-rst 等
	DescriptionContentType string `json:"description_content_type"`

	// DocsURL 文档URL
	// 可能为nil或不同类型，因此使用interface{}
	DocsURL interface{} `json:"docs_url"`

	// DownloadURL 包的下载URL
	DownloadURL string `json:"download_url"`

	// Downloads 包被下载的统计信息
	// 包含最近一天、一周和一个月的下载次数
	Downloads *Download `json:"downloads"`

	// HomePage 项目的主页URL
	HomePage string `json:"home_page"`

	// Keywords 关键字，用于搜索
	Keywords string `json:"keywords"`

	// License 包的许可证
	License string `json:"license"`

	// Maintainer 包的维护者名称
	Maintainer string `json:"maintainer"`

	// MaintainerEmail 维护者的电子邮箱
	MaintainerEmail string `json:"maintainer_email"`

	// Name 包名
	Name string `json:"name"`

	// PackageURL 包在PyPI上的URL
	PackageURL string `json:"package_url"`

	// Platform 支持的平台信息
	// 可能为nil或不同类型，因此使用interface{}
	Platform interface{} `json:"platform"`

	// ProjectURL 项目URL（可能与HomePage相同）
	ProjectURL string `json:"project_url"`

	// ProjectUrls 项目相关URL的映射
	// 如文档、源码、问题追踪等
	ProjectUrls *ProjectUrls `json:"project_urls"`

	// ReleaseURL 发布URL
	ReleaseURL string `json:"release_url"`

	// RequiresDist 当前包依赖的其他包列表
	// 使用PEP 508规范的依赖说明符格式
	RequiresDist []string `json:"requires_dist"`

	// RequiresPython 当前包要求的Python版本范围
	// 使用PEP 440规范的版本说明符格式，如">=3.6,<4.0"
	RequiresPython string `json:"requires_python"`

	// Summary 包的简短摘要描述
	Summary string `json:"summary"`

	// Version 包的当前版本
	// 遵循PEP 440的版本规范
	Version string `json:"version"`

	// Yanked 指示包是否已被撤回
	// 被撤回的包通常是因为发现严重bug或安全问题
	Yanked bool `json:"yanked"`

	// YankedReason 包被撤回的原因
	// 如果Yanked为true，这里可能包含撤回原因
	YankedReason interface{} `json:"yanked_reason"`
}

// Download 包含包的下载统计信息
//
// 记录了包在不同时间段内的下载次数
type Download struct {
	// LastDay 最近一天的下载次数
	LastDay int `json:"last_day"`

	// LastMonth 最近一个月的下载次数
	LastMonth int `json:"last_month"`

	// LastWeek 最近一周的下载次数
	LastWeek int `json:"last_week"`
}

// ProjectUrls 包含项目相关的URL链接集合
//
// 提供了项目各方面资源的链接，如文档、源代码等
type ProjectUrls struct {
	// Documentation 文档URL
	Documentation string `json:"Documentation"`

	// Homepage 项目主页URL
	Homepage string `json:"Homepage"`

	// Source 源代码仓库URL
	Source string `json:"Source"`
}
