# 常见问题

本文档收集了使用 PyPI Crawler 时的常见问题和解决方案。

## 📋 目录

- [安装和配置](#安装和配置)
- [网络和连接](#网络和连接)
- [性能问题](#性能问题)
- [错误处理](#错误处理)
- [功能使用](#功能使用)
- [故障排除](#故障排除)

## 安装和配置

### Q: 如何安装 PyPI Crawler？

**A:** 使用 Go 模块安装：

```bash
go get -u github.com/scagogogo/pypi-crawler
```

确保您的 Go 版本为 1.19 或更高。

### Q: 支持哪些 Go 版本？

**A:** PyPI Crawler 支持 Go 1.19 及以上版本。推荐使用最新的稳定版本。

### Q: 如何在项目中导入？

**A:** 在您的 Go 代码中导入：

```go
import (
    "github.com/scagogogo/pypi-crawler/pkg/pypi/mirrors"
    "github.com/scagogogo/pypi-crawler/pkg/pypi/client"
    "github.com/scagogogo/pypi-crawler/pkg/pypi/models"
)
```

## 网络和连接

### Q: 在中国大陆访问速度很慢怎么办？

**A:** 建议使用国内镜像源：

```go
// 推荐使用清华大学镜像
client := mirrors.NewTsinghuaClient()

// 或者阿里云镜像
client := mirrors.NewAliyunClient()

// 或者豆瓣镜像
client := mirrors.NewDoubanClient()
```

### Q: 如何配置代理？

**A:** 通过客户端选项配置代理：

```go
options := client.NewOptions().
    WithProxy("http://proxy.company.com:8080")

client := mirrors.NewOfficialClient(options)
```

支持的代理类型：
- HTTP 代理：`http://proxy:8080`
- HTTPS 代理：`https://proxy:8080`
- SOCKS5 代理：`socks5://proxy:1080`
- 带认证的代理：`http://user:pass@proxy:8080`

### Q: 连接超时怎么办？

**A:** 增加超时时间：

```go
options := client.NewOptions().
    WithTimeout(60 * time.Second) // 设置60秒超时

client := mirrors.NewTsinghuaClient(options)
```

对于不同操作建议的超时时间：
- 单个包查询：10-30 秒
- 获取包列表：60-120 秒
- 批量操作：5-10 分钟

### Q: 如何处理网络不稳定的情况？

**A:** 配置重试机制：

```go
options := client.NewOptions().
    WithMaxRetries(5).                    // 最多重试5次
    WithRetryDelay(2 * time.Second).      // 重试间隔2秒
    WithTimeout(30 * time.Second)         // 30秒超时

client := mirrors.NewTsinghuaClient(options)
```

## 性能问题

### Q: 批量获取包信息很慢怎么办？

**A:** 使用并发处理：

```go
func batchGetPackages(client api.PyPIClient, packages []string) {
    const maxConcurrency = 5
    semaphore := make(chan struct{}, maxConcurrency)
    var wg sync.WaitGroup

    for _, pkg := range packages {
        wg.Add(1)
        go func(packageName string) {
            defer wg.Done()
            
            semaphore <- struct{}{} // 获取信号量
            defer func() { <-semaphore }() // 释放信号量
            
            // 处理包
            processPackage(client, packageName)
        }(pkg)
    }
    
    wg.Wait()
}
```

### Q: 如何减少重复请求？

**A:** 实现缓存机制：

```go
type CachedClient struct {
    client api.PyPIClient
    cache  map[string]*models.Package
    mutex  sync.RWMutex
}

func (c *CachedClient) GetPackageInfo(ctx context.Context, packageName string) (*models.Package, error) {
    // 检查缓存
    c.mutex.RLock()
    if cached, exists := c.cache[packageName]; exists {
        c.mutex.RUnlock()
        return cached, nil
    }
    c.mutex.RUnlock()
    
    // 从API获取
    pkg, err := c.client.GetPackageInfo(ctx, packageName)
    if err != nil {
        return nil, err
    }
    
    // 存入缓存
    c.mutex.Lock()
    c.cache[packageName] = pkg
    c.mutex.Unlock()
    
    return pkg, nil
}
```

### Q: 内存使用过多怎么办？

**A:** 
1. 避免一次性加载所有包信息
2. 使用流式处理
3. 定期清理缓存
4. 限制并发数量

```go
// 分批处理
func processBatches(packages []string, batchSize int) {
    for i := 0; i < len(packages); i += batchSize {
        end := i + batchSize
        if end > len(packages) {
            end = len(packages)
        }
        
        batch := packages[i:end]
        processBatch(batch)
        
        // 强制垃圾回收
        runtime.GC()
    }
}
```

## 错误处理

### Q: 如何判断包是否存在？

**A:** 检查 404 错误：

```go
pkg, err := client.GetPackageInfo(ctx, "nonexistent-package")
if err != nil {
    if strings.Contains(err.Error(), "404") {
        fmt.Println("包不存在")
        return
    }
    // 其他错误
    return
}
```

### Q: 如何处理服务器错误？

**A:** 实现重试逻辑：

```go
func getPackageWithRetry(client api.PyPIClient, packageName string) (*models.Package, error) {
    maxRetries := 3
    
    for attempt := 0; attempt < maxRetries; attempt++ {
        pkg, err := client.GetPackageInfo(ctx, packageName)
        if err == nil {
            return pkg, nil
        }
        
        // 检查是否为服务器错误（5xx）
        if strings.Contains(err.Error(), "5") {
            time.Sleep(time.Duration(attempt+1) * time.Second)
            continue
        }
        
        // 非服务器错误，不重试
        return nil, err
    }
    
    return nil, fmt.Errorf("重试 %d 次后仍然失败", maxRetries)
}
```

### Q: 如何处理 JSON 解析错误？

**A:** 添加数据验证：

```go
pkg, err := client.GetPackageInfo(ctx, packageName)
if err != nil {
    return nil, err
}

// 验证关键字段
if pkg.Info == nil {
    return nil, fmt.Errorf("包信息为空")
}

if pkg.Info.Name == "" {
    return nil, fmt.Errorf("包名为空")
}

if pkg.Info.Version == "" {
    return nil, fmt.Errorf("版本号为空")
}
```

## 功能使用

### Q: 如何搜索包？

**A:** 使用 SearchPackages 方法：

```go
// 搜索包含 "web" 关键词的包，最多返回 10 个结果
results, err := client.SearchPackages(ctx, "web", 10)
if err != nil {
    log.Fatal(err)
}

for _, pkg := range results {
    fmt.Println(pkg)
}
```

### Q: 如何获取包的所有版本？

**A:** 使用 GetPackageReleases 方法：

```go
versions, err := client.GetPackageReleases(ctx, "requests")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("共有 %d 个版本\n", len(versions))
for _, version := range versions {
    fmt.Println(version)
}
```

### Q: 如何检查包的安全漏洞？

**A:** 使用 CheckPackageVulnerabilities 方法：

```go
vulns, err := client.CheckPackageVulnerabilities(ctx, "requests", "2.25.0")
if err != nil {
    log.Fatal(err)
}

if len(vulns) == 0 {
    fmt.Println("未发现漏洞")
} else {
    fmt.Printf("发现 %d 个漏洞\n", len(vulns))
    for _, vuln := range vulns {
        fmt.Printf("- %s: %s\n", vuln.ID, vuln.Summary)
    }
}
```

### Q: 如何获取包的依赖信息？

**A:** 通过包信息获取依赖：

```go
pkg, err := client.GetPackageInfo(ctx, "flask")
if err != nil {
    log.Fatal(err)
}

dependencies := pkg.Info.GetAllDependencies()
fmt.Printf("依赖项 (%d):\n", len(dependencies))
for _, dep := range dependencies {
    fmt.Printf("- %s\n", dep)
}
```

### Q: 如何获取包的下载文件信息？

**A:** 通过 Urls 字段获取：

```go
pkg, err := client.GetPackageInfo(ctx, "numpy")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("发布文件 (%d):\n", len(pkg.Urls))
for _, file := range pkg.Urls {
    fmt.Printf("- %s (%s, %.2f MB)\n", 
        file.Filename, 
        file.PackageType, 
        float64(file.Size)/(1024*1024))
}
```

## 故障排除

### Q: 程序运行时出现 panic 怎么办？

**A:** 
1. 检查是否正确处理了错误
2. 确保不访问 nil 指针
3. 添加适当的错误检查

```go
// 安全的访问方式
pkg, err := client.GetPackageInfo(ctx, packageName)
if err != nil {
    return err
}

if pkg == nil || pkg.Info == nil {
    return fmt.Errorf("包信息为空")
}

// 现在可以安全访问 pkg.Info
fmt.Println(pkg.Info.Name)
```

### Q: 如何启用调试日志？

**A:** 使用日志包装器：

```go
import "log/slog"

type LoggedClient struct {
    client api.PyPIClient
    logger *slog.Logger
}

func (lc *LoggedClient) GetPackageInfo(ctx context.Context, packageName string) (*models.Package, error) {
    lc.logger.Info("获取包信息", slog.String("package", packageName))
    
    pkg, err := lc.client.GetPackageInfo(ctx, packageName)
    if err != nil {
        lc.logger.Error("获取失败", slog.String("error", err.Error()))
        return nil, err
    }
    
    lc.logger.Info("获取成功", slog.String("version", pkg.Info.Version))
    return pkg, nil
}
```

### Q: 如何测试网络连接？

**A:** 实现健康检查：

```go
func testConnection(client api.PyPIClient) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    // 尝试获取一个知名包
    _, err := client.GetPackageInfo(ctx, "requests")
    return err
}

// 使用示例
if err := testConnection(client); err != nil {
    fmt.Printf("连接测试失败: %v\n", err)
    fmt.Println("请检查网络连接或更换镜像源")
} else {
    fmt.Println("连接正常")
}
```

### Q: 如何报告 Bug？

**A:** 请在 GitHub 仓库提交 Issue，包含以下信息：

1. Go 版本
2. PyPI Crawler 版本
3. 操作系统和版本
4. 完整的错误信息
5. 重现步骤
6. 最小化的示例代码

```bash
# 获取版本信息
go version
go list -m github.com/scagogogo/pypi-crawler
```

### Q: 如何获取帮助？

**A:** 
1. 查看项目文档：[GitHub 仓库](https://github.com/scagogogo/pypi-crawler)
2. 查看示例代码：`examples/` 目录
3. 提交 Issue：描述您的问题
4. 参与讨论：GitHub Discussions

---

**还有其他问题？** 请在 [GitHub Issues](https://github.com/scagogogo/pypi-crawler/issues) 中提出，我们会及时回复。
