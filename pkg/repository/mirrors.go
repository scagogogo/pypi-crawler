package repository

// ------------------------------------------------- --------------------------------------------------------------------

const NetEaseRepositoryServerURL = "https://mirrors.163.com/pypi"

// NewNetEaseRepository 网易
func NewNetEaseRepository() *Repository {
	return NewRepository(NewOptions().SetServerURL(NetEaseRepositoryServerURL))
}

// ------------------------------------------------- --------------------------------------------------------------------

const TencentCloudRepositoryServerURL = "https://mirrors.cloud.tencent.com/pypi"

// NewTencentCloudRepository 腾讯云
func NewTencentCloudRepository() *Repository {
	return NewRepository(NewOptions().SetServerURL(TencentCloudRepositoryServerURL))
}

// ------------------------------------------------- --------------------------------------------------------------------

//const AliCloudRepositoryServerURL = "http://mirrors.aliyun.com/pypi"
//
//// NewAliCloudRepository 中国科学技术大学
//func NewAliCloudRepository() *Repository {
//	return NewRepository(NewOptions().SetServerURL(AliCloudRepositoryServerURL))
//}

// ------------------------------------------------- --------------------------------------------------------------------

const UstcRepositoryServerURL = "https://pypi.mirrors.ustc.edu.cn"

// NewUstcRepository 中国科学技术大学
func NewUstcRepository() *Repository {
	return NewRepository(NewOptions().SetServerURL(UstcRepositoryServerURL))
}

// ------------------------------------------------- --------------------------------------------------------------------

const DouBanRepositoryServerURL = "http://pypi.douban.com"

// NewDouBanRepository 豆瓣镜像
func NewDouBanRepository() *Repository {
	return NewRepository(NewOptions().SetServerURL(DouBanRepositoryServerURL))
}

// ------------------------------------------------- --------------------------------------------------------------------

const SingHuaRepositoryServerURL = "https://pypi.tuna.tsinghua.edu.cn"

// NewTSingHuaRepository 清华镜像仓库
func NewTSingHuaRepository() *Repository {
	return NewRepository(NewOptions().SetServerURL(SingHuaRepositoryServerURL))
}

// ------------------------------------------------- --------------------------------------------------------------------
