package model

// Information 包的相关信息
type Information struct {

	// 包的作者
	Author string `json:"author"`

	// 作者的邮箱
	AuthorEmail string `json:"author_email"`

	// TODO
	BugtrackURL interface{} `json:"bugtrack_url"`

	// 包的分类信息
	Classifiers []string `json:"classifiers"`

	// 包的描述信息
	Description string `json:"description"`

	// 上面那个描述信息的格式，比如html还是text之类的
	DescriptionContentType string `json:"description_content_type"`

	// TODO
	DocsURL interface{} `json:"docs_url"`

	// 下载地址
	DownloadURL string `json:"download_url"`

	// 包被下载的次数
	Downloads       *Download `json:"downloads"`
	HomePage        string    `json:"home_page"`
	Keywords        string    `json:"keywords"`
	License         string    `json:"license"`
	Maintainer      string    `json:"maintainer"`
	MaintainerEmail string    `json:"maintainer_email"`
	Name            string    `json:"name"`
	PackageURL      string    `json:"package_url"`

	// TODO
	Platform interface{} `json:"platform"`

	ProjectURL  string       `json:"project_url"`
	ProjectUrls *ProjectUrls `json:"project_urls"`
	ReleaseURL  string       `json:"release_url"`

	// 当前包依赖的包
	RequiresDist []string `json:"requires_dist"`

	// 当前包要求的python版本
	RequiresPython string `json:"requires_python"`

	Summary string `json:"summary"`

	// 包的最新版本
	Version string `json:"version"`

	Yanked bool `json:"yanked"`

	// TODO
	YankedReason interface{} `json:"yanked_reason"`
}

type Download struct {
	LastDay   int `json:"last_day"`
	LastMonth int `json:"last_month"`
	LastWeek  int `json:"last_week"`
}

type ProjectUrls struct {
	Documentation string `json:"Documentation"`
	Homepage      string `json:"Homepage"`
	Source        string `json:"Source"`
}
