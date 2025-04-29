# PyPi Crawler

# 一、这是什么？

这是一个pypi的爬虫库，能够让你获取pypi上的包的信息。本库提供了完整的PyPI访问接口，支持多种镜像源，并有丰富的配置选项。

# 二、安装依赖

```bash
go get -u github.com/scagogogo/pypi-crawler
```

# 三、API使用指南

## 3.1 创建客户端

首先需要创建一个PyPI客户端实例：

```go
package main

import (
	"fmt"
	"github.com/scagogogo/pypi-crawler/pkg/pypi/mirrors"
	"github.com/scagogogo/pypi-crawler/pkg/pypi/client"
)

func main() {
	// 创建一个使用官方PyPI源的客户端
	pypiClient := mirrors.NewOfficialClient()
	fmt.Println("已创建官方PyPI客户端")
	
	// 也可以使用国内镜像源以提高访问速度
	// 清华大学镜像
	tsinghuaClient := mirrors.NewTsinghuaClient()
	
	// 豆瓣镜像
	doubanClient := mirrors.NewDoubanClient()
	
	// 阿里云镜像
	aliyunClient := mirrors.NewAliyunClient()
	
	// 其他可用镜像源
	tencentClient := mirrors.NewTencentClient()  // 腾讯云镜像
	ustcClient := mirrors.NewUstcClient()        // 中国科技大学镜像
	neteaseClient := mirrors.NewNeteaseClient()  // 网易镜像
	
	// 可以使用自定义选项创建客户端
	customOptions := client.NewOptions().
		WithUserAgent("MyPyPIClient/1.0").     // 设置User-Agent
		WithTimeout(30).                       // 设置超时时间(秒)
		WithMaxRetries(3).                     // 设置最大重试次数
		WithProxy("http://127.0.0.1:8080")     // 设置HTTP代理
	
	customClient := mirrors.NewOfficialClient(customOptions)
	fmt.Println("已创建自定义客户端")
}
```

## 3.2 获取包信息

获取特定包的最新版本信息：

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
	pypiClient := mirrors.NewOfficialClient()
	
	// 创建上下文
	ctx := context.Background()
	
	// 获取包信息
	packageName := "requests"
	pkg, err := pypiClient.GetPackageInfo(ctx, packageName)
	if err != nil {
		log.Fatalf("获取包信息失败: %v", err)
	}
	
	// 访问包的基本信息
	fmt.Printf("包名: %s\n", pkg.Info.Name)
	fmt.Printf("版本: %s\n", pkg.Info.Version)
	fmt.Printf("摘要: %s\n", pkg.Info.Summary)
	fmt.Printf("作者: %s (%s)\n", pkg.Info.Author, pkg.Info.AuthorEmail)
	fmt.Printf("许可证: %s\n", pkg.Info.License)
	
	// 获取Python版本要求
	if pkg.Info.HasPythonRequirement() {
		fmt.Printf("Python版本要求: %s\n", pkg.Info.RequiresPython)
	}
	
	// 获取项目URL
	projectURLs := pkg.Info.GetProjectURLs()
	if len(projectURLs) > 0 {
		fmt.Println("\n项目链接:")
		for name, url := range projectURLs {
			fmt.Printf("  %s: %s\n", name, url)
		}
	}
	
	// 获取依赖项
	dependencies := pkg.Info.GetAllDependencies()
	if len(dependencies) > 0 {
		fmt.Printf("\n依赖项 (%d):\n", len(dependencies))
		for i, dep := range dependencies {
			if i < 10 {
				fmt.Printf("  %d. %s\n", i+1, dep)
			} else {
				fmt.Printf("  ...以及其他 %d 个依赖\n", len(dependencies)-10)
				break
			}
		}
	}
}
```

## 3.3 获取特定版本信息

获取包的特定版本信息：

```go
package main

import (
	"context"
	"fmt"
	"log"
	"github.com/scagogogo/pypi-crawler/pkg/pypi/mirrors"
)

func main() {
	pypiClient := mirrors.NewOfficialClient()
	ctx := context.Background()
	
	packageName := "requests"
	version := "2.25.0"  // 指定版本号
	
	versionPkg, err := pypiClient.GetPackageVersion(ctx, packageName, version)
	if err != nil {
		log.Fatalf("获取版本信息失败: %v", err)
	}
	
	fmt.Printf("包名: %s\n", versionPkg.Info.Name)
	fmt.Printf("版本: %s\n", versionPkg.Info.Version)
	
	// 获取发布文件信息
	if len(versionPkg.Urls) > 0 {
		fmt.Printf("发布文件: %d 个\n", len(versionPkg.Urls))
		for i, file := range versionPkg.Urls {
			if i < 3 { // 最多显示3个文件
				fmt.Printf("  %d. %s (%s, %d 字节)\n",
					i+1, file.Filename, file.PackageType, file.Size)
			} else {
				fmt.Printf("  ... 以及其他 %d 个文件\n", len(versionPkg.Urls)-3)
				break
			}
		}
	}
}
```

## 3.4 获取包的所有版本

列出一个包的所有发布版本：

```go
package main

