package repository

// ------------------------------------------------- --------------------------------------------------------------------

// NetEaseRepositoryServerURL 网易PyPI镜像仓库地址
// 网易提供的开源镜像站点，访问速度较快
const NetEaseRepositoryServerURL = "https://mirrors.163.com/pypi"

// NewNetEaseRepository 创建一个使用网易镜像的PyPI仓库实例
//
// 返回值:
//   - *Repository: 使用网易镜像的仓库实例
//
// 使用示例:
//
//	repo := repository.NewNetEaseRepository()
//	// 现在可以使用repo访问网易的PyPI镜像
//	indexes, err := repo.DownloadIndex(context.Background())
func NewNetEaseRepository() *Repository {
	return NewRepository(NewOptions().SetServerURL(NetEaseRepositoryServerURL))
}

// ------------------------------------------------- --------------------------------------------------------------------

// TencentCloudRepositoryServerURL 腾讯云PyPI镜像仓库地址
// 腾讯云提供的开源镜像站点，服务稳定
const TencentCloudRepositoryServerURL = "https://mirrors.cloud.tencent.com/pypi"

// NewTencentCloudRepository 创建一个使用腾讯云镜像的PyPI仓库实例
//
// 返回值:
//   - *Repository: 使用腾讯云镜像的仓库实例
//
// 使用示例:
//
//	repo := repository.NewTencentCloudRepository()
//	// 现在可以使用repo访问腾讯云的PyPI镜像
//	pkg, err := repo.GetPackage(context.Background(), "requests")
func NewTencentCloudRepository() *Repository {
	return NewRepository(NewOptions().SetServerURL(TencentCloudRepositoryServerURL))
}

// ------------------------------------------------- --------------------------------------------------------------------

// AliCloudRepositoryServerURL 阿里云PyPI镜像仓库地址（已注释，可能不可用）
//const AliCloudRepositoryServerURL = "http://mirrors.aliyun.com/pypi"
//
//// NewAliCloudRepository 创建一个使用阿里云镜像的PyPI仓库实例
//// 注意：此方法已被注释，因为相应的镜像可能不再可用
//func NewAliCloudRepository() *Repository {
//	return NewRepository(NewOptions().SetServerURL(AliCloudRepositoryServerURL))
//}

// ------------------------------------------------- --------------------------------------------------------------------

// UstcRepositoryServerURL 中国科学技术大学PyPI镜像仓库地址
// 中科大提供的开源镜像站点，稳定性较好
const UstcRepositoryServerURL = "https://pypi.mirrors.ustc.edu.cn"

// NewUstcRepository 创建一个使用中国科学技术大学镜像的PyPI仓库实例
//
// 返回值:
//   - *Repository: 使用中科大镜像的仓库实例
//
// 使用示例:
//
//	repo := repository.NewUstcRepository()
//	// 现在可以使用repo访问中科大的PyPI镜像
func NewUstcRepository() *Repository {
	return NewRepository(NewOptions().SetServerURL(UstcRepositoryServerURL))
}

// ------------------------------------------------- --------------------------------------------------------------------

// DouBanRepositoryServerURL 豆瓣PyPI镜像仓库地址
// 豆瓣提供的开源镜像站点
const DouBanRepositoryServerURL = "http://pypi.douban.com"

// NewDouBanRepository 创建一个使用豆瓣镜像的PyPI仓库实例
//
// 返回值:
//   - *Repository: 使用豆瓣镜像的仓库实例
//
// 使用示例:
//
//	repo := repository.NewDouBanRepository()
//	// 现在可以使用repo访问豆瓣的PyPI镜像
//
// 注意:
//
//	豆瓣镜像有时可能不稳定，如果出现连接问题，请考虑使用其他镜像源
func NewDouBanRepository() *Repository {
	return NewRepository(NewOptions().SetServerURL(DouBanRepositoryServerURL))
}

// ------------------------------------------------- --------------------------------------------------------------------

// SingHuaRepositoryServerURL 清华大学PyPI镜像仓库地址
// 清华大学TUNA协会提供的开源镜像站点，更新迅速，是国内较为推荐的PyPI镜像
const SingHuaRepositoryServerURL = "https://pypi.tuna.tsinghua.edu.cn"

// NewTSingHuaRepository 创建一个使用清华镜像仓库的PyPI仓库实例
//
// 返回值:
//   - *Repository: 使用清华镜像的仓库实例
//
// 使用示例:
//
//	repo := repository.NewTSingHuaRepository()
//	// 现在可以使用repo访问清华的PyPI镜像
//	// 这是国内较为推荐的PyPI镜像源之一
func NewTSingHuaRepository() *Repository {
	return NewRepository(NewOptions().SetServerURL(SingHuaRepositoryServerURL))
}

// ------------------------------------------------- --------------------------------------------------------------------
