# 快速开始

本指南将帮助您在 5 分钟内开始使用 PyPI Crawler。

## 📦 安装

```bash
go get -u github.com/scagogogo/pypi-crawler
```

## 🚀 第一个程序

创建一个简单的程序来获取包信息：

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/scagogogo/pypi-crawler/pkg/pypi/mirrors"
)

func main() {
    // 1. 创建客户端
    client := mirrors.NewOfficialClient()

    // 2. 创建上下文
    ctx := context.Background()

    // 3. 获取包信息
    pkg, err := client.GetPackageInfo(ctx, "requests")
    if err != nil {
        log.Fatalf("获取包信息失败: %v", err)
    }

    // 4. 显示结果
    fmt.Printf("包名: %s\n", pkg.Info.Name)
    fmt.Printf("版本: %s\n", pkg.Info.Version)
    fmt.Printf("摘要: %s\n", pkg.Info.Summary)
    fmt.Printf("作者: %s\n", pkg.Info.Author)
}
```

运行程序：

```bash
go run main.go
```

输出示例：
```
包名: requests
版本: 2.31.0
摘要: Python HTTP for Humans.
作者: Kenneth Reitz
```

## 🌐 使用国内镜像源

如果您在中国，建议使用国内镜像源以提高访问速度：

```go
// 使用清华大学镜像源
client := mirrors.NewTsinghuaClient()

// 使用阿里云镜像源
client := mirrors.NewAliyunClient()

// 使用豆瓣镜像源
client := mirrors.NewDoubanClient()
```

## 🔧 自定义配置

您可以自定义客户端配置：

```go
package main

import (
    "context"
    "time"

    "github.com/scagogogo/pypi-crawler/pkg/pypi/client"
    "github.com/scagogogo/pypi-crawler/pkg/pypi/mirrors"
)

func main() {
    // 创建自定义配置
    options := client.NewOptions().
        WithTimeout(15 * time.Second).        // 设置15秒超时
        WithMaxRetries(5).                    // 最多重试5次
        WithUserAgent("MyApp/1.0")            // 自定义User-Agent

    // 使用自定义配置创建客户端
    pypiClient := mirrors.NewOfficialClient(options)

    // 使用客户端...
    ctx := context.Background()
    pkg, err := pypiClient.GetPackageInfo(ctx, "flask")
    // ...
}
```

## 📋 常用操作

### 获取特定版本信息

```go
// 获取 requests 包的 2.28.0 版本信息
pkg, err := client.GetPackageVersion(ctx, "requests", "2.28.0")
```

### 获取所有版本列表

```go
// 获取 requests 包的所有版本
versions, err := client.GetPackageReleases(ctx, "requests")
fmt.Printf("共有 %d 个版本\n", len(versions))
```

### 搜索包

```go
// 搜索包含 "flask" 关键词的包，最多返回 10 个结果
results, err := client.SearchPackages(ctx, "flask", 10)
```

### 检查安全漏洞

```go
// 检查 requests 2.25.0 版本的安全漏洞
vulns, err := client.CheckPackageVulnerabilities(ctx, "requests", "2.25.0")
if len(vulns) > 0 {
    fmt.Printf("发现 %d 个漏洞\n", len(vulns))
}
```

## 🛠️ 错误处理

始终检查错误并适当处理：

```go
pkg, err := client.GetPackageInfo(ctx, "nonexistent-package")
if err != nil {
    // 处理错误
    fmt.Printf("错误: %v\n", err)
    return
}
```

## 📚 下一步

- 查看 [API 参考](./api-reference.md) 了解所有可用的方法
- 查看 [数据模型](./data-models.md) 了解返回数据的结构
- 查看 [示例代码](./examples.md) 获取更多使用示例
- 查看 [最佳实践](./best-practices.md) 了解性能优化建议

## 💡 提示

1. **使用上下文**: 始终传递 `context.Context` 以便控制请求生命周期
2. **选择合适的镜像源**: 根据您的地理位置选择最快的镜像源
3. **处理错误**: 网络请求可能失败，请妥善处理错误
4. **设置超时**: 为长时间运行的操作设置合适的超时时间

---

**恭喜！** 您已经学会了 PyPI Crawler 的基本使用方法。现在可以开始构建您的应用程序了！
