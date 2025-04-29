package mirrors

import (
	"github.com/scagogogo/pypi-crawler/pkg/pypi/api"
	"github.com/scagogogo/pypi-crawler/pkg/pypi/client"
)

// 预定义的PyPI镜像源URL常量
const (
	// OfficialURL PyPI官方源URL
	OfficialURL = "https://pypi.org"

	// TsinghuaURL 清华大学镜像源
	TsinghuaURL = "https://pypi.tuna.tsinghua.edu.cn"

	// DoubanURL 豆瓣镜像源
	DoubanURL = "https://pypi.doubanio.com"

	// AliyunURL 阿里云镜像源
	AliyunURL = "https://mirrors.aliyun.com/pypi"

	// TencentURL 腾讯云镜像源
	TencentURL = "https://mirrors.cloud.tencent.com/pypi"

	// UstcURL 中国科技大学镜像源
	UstcURL = "https://pypi.mirrors.ustc.edu.cn"

	// NeteaseURL 网易镜像源
	NeteaseURL = "https://mirrors.163.com/pypi"
)

// NewOfficialClient 创建使用官方源的客户端
//
// 参数:
//   - options: 可选配置选项，如代理、超时等
//
// 返回值:
//   - api.PyPIClient: 客户端接口实例
//
// 使用示例:
//
//	// 创建使用官方源的客户端
//	client := mirrors.NewOfficialClient()
func NewOfficialClient(options ...*client.Options) api.PyPIClient {
	var clientOptions *client.Options
	if len(options) > 0 {
		clientOptions = options[0]
	} else {
		clientOptions = client.NewOptions()
	}

	clientOptions.WithBaseURL(OfficialURL)
	return client.NewClient(clientOptions)
}

// NewTsinghuaClient 创建使用清华大学镜像源的客户端
//
// 参数:
//   - options: 可选配置选项，如代理、超时等
//
// 返回值:
//   - api.PyPIClient: 客户端接口实例
//
// 使用示例:
//
//	// 创建使用清华大学镜像源的客户端
//	client := mirrors.NewTsinghuaClient()
func NewTsinghuaClient(options ...*client.Options) api.PyPIClient {
	var clientOptions *client.Options
	if len(options) > 0 {
		clientOptions = options[0]
	} else {
		clientOptions = client.NewOptions()
	}

	clientOptions.WithBaseURL(TsinghuaURL)
	return client.NewClient(clientOptions)
}

// NewDoubanClient 创建使用豆瓣镜像源的客户端
//
// 参数:
//   - options: 可选配置选项，如代理、超时等
//
// 返回值:
//   - api.PyPIClient: 客户端接口实例
//
// 使用示例:
//
//	// 创建使用豆瓣镜像源的客户端
//	client := mirrors.NewDoubanClient()
func NewDoubanClient(options ...*client.Options) api.PyPIClient {
	var clientOptions *client.Options
	if len(options) > 0 {
		clientOptions = options[0]
	} else {
		clientOptions = client.NewOptions()
	}

	clientOptions.WithBaseURL(DoubanURL)
	return client.NewClient(clientOptions)
}

// NewAliyunClient 创建使用阿里云镜像源的客户端
//
// 参数:
//   - options: 可选配置选项，如代理、超时等
//
// 返回值:
//   - api.PyPIClient: 客户端接口实例
//
// 使用示例:
//
//	// 创建使用阿里云镜像源的客户端
//	client := mirrors.NewAliyunClient()
func NewAliyunClient(options ...*client.Options) api.PyPIClient {
	var clientOptions *client.Options
	if len(options) > 0 {
		clientOptions = options[0]
	} else {
		clientOptions = client.NewOptions()
	}

	clientOptions.WithBaseURL(AliyunURL)
	return client.NewClient(clientOptions)
}

// NewTencentClient 创建使用腾讯云镜像源的客户端
//
// 参数:
//   - options: 可选配置选项，如代理、超时等
//
// 返回值:
//   - api.PyPIClient: 客户端接口实例
//
// 使用示例:
//
//	// 创建使用腾讯云镜像源的客户端
//	client := mirrors.NewTencentClient()
func NewTencentClient(options ...*client.Options) api.PyPIClient {
	var clientOptions *client.Options
	if len(options) > 0 {
		clientOptions = options[0]
	} else {
		clientOptions = client.NewOptions()
	}

	clientOptions.WithBaseURL(TencentURL)
	return client.NewClient(clientOptions)
}

// NewUstcClient 创建使用中国科技大学镜像源的客户端
//
// 参数:
//   - options: 可选配置选项，如代理、超时等
//
// 返回值:
//   - api.PyPIClient: 客户端接口实例
//
// 使用示例:
//
//	// 创建使用中国科技大学镜像源的客户端
//	client := mirrors.NewUstcClient()
func NewUstcClient(options ...*client.Options) api.PyPIClient {
	var clientOptions *client.Options
	if len(options) > 0 {
		clientOptions = options[0]
	} else {
		clientOptions = client.NewOptions()
	}

	clientOptions.WithBaseURL(UstcURL)
	return client.NewClient(clientOptions)
}

// NewNeteaseClient 创建使用网易镜像源的客户端
//
// 参数:
//   - options: 可选配置选项，如代理、超时等
//
// 返回值:
//   - api.PyPIClient: 客户端接口实例
//
// 使用示例:
//
//	// 创建使用网易镜像源的客户端
//	client := mirrors.NewNeteaseClient()
func NewNeteaseClient(options ...*client.Options) api.PyPIClient {
	var clientOptions *client.Options
	if len(options) > 0 {
		clientOptions = options[0]
	} else {
		clientOptions = client.NewOptions()
	}

	clientOptions.WithBaseURL(NeteaseURL)
	return client.NewClient(clientOptions)
}
