package repository

// Options 配置PyPI仓库的选项
//
// 该结构体包含与仓库连接相关的配置选项，如服务器URL和代理设置
type Options struct {
	// ServerURL 仓库服务器的地址
	// 默认为 "https://pypi.org"，可以修改为其他PyPI镜像地址
	ServerURL string

	// Proxy 请求时使用的代理IP
	// 格式如："http://127.0.0.1:8080" 或 "socks5://127.0.0.1:1080"
	// 如果为空，则不使用代理
	Proxy string
}

// DefaultServerURL 默认的PyPI官方仓库地址
// 当未指定ServerURL时，将使用此值
const DefaultServerURL = "https://pypi.org"

// NewOptions 创建一个新的仓库选项实例
//
// 返回值:
//   - *Options: 初始化的选项实例，使用默认的PyPI官方仓库地址
//
// 使用示例:
//
//	options := repository.NewOptions()
//	// 此时 options.ServerURL 的值为 "https://pypi.org"
//	// 且 options.Proxy 的值为空字符串
func NewOptions() *Options {
	return &Options{
		ServerURL: DefaultServerURL,
	}
}

// SetServerURL 设置仓库服务器的URL地址
//
// 参数:
//   - serverURL: PyPI仓库服务器地址
//
// 返回值:
//   - *Options: 更新后的选项实例，用于链式调用
//
// 使用示例:
//
//	options := repository.NewOptions().SetServerURL("https://mirrors.aliyun.com/pypi")
//	// 使用阿里云镜像作为仓库地址
func (x *Options) SetServerURL(serverURL string) *Options {
	x.ServerURL = serverURL
	return x
}

// SetProxy 设置请求使用的代理
//
// 参数:
//   - proxy: 代理服务器地址
//
// 返回值:
//   - *Options: 更新后的选项实例，用于链式调用
//
// 使用示例:
//
//	options := repository.NewOptions().SetProxy("http://127.0.0.1:8080")
//	// 使用本地的HTTP代理
//
//	options := repository.NewOptions().SetProxy("socks5://127.0.0.1:1080")
//	// 使用本地的SOCKS5代理
func (x *Options) SetProxy(proxy string) *Options {
	x.Proxy = proxy
	return x
}
