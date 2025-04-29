package api

import (
	"context"

	"github.com/scagogogo/pypi-crawler/pkg/pypi/models"
)

// PyPIClient 定义了与PyPI API交互的客户端接口
// 此接口定义了所有可用的操作，使调用方能够更容易地进行测试和模拟
type PyPIClient interface {
	// GetPackageInfo 获取指定包的最新版本信息
	// 参数:
	//   - ctx: 上下文，用于控制请求的生命周期
	//   - packageName: 要获取信息的包名
	//
	// 返回值:
	//   - *models.Package: 包含包详细信息的结构体指针
	//   - error: 如有错误则返回，否则为nil
	GetPackageInfo(ctx context.Context, packageName string) (*models.Package, error)

	// GetPackageVersion 获取指定包的特定版本信息
	// 参数:
	//   - ctx: 上下文，用于控制请求的生命周期
	//   - packageName: 包名
	//   - version: 版本号
	//
	// 返回值:
	//   - *models.Package: 包含特定版本详细信息的结构体指针
	//   - error: 如有错误则返回，否则为nil
	GetPackageVersion(ctx context.Context, packageName string, version string) (*models.Package, error)

	// GetPackageReleases 获取指定包的所有发布版本
	// 参数:
	//   - ctx: 上下文，用于控制请求的生命周期
	//   - packageName: 包名
	//
	// 返回值:
	//   - []string: 包含版本号的字符串切片
	//   - error: 如有错误则返回，否则为nil
	GetPackageReleases(ctx context.Context, packageName string) ([]string, error)

	// CheckPackageVulnerabilities 检查指定包和版本是否存在已知漏洞
	// 参数:
	//   - ctx: 上下文，用于控制请求的生命周期
	//   - packageName: 包名
	//   - version: 版本号
	//
	// 返回值:
	//   - []models.Vulnerability: 包含漏洞信息的切片
	//   - error: 如有错误则返回，否则为nil
	CheckPackageVulnerabilities(ctx context.Context, packageName string, version string) ([]models.Vulnerability, error)

	// GetAllPackages 获取PyPI仓库中所有包的列表
	// 该方法通过调用PyPI的Simple API获取所有可用包的索引列表
	//
	// 参数:
	//   - ctx: 上下文，用于控制请求的生命周期
	//
	// 返回值:
	//   - []string: 包含所有包名的切片
	//   - error: 如有错误则返回，否则为nil
	GetAllPackages(ctx context.Context) ([]string, error)

	// GetPackageList 获取PyPI仓库中所有包的列表（以map形式返回）
	// 该方法是GetAllPackages的变种，返回map格式便于查询和遍历
	//
	// 参数:
	//   - ctx: 上下文，用于控制请求的生命周期
	//
	// 返回值:
	//   - map[string]struct{}: 包含所有包名的map，值为空结构体
	//   - error: 如有错误则返回，否则为nil
	GetPackageList(ctx context.Context) (map[string]struct{}, error)

	// SearchPackages 根据关键词搜索包
	// 该方法通过关键词搜索PyPI仓库中的包
	//
	// 参数:
	//   - ctx: 上下文，用于控制请求的生命周期
	//   - keyword: 搜索关键词
	//   - limit: 最大返回结果数，如果为0则使用默认值(100)
	//
	// 返回值:
	//   - []string: 匹配的包名列表
	//   - error: 如有错误则返回，否则为nil
	SearchPackages(ctx context.Context, keyword string, limit int) ([]string, error)
}
