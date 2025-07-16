# 客户端配置

本文档详细介绍如何配置 PyPI Crawler 客户端以满足不同的使用需求。

## 📋 目录

- [配置选项概览](#配置选项概览)
- [创建配置](#创建配置)
- [配置选项详解](#配置选项详解)
- [常用配置场景](#常用配置场景)
- [最佳实践](#最佳实践)

## 配置选项概览

PyPI Crawler 提供了丰富的配置选项来自定义客户端行为：

```go
type Options struct {
    BaseURL     string        // API 基础 URL
    Timeout     time.Duration // HTTP 请求超时时间
    Proxy       string        // HTTP 代理地址
    UserAgent   string        // User-Agent 头部
    MaxRetries  int           // 最大重试次数
    RetryDelay  time.Duration // 重试间隔时间
    RespectETag bool          // 是否遵循 ETag 缓存
}
```

## 创建配置

### 使用默认配置

```go
import "github.com/scagogogo/pypi-crawler/pkg/pypi/client"

// 创建默认配置
options := client.NewOptions()

// 默认值:
// - BaseURL: "https://pypi.org"
// - Timeout: 30 秒
// - UserAgent: "PyPIClient/2.0 (github.com/scagogogo/pypi-crawler)"
// - MaxRetries: 3
// - RetryDelay: 1 秒
// - RespectETag: true
```

### 链式配置

```go
options := client.NewOptions().
    WithBaseURL("https://pypi.tuna.tsinghua.edu.cn").
    WithTimeout(15 * time.Second).
    WithMaxRetries(5).
    WithUserAgent("MyApp/1.0 (contact@example.com)")
```

## 配置选项详解

### BaseURL - API 基础地址

设置 PyPI API 的基础 URL。

```go
// 使用官方源
options := client.NewOptions().WithBaseURL("https://pypi.org")

// 使用清华镜像
options := client.NewOptions().WithBaseURL("https://pypi.tuna.tsinghua.edu.cn")

// 使用自定义源
options := client.NewOptions().WithBaseURL("https://my-pypi-mirror.com")
```

**注意事项:**
- URL 不应包含尾部斜杠
- 必须是有效的 HTTP/HTTPS URL
- 镜像源必须兼容 PyPI API

### Timeout - 请求超时

设置 HTTP 请求的超时时间。

```go
import "time"

// 设置 10 秒超时
options := client.NewOptions().WithTimeout(10 * time.Second)

// 设置 1 分钟超时（适用于获取所有包列表）
options := client.NewOptions().WithTimeout(60 * time.Second)

// 设置 5 分钟超时（适用于大批量操作）
options := client.NewOptions().WithTimeout(5 * time.Minute)
```

**建议值:**
- 普通查询: 10-30 秒
- 获取包列表: 60-120 秒
- 批量操作: 5-10 分钟

### Proxy - HTTP 代理

配置 HTTP 代理服务器。

```go
// HTTP 代理
options := client.NewOptions().WithProxy("http://proxy.company.com:8080")

// HTTPS 代理
options := client.NewOptions().WithProxy("https://proxy.company.com:8080")

// SOCKS5 代理
options := client.NewOptions().WithProxy("socks5://127.0.0.1:1080")

// 带认证的代理
options := client.NewOptions().WithProxy("http://username:password@proxy.com:8080")
```

**支持的代理类型:**
- HTTP 代理
- HTTPS 代理
- SOCKS5 代理

### UserAgent - 用户代理

设置 HTTP 请求的 User-Agent 头部。

```go
// 自定义 User-Agent
options := client.NewOptions().WithUserAgent("MyApp/1.0 (contact@example.com)")

// 包含版本信息
options := client.NewOptions().WithUserAgent("DataAnalyzer/2.1.0 (Python Package Scanner)")

// 遵循 PyPI 建议格式
options := client.NewOptions().WithUserAgent("CompanyTool/1.0 (admin@company.com)")
```

**PyPI 建议格式:**
```
ApplicationName/Version (contact-info)
```

**注意事项:**
- 应包含应用名称和版本
- 建议包含联系信息
- 避免使用通用的 User-Agent

### MaxRetries - 最大重试次数

设置请求失败后的最大重试次数。

```go
// 不重试
options := client.NewOptions().WithMaxRetries(0)

// 重试 3 次（默认）
options := client.NewOptions().WithMaxRetries(3)

// 重试 10 次（适用于不稳定网络）
options := client.NewOptions().WithMaxRetries(10)
```

**重试条件:**
- 网络连接失败
- HTTP 5xx 服务器错误
- 请求超时

**不重试的情况:**
- HTTP 4xx 客户端错误（如 404）
- 上下文取消

### RetryDelay - 重试间隔

设置重试之间的等待时间。

```go
// 立即重试
options := client.NewOptions().WithRetryDelay(0)

// 等待 1 秒后重试（默认）
options := client.NewOptions().WithRetryDelay(1 * time.Second)

// 等待 5 秒后重试
options := client.NewOptions().WithRetryDelay(5 * time.Second)
```

**建议值:**
- 快速重试: 0-1 秒
- 普通重试: 1-3 秒
- 保守重试: 5-10 秒

### RespectETag - ETag 缓存

控制是否遵循 HTTP ETag 缓存机制。

```go
// 启用 ETag 缓存（默认）
options := client.NewOptions().WithRespectETag(true)

// 禁用 ETag 缓存
options := client.NewOptions().WithRespectETag(false)
```

**ETag 缓存的作用:**
- 减少不必要的数据传输
- 提高响应速度
- 减轻服务器负载

## 常用配置场景

### 开发环境配置

```go
// 开发环境：快速响应，详细错误信息
devOptions := client.NewOptions().
    WithTimeout(10 * time.Second).
    WithMaxRetries(1).
    WithUserAgent("DevApp/0.1.0 (dev@company.com)")

client := mirrors.NewOfficialClient(devOptions)
```

### 生产环境配置

```go
// 生产环境：稳定可靠，适当重试
prodOptions := client.NewOptions().
    WithTimeout(30 * time.Second).
    WithMaxRetries(3).
    WithRetryDelay(2 * time.Second).
    WithUserAgent("ProdApp/1.0.0 (ops@company.com)")

client := mirrors.NewTsinghuaClient(prodOptions)
```

### 批量处理配置

```go
// 批量处理：长超时，更多重试
batchOptions := client.NewOptions().
    WithTimeout(5 * time.Minute).
    WithMaxRetries(5).
    WithRetryDelay(3 * time.Second).
    WithUserAgent("BatchProcessor/1.0 (batch@company.com)")

client := mirrors.NewAliyunClient(batchOptions)
```

### 网络受限环境配置

```go
// 网络受限：使用代理，增加重试
restrictedOptions := client.NewOptions().
    WithProxy("http://proxy.company.com:8080").
    WithTimeout(60 * time.Second).
    WithMaxRetries(10).
    WithRetryDelay(5 * time.Second).
    WithUserAgent("RestrictedApp/1.0 (admin@company.com)")

client := mirrors.NewOfficialClient(restrictedOptions)
```

### 高频访问配置

```go
// 高频访问：启用缓存，合理限制
highFreqOptions := client.NewOptions().
    WithTimeout(15 * time.Second).
    WithMaxRetries(2).
    WithRetryDelay(1 * time.Second).
    WithRespectETag(true).
    WithUserAgent("HighFreqApp/1.0 (api@company.com)")

client := mirrors.NewTsinghuaClient(highFreqOptions)
```

## 最佳实践

### 1. 选择合适的超时时间

```go
// 根据操作类型设置超时
var options *client.Options

switch operationType {
case "single_package":
    options = client.NewOptions().WithTimeout(10 * time.Second)
case "package_list":
    options = client.NewOptions().WithTimeout(60 * time.Second)
case "batch_operation":
    options = client.NewOptions().WithTimeout(5 * time.Minute)
}
```

### 2. 合理设置重试策略

```go
// 根据网络环境调整重试
networkQuality := getNetworkQuality() // 假设的函数

var maxRetries int
var retryDelay time.Duration

switch networkQuality {
case "excellent":
    maxRetries = 1
    retryDelay = 500 * time.Millisecond
case "good":
    maxRetries = 3
    retryDelay = 1 * time.Second
case "poor":
    maxRetries = 5
    retryDelay = 3 * time.Second
}

options := client.NewOptions().
    WithMaxRetries(maxRetries).
    WithRetryDelay(retryDelay)
```

### 3. 使用有意义的 User-Agent

```go
// 包含应用信息和联系方式
userAgent := fmt.Sprintf("%s/%s (%s)", 
    appName, appVersion, contactEmail)

options := client.NewOptions().WithUserAgent(userAgent)
```

### 4. 环境变量配置

```go
import "os"

// 从环境变量读取配置
func createOptionsFromEnv() *client.Options {
    options := client.NewOptions()
    
    if proxy := os.Getenv("HTTP_PROXY"); proxy != "" {
        options.WithProxy(proxy)
    }
    
    if userAgent := os.Getenv("PYPI_USER_AGENT"); userAgent != "" {
        options.WithUserAgent(userAgent)
    }
    
    return options
}
```

### 5. 配置验证

```go
func validateOptions(options *client.Options) error {
    if options.Timeout <= 0 {
        return fmt.Errorf("timeout must be positive")
    }
    
    if options.MaxRetries < 0 {
        return fmt.Errorf("max retries cannot be negative")
    }
    
    if options.BaseURL == "" {
        return fmt.Errorf("base URL cannot be empty")
    }
    
    return nil
}
```

---

**下一步**: 查看 [镜像源配置](./mirrors.md) 了解如何选择和配置不同的镜像源。
