# 示例代码

本文档提供了丰富的 PyPI Crawler 使用示例，涵盖各种常见场景。

## 📋 目录

- [基础示例](#基础示例)
- [高级查询](#高级查询)
- [批量操作](#批量操作)
- [安全检查](#安全检查)
- [数据分析](#数据分析)
- [实用工具](#实用工具)
- [性能优化](#性能优化)

## 基础示例

### 获取包信息

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/scagogogo/pypi-crawler/pkg/pypi/mirrors"
)

func main() {
    // 创建客户端
    client := mirrors.NewTsinghuaClient()
    ctx := context.Background()

    // 获取包信息
    pkg, err := client.GetPackageInfo(ctx, "requests")
    if err != nil {
        log.Fatalf("获取包信息失败: %v", err)
    }

    // 显示基本信息
    fmt.Printf("包名: %s\n", pkg.Info.Name)
    fmt.Printf("版本: %s\n", pkg.Info.Version)
    fmt.Printf("摘要: %s\n", pkg.Info.Summary)
    fmt.Printf("作者: %s <%s>\n", pkg.Info.Author, pkg.Info.AuthorEmail)
    fmt.Printf("许可证: %s\n", pkg.Info.License)
    fmt.Printf("主页: %s\n", pkg.Info.HomePage)

    // 显示依赖信息
    deps := pkg.Info.GetAllDependencies()
    if len(deps) > 0 {
        fmt.Printf("\n依赖项 (%d):\n", len(deps))
        for i, dep := range deps {
            if i < 5 {
                fmt.Printf("  %d. %s\n", i+1, dep)
            } else {
                fmt.Printf("  ... 以及其他 %d 个依赖\n", len(deps)-5)
                break
            }
        }
    }

    // 显示项目链接
    urls := pkg.Info.GetProjectURLs()
    if len(urls) > 0 {
        fmt.Println("\n项目链接:")
        for name, url := range urls {
            fmt.Printf("  %s: %s\n", name, url)
        }
    }
}
```

### 搜索包

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/scagogogo/pypi-crawler/pkg/pypi/mirrors"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("用法: go run search.go <关键词>")
        return
    }

    keyword := os.Args[1]
    client := mirrors.NewTsinghuaClient()
    ctx := context.Background()

    // 搜索包
    results, err := client.SearchPackages(ctx, keyword, 20)
    if err != nil {
        log.Fatalf("搜索失败: %v", err)
    }

    fmt.Printf("搜索关键词 '%s' 找到 %d 个结果:\n\n", keyword, len(results))

    // 显示搜索结果
    for i, pkgName := range results {
        fmt.Printf("%d. %s\n", i+1, pkgName)
        
        // 获取包的简要信息
        if pkg, err := client.GetPackageInfo(ctx, pkgName); err == nil {
            fmt.Printf("   摘要: %s\n", pkg.Info.Summary)
            fmt.Printf("   版本: %s\n", pkg.Info.Version)
        }
        fmt.Println()
    }
}
```

## 高级查询

### 比较包版本

```go
package main

import (
    "context"
    "fmt"
    "log"
    "sort"
    "strings"

    "github.com/scagogogo/pypi-crawler/pkg/pypi/mirrors"
)

func main() {
    client := mirrors.NewTsinghuaClient()
    ctx := context.Background()

    packageName := "django"

    // 获取所有版本
    versions, err := client.GetPackageReleases(ctx, packageName)
    if err != nil {
        log.Fatalf("获取版本列表失败: %v", err)
    }

    fmt.Printf("包 %s 共有 %d 个版本\n\n", packageName, len(versions))

    // 过滤稳定版本（不包含 alpha, beta, rc）
    stableVersions := filterStableVersions(versions)
    fmt.Printf("稳定版本 (%d):\n", len(stableVersions))
    
    // 显示最近的10个稳定版本
    sort.Strings(stableVersions)
    start := len(stableVersions) - 10
    if start < 0 {
        start = 0
    }

    for i := len(stableVersions) - 1; i >= start; i-- {
        version := stableVersions[i]
        fmt.Printf("  %s", version)
        
        // 获取版本详细信息
        if pkg, err := client.GetPackageVersion(ctx, packageName, version); err == nil {
            if len(pkg.Urls) > 0 {
                uploadTime, _ := pkg.Urls[0].GetUploadTimeISO()
                fmt.Printf(" (发布于 %s)", uploadTime.Format("2006-01-02"))
            }
        }
        fmt.Println()
    }
}

func filterStableVersions(versions []string) []string {
    var stable []string
    for _, version := range versions {
        v := strings.ToLower(version)
        if !strings.Contains(v, "alpha") && 
           !strings.Contains(v, "beta") && 
           !strings.Contains(v, "rc") && 
           !strings.Contains(v, "dev") {
            stable = append(stable, version)
        }
    }
    return stable
}
```

### 分析包的发布文件

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/scagogogo/pypi-crawler/pkg/pypi/mirrors"
)

func main() {
    client := mirrors.NewTsinghuaClient()
    ctx := context.Background()

    packageName := "numpy"

    // 获取包信息
    pkg, err := client.GetPackageInfo(ctx, packageName)
    if err != nil {
        log.Fatalf("获取包信息失败: %v", err)
    }

    fmt.Printf("包 %s 版本 %s 的发布文件分析:\n\n", pkg.Info.Name, pkg.Info.Version)

    // 分析发布文件
    var totalSize int64
    wheelCount := 0
    sdistCount := 0
    
    fmt.Println("发布文件列表:")
    for i, file := range pkg.Urls {
        fmt.Printf("%d. %s\n", i+1, file.Filename)
        fmt.Printf("   类型: %s\n", file.PackageType)
        fmt.Printf("   大小: %.2f MB\n", float64(file.Size)/(1024*1024))
        fmt.Printf("   Python版本: %s\n", file.PythonVersion)
        
        if file.RequiresPython != "" {
            fmt.Printf("   需要Python: %s\n", file.RequiresPython)
        }
        
        // 显示哈希值
        fmt.Printf("   SHA256: %s\n", file.Digests.SHA256)
        
        if file.IsYanked() {
            fmt.Printf("   ⚠️  已撤回: %s\n", file.YankedReason)
        }
        
        fmt.Println()
        
        totalSize += file.Size
        if file.IsWheel() {
            wheelCount++
        } else if file.IsSourceDist() {
            sdistCount++
        }
    }

    // 统计信息
    fmt.Printf("统计信息:\n")
    fmt.Printf("  总文件数: %d\n", len(pkg.Urls))
    fmt.Printf("  Wheel文件: %d\n", wheelCount)
    fmt.Printf("  源码包: %d\n", sdistCount)
    fmt.Printf("  总大小: %.2f MB\n", float64(totalSize)/(1024*1024))
}
```

## 批量操作

### 批量获取包信息

```go
package main

import (
    "context"
    "fmt"
    "log"
    "sync"
    "time"

    "github.com/scagogogo/pypi-crawler/pkg/pypi/mirrors"
    "github.com/scagogogo/pypi-crawler/pkg/pypi/models"
)

func main() {
    packages := []string{
        "requests", "flask", "django", "numpy", "pandas",
        "matplotlib", "scipy", "scikit-learn", "tensorflow", "pytorch",
    }

    client := mirrors.NewTsinghuaClient()
    
    // 并发获取包信息
    results := batchGetPackageInfo(client, packages, 3) // 最多3个并发

    // 显示结果
    fmt.Printf("成功获取 %d 个包的信息:\n\n", len(results))
    
    for _, result := range results {
        if result.Error != nil {
            fmt.Printf("❌ %s: %v\n", result.PackageName, result.Error)
        } else {
            pkg := result.Package
            fmt.Printf("✅ %s (%s)\n", pkg.Info.Name, pkg.Info.Version)
            fmt.Printf("   摘要: %s\n", pkg.Info.Summary)
            fmt.Printf("   作者: %s\n", pkg.Info.Author)
            fmt.Println()
        }
    }
}

type PackageResult struct {
    PackageName string
    Package     *models.Package
    Error       error
}

func batchGetPackageInfo(client api.PyPIClient, packages []string, concurrency int) []PackageResult {
    results := make([]PackageResult, len(packages))
    
    // 创建工作池
    jobs := make(chan int, len(packages))
    var wg sync.WaitGroup
    
    // 启动工作协程
    for i := 0; i < concurrency; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            
            for index := range jobs {
                packageName := packages[index]
                
                ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
                pkg, err := client.GetPackageInfo(ctx, packageName)
                cancel()
                
                results[index] = PackageResult{
                    PackageName: packageName,
                    Package:     pkg,
                    Error:       err,
                }
                
                // 避免请求过于频繁
                time.Sleep(100 * time.Millisecond)
            }
        }()
    }
    
    // 发送任务
    for i := range packages {
        jobs <- i
    }
    close(jobs)
    
    // 等待完成
    wg.Wait()
    
    return results
}
```

### 导出包信息到CSV

```go
package main

import (
    "context"
    "encoding/csv"
    "fmt"
    "log"
    "os"
    "strconv"

    "github.com/scagogogo/pypi-crawler/pkg/pypi/mirrors"
)

func main() {
    packages := []string{"requests", "flask", "django", "fastapi", "tornado"}
    
    client := mirrors.NewTsinghuaClient()
    ctx := context.Background()

    // 创建CSV文件
    file, err := os.Create("packages.csv")
    if err != nil {
        log.Fatalf("创建文件失败: %v", err)
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    // 写入CSV头部
    headers := []string{
        "包名", "版本", "摘要", "作者", "许可证", 
        "主页", "依赖数量", "文件数量", "总大小(MB)",
    }
    writer.Write(headers)

    // 获取并写入包信息
    for _, packageName := range packages {
        fmt.Printf("正在处理 %s...\n", packageName)
        
        pkg, err := client.GetPackageInfo(ctx, packageName)
        if err != nil {
            fmt.Printf("获取 %s 失败: %v\n", packageName, err)
            continue
        }

        // 计算总大小
        var totalSize int64
        for _, file := range pkg.Urls {
            totalSize += file.Size
        }

        // 准备CSV行数据
        row := []string{
            pkg.Info.Name,
            pkg.Info.Version,
            pkg.Info.Summary,
            pkg.Info.Author,
            pkg.Info.License,
            pkg.Info.HomePage,
            strconv.Itoa(len(pkg.Info.GetAllDependencies())),
            strconv.Itoa(len(pkg.Urls)),
            fmt.Sprintf("%.2f", float64(totalSize)/(1024*1024)),
        }

        writer.Write(row)
    }

    fmt.Println("导出完成: packages.csv")
}
```

## 安全检查

### 检查包的安全漏洞

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/scagogogo/pypi-crawler/pkg/pypi/mirrors"
)

func main() {
    client := mirrors.NewTsinghuaClient()
    ctx := context.Background()

    // 要检查的包和版本
    packagesToCheck := map[string][]string{
        "requests": {"2.25.0", "2.28.0", "2.31.0"},
        "django":   {"3.0.0", "3.2.0", "4.0.0"},
        "flask":    {"1.0.0", "2.0.0", "2.3.0"},
    }

    fmt.Println("安全漏洞检查报告")
    fmt.Println("================\n")

    totalVulns := 0

    for packageName, versions := range packagesToCheck {
        fmt.Printf("📦 检查包: %s\n", packageName)
        fmt.Println(strings.Repeat("-", 40))

        for _, version := range versions {
            fmt.Printf("🔍 版本 %s: ", version)

            vulns, err := client.CheckPackageVulnerabilities(ctx, packageName, version)
            if err != nil {
                fmt.Printf("检查失败 - %v\n", err)
                continue
            }

            if len(vulns) == 0 {
                fmt.Println("✅ 未发现漏洞")
            } else {
                fmt.Printf("⚠️  发现 %d 个漏洞\n", len(vulns))
                totalVulns += len(vulns)

                for i, vuln := range vulns {
                    fmt.Printf("    %d. %s\n", i+1, vuln.ID)
                    fmt.Printf("       摘要: %s\n", vuln.Summary)
                    
                    if vuln.HasCVE() {
                        fmt.Printf("       CVE: %v\n", vuln.GetCVEs())
                    }
                    
                    if len(vuln.FixedIn) > 0 {
                        fmt.Printf("       已修复版本: %v\n", vuln.FixedIn)
                    }
                    
                    if vuln.Link != "" {
                        fmt.Printf("       详情: %s\n", vuln.Link)
                    }
                    fmt.Println()
                }
            }
        }
        fmt.Println()
    }

    fmt.Printf("总计发现 %d 个安全漏洞\n", totalVulns)
    
    if totalVulns > 0 {
        fmt.Println("\n建议:")
        fmt.Println("1. 升级到已修复的版本")
        fmt.Println("2. 查看漏洞详情链接了解更多信息")
        fmt.Println("3. 评估漏洞对您应用的影响")
    }
}
```

## 数据分析

### 分析包的流行度趋势

```go
package main

import (
    "context"
    "fmt"
    "log"
    "sort"
    "strings"
    "time"

    "github.com/scagogogo/pypi-crawler/pkg/pypi/mirrors"
)

func main() {
    client := mirrors.NewTsinghuaClient()
    ctx := context.Background()

    packageName := "flask"

    // 获取所有版本
    versions, err := client.GetPackageReleases(ctx, packageName)
    if err != nil {
        log.Fatalf("获取版本失败: %v", err)
    }

    fmt.Printf("包 %s 的发布历史分析\n", packageName)
    fmt.Println(strings.Repeat("=", 50))

    // 分析版本发布模式
    versionInfo := analyzeVersionHistory(client, ctx, packageName, versions)

    // 按年份统计
    yearStats := make(map[int]int)
    for _, info := range versionInfo {
        if !info.ReleaseTime.IsZero() {
            year := info.ReleaseTime.Year()
            yearStats[year]++
        }
    }

    fmt.Println("\n📊 按年份发布统计:")
    var years []int
    for year := range yearStats {
        years = append(years, year)
    }
    sort.Ints(years)

    for _, year := range years {
        count := yearStats[year]
        fmt.Printf("  %d: %d 个版本\n", year, count)
    }

    // 最近的版本
    fmt.Println("\n🕒 最近的10个版本:")
    sort.Slice(versionInfo, func(i, j int) bool {
        return versionInfo[i].ReleaseTime.After(versionInfo[j].ReleaseTime)
    })

    for i, info := range versionInfo {
        if i >= 10 {
            break
        }
        if !info.ReleaseTime.IsZero() {
            fmt.Printf("  %s - %s\n", info.Version, info.ReleaseTime.Format("2006-01-02"))
        }
    }
}

type VersionInfo struct {
    Version     string
    ReleaseTime time.Time
    FileCount   int
}

func analyzeVersionHistory(client api.PyPIClient, ctx context.Context, packageName string, versions []string) []VersionInfo {
    var versionInfo []VersionInfo

    for _, version := range versions {
        pkg, err := client.GetPackageVersion(ctx, packageName, version)
        if err != nil {
            continue
        }

        info := VersionInfo{
            Version:   version,
            FileCount: len(pkg.Urls),
        }

        // 获取发布时间
        if len(pkg.Urls) > 0 {
            if releaseTime, err := pkg.Urls[0].GetUploadTimeISO(); err == nil {
                info.ReleaseTime = releaseTime
            }
        }

        versionInfo = append(versionInfo, info)

        // 避免请求过于频繁
        time.Sleep(100 * time.Millisecond)
    }

    return versionInfo
}
```

## 实用工具

### 包依赖分析器

```go
package main

import (
    "context"
    "fmt"
    "log"
    "strings"

    "github.com/scagogogo/pypi-crawler/pkg/pypi/mirrors"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("用法: go run deps.go <包名>")
        return
    }

    packageName := os.Args[1]
    client := mirrors.NewTsinghuaClient()
    ctx := context.Background()

    // 分析依赖
    deps, err := analyzeDependencies(client, ctx, packageName, 0, make(map[string]bool))
    if err != nil {
        log.Fatalf("分析依赖失败: %v", err)
    }

    fmt.Printf("包 %s 的依赖分析:\n", packageName)
    fmt.Println(strings.Repeat("=", 50))
    
    printDependencyTree(deps, 0)
}

type Dependency struct {
    Name         string
    Version      string
    Dependencies []*Dependency
}

func analyzeDependencies(client api.PyPIClient, ctx context.Context, packageName string, depth int, visited map[string]bool) (*Dependency, error) {
    // 防止循环依赖
    if visited[packageName] || depth > 3 {
        return &Dependency{Name: packageName, Version: "..."}, nil
    }
    
    visited[packageName] = true
    defer func() { visited[packageName] = false }()

    // 获取包信息
    pkg, err := client.GetPackageInfo(ctx, packageName)
    if err != nil {
        return nil, err
    }

    dep := &Dependency{
        Name:    pkg.Info.Name,
        Version: pkg.Info.Version,
    }

    // 解析依赖
    for _, depStr := range pkg.Info.GetAllDependencies() {
        depName := parseDependencyName(depStr)
        if depName != "" && !visited[depName] {
            if subDep, err := analyzeDependencies(client, ctx, depName, depth+1, visited); err == nil {
                dep.Dependencies = append(dep.Dependencies, subDep)
            }
        }
    }

    return dep, nil
}

func parseDependencyName(depStr string) string {
    // 简单解析依赖名称（去除版本约束）
    parts := strings.FieldsFunc(depStr, func(r rune) bool {
        return r == '>' || r == '<' || r == '=' || r == '!' || r == '~' || r == ' '
    })
    
    if len(parts) > 0 {
        return strings.TrimSpace(parts[0])
    }
    return ""
}

func printDependencyTree(dep *Dependency, indent int) {
    prefix := strings.Repeat("  ", indent)
    if indent > 0 {
        prefix += "└─ "
    }
    
    fmt.Printf("%s%s (%s)\n", prefix, dep.Name, dep.Version)
    
    for _, subDep := range dep.Dependencies {
        printDependencyTree(subDep, indent+1)
    }
}
```

---

**下一步**: 查看 [最佳实践](./best-practices.md) 了解性能优化和使用建议。
