package models

// Package 表示从PyPI获取的包信息
// 包含包的基本元数据、发布版本信息和漏洞信息
type Package struct {
	// Info 包含包的基本元数据
	Info *PackageInfo `json:"info"`

	// LastSerial 包的最后序列号
	LastSerial int `json:"last_serial"`

	// Releases 包含所有发布版本及其文件信息
	// 键为版本字符串，值为该版本的发布文件列表
	Releases map[string][]*ReleaseFile `json:"releases"`

	// Urls 包含最新版本的发布文件信息
	Urls []*ReleaseFile `json:"urls"`

	// Vulnerabilities 已知的安全漏洞信息
	Vulnerabilities []Vulnerability `json:"vulnerabilities"`
}

// PackageInfo 包含包的详细元数据
type PackageInfo struct {
	// Name 包名
	Name string `json:"name"`

	// Version 当前版本
	Version string `json:"version"`

	// Summary 包的简短描述
	Summary string `json:"summary"`

	// Description 包的详细描述
	Description string `json:"description"`

	// DescriptionContentType 描述内容的MIME类型
	DescriptionContentType string `json:"description_content_type"`

	// Author 作者信息
	Author string `json:"author"`

	// AuthorEmail 作者的电子邮箱
	AuthorEmail string `json:"author_email"`

	// Maintainer 维护者信息
	Maintainer string `json:"maintainer"`

	// MaintainerEmail 维护者的电子邮箱
	MaintainerEmail string `json:"maintainer_email"`

	// License 包的许可证
	License string `json:"license"`

	// Keywords 关键字
	Keywords string `json:"keywords"`

	// ClassifiersArray 包的分类标签列表
	ClassifiersArray []string `json:"classifiers"`

	// ProjectURL 项目URL
	ProjectURL string `json:"project_url"`

	// ProjectURLs 项目相关URL，如文档、源码等
	ProjectURLs map[string]string `json:"project_urls"`

	// RequiresDist 依赖的其他包
	RequiresDist []string `json:"requires_dist"`

	// RequiresPython 需要的Python版本
	RequiresPython string `json:"requires_python"`

	// HomePage 主页URL
	HomePage string `json:"home_page"`

	// DocsURL 文档URL
	DocsURL string `json:"docs_url"`

	// DownloadURL 下载URL
	DownloadURL string `json:"download_url"`

	// Yanked 包是否被撤回
	Yanked bool `json:"yanked"`

	// YankedReason 包被撤回的原因
	YankedReason string `json:"yanked_reason,omitempty"`
}

// GetAllDependencies 返回包的所有依赖列表
// 这是一个便捷方法，可以直接获取依赖而不需要检查nil
func (p *PackageInfo) GetAllDependencies() []string {
	if p.RequiresDist == nil {
		return []string{}
	}
	return p.RequiresDist
}

// HasPythonRequirement 检查包是否具有Python版本要求
func (p *PackageInfo) HasPythonRequirement() bool {
	return p.RequiresPython != ""
}

// IsYanked 检查包是否被撤回
func (p *PackageInfo) IsYanked() bool {
	return p.Yanked
}

// GetProjectURLs 获取项目相关URL
// 返回一个新的map，避免nil检查
func (p *PackageInfo) GetProjectURLs() map[string]string {
	if p.ProjectURLs == nil {
		return map[string]string{}
	}
	return p.ProjectURLs
}
