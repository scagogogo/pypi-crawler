# API 参考

本文档详细介绍了 PyPI Crawler 提供的所有 API 接口。

## 📋 目录

- [PyPIClient 接口](#pypiclient-接口)
- [客户端创建](#客户端创建)
- [包信息 API](#包信息-api)
- [搜索 API](#搜索-api)
- [安全 API](#安全-api)
- [索引 API](#索引-api)

## PyPIClient 接口

所有客户端都实现了 `api.PyPIClient` 接口：

```go
type PyPIClient interface {
    GetPackageInfo(ctx context.Context, packageName string) (*models.Package, error)
    GetPackageVersion(ctx context.Context, packageName string, version string) (*models.Package, error)
    GetPackageReleases(ctx context.Context, packageName string) ([]string, error)
    CheckPackageVulnerabilities(ctx context.Context, packageName string, version string) ([]models.Vulnerability, error)
    GetAllPackages(ctx context.Context) ([]string, error)
    GetPackageList(ctx context.Context) (map[string]struct{}, error)
    SearchPackages(ctx context.Context, keyword string, limit int) ([]string, error)
}
```

## 客户端创建

### 使用镜像源工厂

```go
import "github.com/scagogogo/pypi-crawler/pkg/pypi/mirrors"

// 官方源
client := mirrors.NewOfficialClient()

// 国内镜像源
client := mirrors.NewTsinghuaClient()  // 清华大学
client := mirrors.NewAliyunClient()    // 阿里云
client := mirrors.NewDoubanClient()    // 豆瓣
client := mirrors.NewTencentClient()   // 腾讯云
client := mirrors.NewUstcClient()      // 中科大
client := mirrors.NewNeteaseClient()   // 网易
```

### 使用自定义配置

```go
import (
    "github.com/scagogogo/pypi-crawler/pkg/pypi/client"
    "github.com/scagogogo/pypi-crawler/pkg/pypi/mirrors"
)

options := client.NewOptions().
    WithTimeout(30 * time.Second).
    WithMaxRetries(3).
    WithProxy("http://proxy:8080")

client := mirrors.NewOfficialClient(options)
```

## 包信息 API

### GetPackageInfo

获取指定包的最新版本信息。

**函数签名:**
```go
GetPackageInfo(ctx context.Context, packageName string) (*models.Package, error)
```

**参数:**
- `ctx`: 上下文，用于控制请求生命周期
- `packageName`: 包名（不区分大小写）

**返回值:**
- `*models.Package`: 包信息结构体
- `error`: 错误信息

**示例:**
```go
pkg, err := client.GetPackageInfo(ctx, "requests")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("包名: %s\n", pkg.Info.Name)
fmt.Printf("版本: %s\n", pkg.Info.Version)
fmt.Printf("摘要: %s\n", pkg.Info.Summary)
```

**可能的错误:**
- 包不存在
- 网络连接失败
- API 响应格式错误

### GetPackageVersion

获取指定包的特定版本信息。

**函数签名:**
```go
GetPackageVersion(ctx context.Context, packageName string, version string) (*models.Package, error)
```

**参数:**
- `ctx`: 上下文
- `packageName`: 包名
- `version`: 版本号（如 "2.28.0"）

**返回值:**
- `*models.Package`: 特定版本的包信息
- `error`: 错误信息

**示例:**
```go
pkg, err := client.GetPackageVersion(ctx, "requests", "2.28.0")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("版本: %s\n", pkg.Info.Version)
fmt.Printf("发布文件数: %d\n", len(pkg.Urls))
```

### GetPackageReleases

获取指定包的所有发布版本列表。

**函数签名:**
```go
GetPackageReleases(ctx context.Context, packageName string) ([]string, error)
```

**参数:**
- `ctx`: 上下文
- `packageName`: 包名

**返回值:**
- `[]string`: 版本号列表
- `error`: 错误信息

**示例:**
```go
versions, err := client.GetPackageReleases(ctx, "requests")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("共有 %d 个版本\n", len(versions))
for i, version := range versions {
    if i < 5 { // 显示前5个版本
        fmt.Printf("  %s\n", version)
    }
}
```

## 搜索 API

### SearchPackages

根据关键词搜索包。

**函数签名:**
```go
SearchPackages(ctx context.Context, keyword string, limit int) ([]string, error)
```

**参数:**
- `ctx`: 上下文
- `keyword`: 搜索关键词
- `limit`: 最大返回结果数（0 表示使用默认值 100）

**返回值:**
- `[]string`: 匹配的包名列表
- `error`: 错误信息

**示例:**
```go
results, err := client.SearchPackages(ctx, "flask", 10)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("找到 %d 个包:\n", len(results))
for _, pkg := range results {
    fmt.Printf("  %s\n", pkg)
}
```

**注意:**
- 搜索是基于包名的简单字符串匹配
- 搜索不区分大小写
- 结果按包名字母顺序返回

## 安全 API

### CheckPackageVulnerabilities

检查指定包和版本是否存在已知安全漏洞。

**函数签名:**
```go
CheckPackageVulnerabilities(ctx context.Context, packageName string, version string) ([]models.Vulnerability, error)
```

**参数:**
- `ctx`: 上下文
- `packageName`: 包名
- `version`: 版本号

**返回值:**
- `[]models.Vulnerability`: 漏洞信息列表
- `error`: 错误信息

**示例:**
```go
vulns, err := client.CheckPackageVulnerabilities(ctx, "requests", "2.25.0")
if err != nil {
    log.Fatal(err)
}

if len(vulns) == 0 {
    fmt.Println("未发现已知漏洞")
} else {
    fmt.Printf("发现 %d 个漏洞:\n", len(vulns))
    for _, vuln := range vulns {
        fmt.Printf("  ID: %s\n", vuln.ID)
        fmt.Printf("  摘要: %s\n", vuln.Summary)
        if vuln.HasCVE() {
            fmt.Printf("  CVE: %v\n", vuln.GetCVEs())
        }
    }
}
```

## 索引 API

### GetAllPackages

获取 PyPI 仓库中所有包的列表。

**函数签名:**
```go
GetAllPackages(ctx context.Context) ([]string, error)
```

**参数:**
- `ctx`: 上下文

**返回值:**
- `[]string`: 所有包名的列表
- `error`: 错误信息

**示例:**
```go
packages, err := client.GetAllPackages(ctx)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("PyPI 中共有 %d 个包\n", len(packages))
```

**注意:**
- 此操作可能需要较长时间（几秒到几十秒）
- 返回的包名数量通常超过 40 万个
- 建议设置较长的超时时间

### GetPackageList

获取所有包的列表，以 map 形式返回，便于快速查找。

**函数签名:**
```go
GetPackageList(ctx context.Context) (map[string]struct{}, error)
```

**参数:**
- `ctx`: 上下文

**返回值:**
- `map[string]struct{}`: 包名集合
- `error`: 错误信息

**示例:**
```go
packageMap, err := client.GetPackageList(ctx)
if err != nil {
    log.Fatal(err)
}

// 检查包是否存在
if _, exists := packageMap["requests"]; exists {
    fmt.Println("requests 包存在")
}
```

## 错误处理

所有 API 方法都可能返回以下类型的错误：

1. **网络错误**: 连接失败、超时等
2. **HTTP 错误**: 404（包不存在）、500（服务器错误）等
3. **解析错误**: JSON 解析失败
4. **上下文错误**: 上下文取消或超时

**错误处理示例:**
```go
pkg, err := client.GetPackageInfo(ctx, "nonexistent")
if err != nil {
    if strings.Contains(err.Error(), "404") {
        fmt.Println("包不存在")
    } else {
        fmt.Printf("其他错误: %v\n", err)
    }
    return
}
```

## 性能建议

1. **复用客户端**: 创建一次客户端，多次使用
2. **设置合适的超时**: 根据网络环境调整超时时间
3. **使用上下文**: 利用上下文控制请求生命周期
4. **选择合适的镜像源**: 使用地理位置最近的镜像源
5. **批量操作**: 避免在循环中频繁调用 API

---

**下一步**: 查看 [数据模型](./data-models.md) 了解返回数据的详细结构。
