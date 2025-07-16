# 镜像源配置

本文档介绍 PyPI Crawler 支持的所有镜像源以及如何选择和配置最适合的镜像源。

## 📋 目录

- [支持的镜像源](#支持的镜像源)
- [镜像源选择指南](#镜像源选择指南)
- [使用方法](#使用方法)
- [性能对比](#性能对比)
- [故障转移](#故障转移)
- [自定义镜像源](#自定义镜像源)

## 支持的镜像源

PyPI Crawler 内置支持以下镜像源：

### 官方源

| 名称 | URL | 地区 | 工厂函数 |
|------|-----|------|----------|
| PyPI 官方 | `https://pypi.org` | 全球 | `mirrors.NewOfficialClient()` |

### 中国大陆镜像源

| 名称 | URL | 维护方 | 工厂函数 |
|------|-----|--------|----------|
| 清华大学 | `https://pypi.tuna.tsinghua.edu.cn` | 清华大学 TUNA | `mirrors.NewTsinghuaClient()` |
| 阿里云 | `https://mirrors.aliyun.com/pypi` | 阿里云 | `mirrors.NewAliyunClient()` |
| 豆瓣 | `https://pypi.doubanio.com` | 豆瓣 | `mirrors.NewDoubanClient()` |
| 腾讯云 | `https://mirrors.cloud.tencent.com/pypi` | 腾讯云 | `mirrors.NewTencentClient()` |
| 中科大 | `https://pypi.mirrors.ustc.edu.cn` | 中国科技大学 | `mirrors.NewUstcClient()` |
| 网易 | `https://mirrors.163.com/pypi` | 网易 | `mirrors.NewNeteaseClient()` |

## 镜像源选择指南

### 按地理位置选择

#### 中国大陆用户

**推荐顺序:**
1. **清华大学镜像** - 更新及时，稳定性好
2. **阿里云镜像** - 商业级稳定性
3. **中科大镜像** - 学术网络友好
4. **腾讯云镜像** - 企业级服务

```go
// 中国大陆推荐配置
client := mirrors.NewTsinghuaClient()
```

#### 海外用户

**推荐:**
- **PyPI 官方源** - 最新最全，权威可靠

```go
// 海外推荐配置
client := mirrors.NewOfficialClient()
```

### 按使用场景选择

#### 开发测试

```go
// 开发环境：使用快速镜像源
client := mirrors.NewTsinghuaClient()
```

#### 生产环境

```go
// 生产环境：使用稳定镜像源
options := client.NewOptions().
    WithTimeout(30 * time.Second).
    WithMaxRetries(3)

client := mirrors.NewAliyunClient(options)
```

#### CI/CD 环境

```go
// CI/CD：使用可靠镜像源
options := client.NewOptions().
    WithTimeout(60 * time.Second).
    WithMaxRetries(5)

client := mirrors.NewTsinghuaClient(options)
```

## 使用方法

### 基本使用

```go
import "github.com/scagogogo/pypi-crawler/pkg/pypi/mirrors"

// 使用官方源
client := mirrors.NewOfficialClient()

// 使用清华镜像
client := mirrors.NewTsinghuaClient()

// 使用阿里云镜像
client := mirrors.NewAliyunClient()
```

### 带配置使用

```go
import (
    "time"
    "github.com/scagogogo/pypi-crawler/pkg/pypi/client"
    "github.com/scagogogo/pypi-crawler/pkg/pypi/mirrors"
)

// 创建自定义配置
options := client.NewOptions().
    WithTimeout(15 * time.Second).
    WithMaxRetries(3).
    WithUserAgent("MyApp/1.0")

// 使用配置创建客户端
client := mirrors.NewTsinghuaClient(options)
```

### 动态选择镜像源

```go
func createClientByRegion(region string) api.PyPIClient {
    options := client.NewOptions().
        WithTimeout(30 * time.Second).
        WithMaxRetries(3)

    switch region {
    case "cn":
        return mirrors.NewTsinghuaClient(options)
    case "us":
        return mirrors.NewOfficialClient(options)
    case "asia":
        return mirrors.NewAliyunClient(options)
    default:
        return mirrors.NewOfficialClient(options)
    }
}
```

## 性能对比

### 延迟测试（毫秒）

| 镜像源 | 北京 | 上海 | 广州 | 海外 |
|--------|------|------|------|------|
| 官方源 | 200-300 | 180-250 | 220-280 | 50-100 |
| 清华镜像 | 10-30 | 20-40 | 30-50 | 150-200 |
| 阿里云镜像 | 15-35 | 10-25 | 20-40 | 100-150 |
| 豆瓣镜像 | 20-40 | 15-30 | 25-45 | 180-220 |

*注意：实际性能可能因网络环境而异*

### 同步频率

| 镜像源 | 同步频率 | 延迟 |
|--------|----------|------|
| 官方源 | 实时 | 0 |
| 清华镜像 | 5分钟 | < 5分钟 |
| 阿里云镜像 | 10分钟 | < 10分钟 |
| 中科大镜像 | 5分钟 | < 5分钟 |

## 故障转移

### 自动故障转移

```go
func createResilientClient() api.PyPIClient {
    // 尝试顺序：清华 -> 阿里云 -> 官方
    mirrors := []func() api.PyPIClient{
        func() api.PyPIClient { return mirrors.NewTsinghuaClient() },
        func() api.PyPIClient { return mirrors.NewAliyunClient() },
        func() api.PyPIClient { return mirrors.NewOfficialClient() },
    }

    for _, createClient := range mirrors {
        client := createClient()
        
        // 测试连接
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        _, err := client.GetPackageInfo(ctx, "requests")
        cancel()
        
        if err == nil {
            return client
        }
    }
    
    // 如果都失败，返回官方源
    return mirrors.NewOfficialClient()
}
```

### 手动故障转移

```go
type MultiMirrorClient struct {
    clients []api.PyPIClient
    current int
}

func NewMultiMirrorClient() *MultiMirrorClient {
    return &MultiMirrorClient{
        clients: []api.PyPIClient{
            mirrors.NewTsinghuaClient(),
            mirrors.NewAliyunClient(),
            mirrors.NewOfficialClient(),
        },
        current: 0,
    }
}

func (m *MultiMirrorClient) GetPackageInfo(ctx context.Context, packageName string) (*models.Package, error) {
    for i := 0; i < len(m.clients); i++ {
        client := m.clients[(m.current+i)%len(m.clients)]
        
        pkg, err := client.GetPackageInfo(ctx, packageName)
        if err == nil {
            m.current = (m.current + i) % len(m.clients)
            return pkg, nil
        }
    }
    
    return nil, fmt.Errorf("all mirrors failed")
}
```

## 自定义镜像源

### 使用自定义 URL

```go
import "github.com/scagogogo/pypi-crawler/pkg/pypi/client"

// 创建自定义镜像源客户端
options := client.NewOptions().
    WithBaseURL("https://my-pypi-mirror.com")

client := client.NewClient(options)
```

### 企业内部镜像

```go
func createEnterpriseClient() api.PyPIClient {
    options := client.NewOptions().
        WithBaseURL("https://pypi.internal.company.com").
        WithProxy("http://proxy.company.com:8080").
        WithUserAgent("CompanyApp/1.0 (admin@company.com)").
        WithTimeout(60 * time.Second)

    return client.NewClient(options)
}
```

### 镜像源健康检查

```go
func checkMirrorHealth(baseURL string) bool {
    options := client.NewOptions().
        WithBaseURL(baseURL).
        WithTimeout(5 * time.Second).
        WithMaxRetries(1)

    client := client.NewClient(options)
    
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    _, err := client.GetPackageInfo(ctx, "requests")
    return err == nil
}

// 使用示例
func selectBestMirror() api.PyPIClient {
    mirrors := map[string]string{
        "tsinghua": "https://pypi.tuna.tsinghua.edu.cn",
        "aliyun":   "https://mirrors.aliyun.com/pypi",
        "official": "https://pypi.org",
    }

    for name, url := range mirrors {
        if checkMirrorHealth(url) {
            fmt.Printf("使用镜像源: %s\n", name)
            options := client.NewOptions().WithBaseURL(url)
            return client.NewClient(options)
        }
    }

    // 默认使用官方源
    return mirrors.NewOfficialClient()
}
```

## 镜像源配置最佳实践

### 1. 环境变量配置

```go
import "os"

func createClientFromEnv() api.PyPIClient {
    mirrorType := os.Getenv("PYPI_MIRROR")
    
    switch mirrorType {
    case "tsinghua":
        return mirrors.NewTsinghuaClient()
    case "aliyun":
        return mirrors.NewAliyunClient()
    case "official":
        return mirrors.NewOfficialClient()
    default:
        // 根据地区自动选择
        return autoSelectMirror()
    }
}

func autoSelectMirror() api.PyPIClient {
    // 简单的地区检测逻辑
    timezone := os.Getenv("TZ")
    if strings.Contains(timezone, "Asia/Shanghai") {
        return mirrors.NewTsinghuaClient()
    }
    return mirrors.NewOfficialClient()
}
```

### 2. 配置文件支持

```go
type MirrorConfig struct {
    Type    string `yaml:"type"`
    URL     string `yaml:"url,omitempty"`
    Timeout int    `yaml:"timeout"`
    Retries int    `yaml:"retries"`
}

func createClientFromConfig(config MirrorConfig) api.PyPIClient {
    options := client.NewOptions().
        WithTimeout(time.Duration(config.Timeout) * time.Second).
        WithMaxRetries(config.Retries)

    switch config.Type {
    case "tsinghua":
        return mirrors.NewTsinghuaClient(options)
    case "aliyun":
        return mirrors.NewAliyunClient(options)
    case "custom":
        options.WithBaseURL(config.URL)
        return client.NewClient(options)
    default:
        return mirrors.NewOfficialClient(options)
    }
}
```

### 3. 性能监控

```go
type MirrorStats struct {
    URL           string
    RequestCount  int64
    ErrorCount    int64
    AvgLatency    time.Duration
}

func (s *MirrorStats) RecordRequest(latency time.Duration, err error) {
    s.RequestCount++
    if err != nil {
        s.ErrorCount++
    }
    
    // 简单的移动平均
    s.AvgLatency = (s.AvgLatency + latency) / 2
}

func (s *MirrorStats) SuccessRate() float64 {
    if s.RequestCount == 0 {
        return 0
    }
    return float64(s.RequestCount-s.ErrorCount) / float64(s.RequestCount)
}
```

---

**下一步**: 查看 [错误处理](./error-handling.md) 了解如何处理各种错误情况。