import (
	"context"
	"fmt"
	"log"
	"github.com/scagogogo/pypi-crawler/pkg/pypi/mirrors"
)

func main() {
	pypiClient := mirrors.NewOfficialClient()
	ctx := context.Background()
	
	packageName := "requests"
	
	versions, err := pypiClient.GetPackageReleases(ctx, packageName)
	if err != nil {
		log.Fatalf("获取版本列表失败: %v", err)
	}
	
	fmt.Printf("包 %s 共有 %d 个版本\n", packageName, len(versions))
	
	// 显示最近的10个版本
	maxVersions := 10
	if len(versions) < maxVersions {
		maxVersions = len(versions)
	}
	
	fmt.Printf("显示前 %d 个版本:\n", maxVersions)
	for i := 0; i < maxVersions; i++ {
		fmt.Printf("  %d. %s\n", i+1, versions[i])
	}
}
```

## 3.5 检查包的漏洞

检查特定版本包的已知漏洞：

```go
package main

import (
	"context"
	"fmt"
	"log"
	"github.com/scagogogo/pypi-crawler/pkg/pypi/mirrors"
)

func main() {
	pypiClient := mirrors.NewOfficialClient()
	ctx := context.Background()
	
	packageName := "requests"
	version := "2.25.0"  // 检查这个版本
	
	vulnerabilities, err := pypiClient.CheckPackageVulnerabilities(ctx, packageName, version)
	if err != nil {
		log.Fatalf("检查漏洞失败: %v", err)
	}
	
	if len(vulnerabilities) == 0 {
		fmt.Println("未发现已知漏洞")
	} else {
		fmt.Printf("发现 %d 个漏洞:\n", len(vulnerabilities))
		for i, vuln := range vulnerabilities {
			fmt.Printf("  %d. [%s] %s\n", i+1, vuln.ID, vuln.Summary)
			
			if len(vuln.FixedIn) > 0 {
				fmt.Printf("     已在以下版本修复: %v\n", vuln.FixedIn)
			}
			
			if vuln.HasCVE() {
				fmt.Printf("     CVE编号: %v\n", vuln.GetCVEs())
			}
		}
	}
}
```

## 3.6 获取所有包列表

获取PyPI索引中的所有包名：

```go
package main

import (
	"context"
	"fmt"
	"log"
	"github.com/scagogogo/pypi-crawler/pkg/pypi/mirrors"
)

func main() {
	pypiClient := mirrors.NewOfficialClient()
	ctx := context.Background()
	
	// 注意：此操作可能需要较长时间
	packages, err := pypiClient.GetAllPackages(ctx)
	if err != nil {
		log.Fatalf("获取包索引失败: %v", err)
	}
	
	fmt.Printf("PyPI中共有 %d 个包\n", len(packages))
	
	// 显示前10个包名
	limit := 10
	if len(packages) < limit {
		limit = len(packages)
	}
	
	fmt.Println("前10个包名示例:")
	for i := 0; i < limit; i++ {
		fmt.Printf("  %d. %s\n", i+1, packages[i])
	}
}
```

## 3.7 搜索包

通过关键词搜索包：

```go
package main

import (
	"context"
	"fmt"
	"log"
	"github.com/scagogogo/pypi-crawler/pkg/pypi/mirrors"
)

func main() {
	pypiClient := mirrors.NewOfficialClient()
	ctx := context.Background()
	
	keyword := "flask"
	limit := 10  // 限制结果数量
	
	results, err := pypiClient.SearchPackages(ctx, keyword, limit)
	if err != nil {
		log.Fatalf("搜索失败: %v", err)
	}
	
	fmt.Printf("搜索关键词 '%s' 找到 %d 个结果:\n", keyword, len(results))
	for i, pkg := range results {
		fmt.Printf("  %d. %s\n", i+1, pkg)
	}
}
```

# 四、更多示例

更多详细示例可以参考项目的`examples`目录，包括：

- `examples/pypi_client` - PyPI客户端基本使用示例
- `examples/api_client` - API客户端使用示例
- `examples/index` - 获取包索引示例
- `examples/package` - 获取包信息示例
- `examples/repository` - 不同镜像源创建示例
- `examples/search` - 搜索包示例
- `examples/combined` - 综合功能命令行工具示例

# 五、项目结构

```
pkg/pypi/
├── api/            - API 接口定义
├── client/         - API 实现
├── mirrors/        - 镜像源工厂
├── models/         - 数据模型
```

# 六、贡献指南

欢迎提交Issue和Pull Request。在提交PR前，请确保：

1. 代码通过所有测试
2. 新功能已添加相应的测试
3. 文档已更新

# 七、许可证

本项目使用MIT许可证。详情请查看LICENSE文件。

