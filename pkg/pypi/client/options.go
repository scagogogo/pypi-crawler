package client

import "time"

// Options 配置PyPI客户端的选项
// 包含API基址、超时设置、代理等配置
type Options struct {
	// BaseURL PyPI API的基础URL
	// 默认为官方PyPI: "https://pypi.org"
	BaseURL string

	// Timeout HTTP请求超时时间
	// 默认为30秒
	Timeout time.Duration

	// Proxy HTTP代理地址
	// 格式如: "http://127.0.0.1:8080" 或 "socks5://127.0.0.1:1080"
	// 为空时不使用代理
	Proxy string

	// UserAgent HTTP请求的User-Agent头
	// 遵循PyPI API指南，应包含唯一标识和联系信息
	UserAgent string

	// MaxRetries 请求失败后的最大重试次数
	// 默认为3次
	MaxRetries int

	// RetryDelay 重试之间的延迟时间
	// 默认为1秒
	RetryDelay time.Duration

	// RespectETag 是否遵循ETag缓存机制
	// 默认为true
	RespectETag bool
}

// 默认值常量
const (
	// DefaultBaseURL 默认的PyPI API基址
	DefaultBaseURL = "https://pypi.org"

	// DefaultUserAgent 默认的用户代理字符串
	DefaultUserAgent = "PyPIClient/2.0 (github.com/scagogogo/pypi-crawler)"

	// DefaultTimeout 默认的HTTP请求超时时间
	DefaultTimeout = 30 * time.Second

	// DefaultMaxRetries 默认的最大重试次数
	DefaultMaxRetries = 3

	// DefaultRetryDelay 默认的重试间隔时间
	DefaultRetryDelay = 1 * time.Second
)

// NewOptions 创建一个新的客户端选项实例，使用默认值
//
// 返回值:
//   - *Options: 初始化的选项实例
//
// 使用示例:
//
//	options := client.NewOptions()
//	// 默认使用官方PyPI，30秒超时，3次重试
func NewOptions() *Options {
	return &Options{
		BaseURL:     DefaultBaseURL,
		Timeout:     DefaultTimeout,
		UserAgent:   DefaultUserAgent,
		MaxRetries:  DefaultMaxRetries,
		RetryDelay:  DefaultRetryDelay,
		RespectETag: true,
	}
}

// WithBaseURL 设置API基础URL
//
// 参数:
//   - baseURL: API服务的基础URL
//
// 返回值:
//   - *Options: 更新后的选项实例，用于链式调用
//
// 使用示例:
//
//	options := client.NewOptions().WithBaseURL("https://pypi.tuna.tsinghua.edu.cn")
//	// 使用清华大学镜像源
func (o *Options) WithBaseURL(baseURL string) *Options {
	o.BaseURL = baseURL
	return o
}

// WithProxy 设置HTTP代理
//
// 参数:
//   - proxy: 代理服务器地址
//
// 返回值:
//   - *Options: 更新后的选项实例，用于链式调用
//
// 使用示例:
//
//	options := client.NewOptions().WithProxy("http://127.0.0.1:8080")
//	// 使用本地HTTP代理
func (o *Options) WithProxy(proxy string) *Options {
	o.Proxy = proxy
	return o
}

// WithTimeout 设置HTTP请求超时时间
//
// 参数:
//   - timeout: HTTP请求超时时间
//
// 返回值:
//   - *Options: 更新后的选项实例，用于链式调用
//
// 使用示例:
//
//	options := client.NewOptions().WithTimeout(10 * time.Second)
//	// 设置10秒超时
func (o *Options) WithTimeout(timeout time.Duration) *Options {
	o.Timeout = timeout
	return o
}

// WithUserAgent 设置请求的User-Agent头部
//
// 参数:
//   - userAgent: 用户代理字符串
//
// 返回值:
//   - *Options: 更新后的选项实例，用于链式调用
//
// 使用示例:
//
//	options := client.NewOptions().WithUserAgent("MyApp/1.0 (contact@example.com)")
//	// 设置自定义User-Agent
func (o *Options) WithUserAgent(userAgent string) *Options {
	o.UserAgent = userAgent
	return o
}

// WithMaxRetries 设置最大重试次数
//
// 参数:
//   - maxRetries: 最大重试次数
//
// 返回值:
//   - *Options: 更新后的选项实例，用于链式调用
//
// 使用示例:
//
//	options := client.NewOptions().WithMaxRetries(5)
//	// 设置失败后最多重试5次
func (o *Options) WithMaxRetries(maxRetries int) *Options {
	o.MaxRetries = maxRetries
	return o
}

// WithRetryDelay 设置重试间隔时间
//
// 参数:
//   - retryDelay: 重试之间的等待时间
//
// 返回值:
//   - *Options: 更新后的选项实例，用于链式调用
//
// 使用示例:
//
//	options := client.NewOptions().WithRetryDelay(2 * time.Second)
//	// 设置失败后等待2秒后重试
func (o *Options) WithRetryDelay(retryDelay time.Duration) *Options {
	o.RetryDelay = retryDelay
	return o
}

// WithRespectETag 设置是否遵循ETag缓存机制
//
// 参数:
//   - respectETag: 是否遵循ETag
//
// 返回值:
//   - *Options: 更新后的选项实例，用于链式调用
//
// 使用示例:
//
//	options := client.NewOptions().WithRespectETag(false)
//	// 禁用ETag缓存机制
func (o *Options) WithRespectETag(respectETag bool) *Options {
	o.RespectETag = respectETag
	return o
}
