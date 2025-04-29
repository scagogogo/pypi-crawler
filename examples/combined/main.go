package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/scagogogo/pypi-crawler/pkg/pypi/client"
	"github.com/scagogogo/pypi-crawler/pkg/pypi/mirrors"
)

// 不同的镜像源选项
const (
	MirrorOfficial = "official"
	MirrorTsinghua = "tsinghua"
	MirrorDouban   = "douban"
	MirrorAliyun   = "aliyun"
	MirrorTencent  = "tencent"
	MirrorUstc     = "ustc"
	MirrorNetease  = "netease"
)

func main() {
	// 如果没有足够的参数，显示使用说明
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]

	// 解析可选的镜像源
	mirrorSource := MirrorOfficial
	timeoutSeconds := 60

	// 检查是否有--mirror标志
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "--mirror" && i+1 < len(os.Args) {
			mirrorSource = os.Args[i+1]
			// 从参数列表中移除这两个参数
			os.Args = append(os.Args[:i], os.Args[i+2:]...)
			break
		}
	}

	// 检查是否有--timeout标志
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "--timeout" && i+1 < len(os.Args) {
			if _, err := fmt.Sscanf(os.Args[i+1], "%d", &timeoutSeconds); err != nil {
				fmt.Printf("警告: 无法解析超时值 '%s'，使用默认值 %d 秒\n", os.Args[i+1], timeoutSeconds)
			}
			// 从参数列表中移除这两个参数
			os.Args = append(os.Args[:i], os.Args[i+2:]...)
			break
		}
	}

	// 创建上下文（带超时）
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	// 创建PyPI客户端
	pypiClient := createClient(mirrorSource)
	fmt.Printf("使用 %s 镜像源\n\n", getMirrorName(mirrorSource))

	switch strings.ToLower(command) {
	case "info":
		handleInfoCommand(ctx, pypiClient, os.Args[2:])
	case "version":
		handleVersionCommand(ctx, pypiClient, os.Args[2:])
	case "releases":
		handleReleasesCommand(ctx, pypiClient, os.Args[2:])
	case "search":
		handleSearchCommand(ctx, pypiClient, os.Args[2:])
	case "list":
		handleListCommand(ctx, pypiClient, os.Args[2:])
	case "check":
		handleCheckCommand(ctx, pypiClient, os.Args[2:])
	case "mirrors":
		handleMirrorsCommand(os.Args[2:])
	default:
		fmt.Printf("错误: 未知命令 '%s'\n", command)
		printUsage()
		os.Exit(1)
	}
}

// 显示使用帮助
func printUsage() {
	fmt.Println("使用: go run examples/combined/main.go <command> [arguments...] [options]")
	fmt.Println("\n可用命令:")
	fmt.Println("  info <package>              - 获取包信息")
	fmt.Println("  version <package> <version> - 获取包的特定版本信息")
	fmt.Println("  releases <package>          - 获取包的所有版本")
	fmt.Println("  search <keyword> [limit]    - 搜索包")
	fmt.Println("  list [limit]                - 列出所有包（显示一部分）")
	fmt.Println("  check <package> <version>   - 检查包的漏洞")
	fmt.Println("  mirrors                     - 显示所有支持的镜像源")

	fmt.Println("\n全局选项:")
	fmt.Println("  --mirror <source>           - 指定镜像源 (默认: official)")
	fmt.Println("                                可选: official, tsinghua, douban, aliyun, tencent, ustc, netease")
	fmt.Println("  --timeout <seconds>         - 设置操作超时时间 (默认: 60)")

	fmt.Println("\n示例:")
	fmt.Println("  go run examples/combined/main.go info requests")
	fmt.Println("  go run examples/combined/main.go search flask --mirror tsinghua")
	fmt.Println("  go run examples/combined/main.go check django 3.2.0 --timeout 30")
}

// 根据选择创建客户端
func createClient(mirrorSource string) client.Client {
	// 设置自定义选项
	options := client.NewOptions().
		WithUserAgent("PyPI-Crawler-Example/1.0").
		WithMaxRetries(2)

	// 根据镜像源选择创建不同的客户端
	switch strings.ToLower(mirrorSource) {
	case MirrorTsinghua:
		return mirrors.NewTsinghuaClient(options)
	case MirrorDouban:
		return mirrors.NewDoubanClient(options)
	case MirrorAliyun:
		return mirrors.NewAliyunClient(options)
	case MirrorTencent:
		return mirrors.NewTencentClient(options)
	case MirrorUstc:
		return mirrors.NewUstcClient(options)
	case MirrorNetease:
		return mirrors.NewNeteaseClient(options)
	default:
		return mirrors.NewOfficialClient(options)
	}
}

// 获取镜像源名称
func getMirrorName(mirrorSource string) string {
	switch strings.ToLower(mirrorSource) {
	case MirrorTsinghua:
		return "清华大学"
	case MirrorDouban:
		return "豆瓣"
	case MirrorAliyun:
		return "阿里云"
	case MirrorTencent:
		return "腾讯云"
	case MirrorUstc:
		return "中国科技大学"
	case MirrorNetease:
		return "网易"
	default:
		return "官方PyPI"
	}
}

// 处理info命令
func handleInfoCommand(ctx context.Context, pypiClient client.Client, args []string) {
	if len(args) < 1 {
		fmt.Println("错误: 缺少包名参数")
		fmt.Println("用法: go run examples/combined/main.go info <package>")
		os.Exit(1)
	}
	packageName := args[0]

	fmt.Printf("正在获取包 '%s' 的信息...\n", packageName)
	pkg, err := pypiClient.GetPackageInfo(ctx, packageName)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}

	// 显示包信息
	fmt.Println("\n包详情:")
	fmt.Println("名称:", pkg.Info.Name)
	fmt.Println("版本:", pkg.Info.Version)
	fmt.Println("摘要:", pkg.Info.Summary)
	fmt.Println("作者:", pkg.Info.Author)
	fmt.Println("主页:", pkg.Info.HomePage)
	fmt.Println("许可证:", pkg.Info.License)

	// 显示依赖项
	dependencies := pkg.Info.GetAllDependencies()
	if len(dependencies) > 0 {
		fmt.Printf("\n依赖项 (%d):\n", len(dependencies))
		for i, dep := range dependencies {
			if i < 10 {
				fmt.Printf("  %d. %s\n", i+1, dep)
			} else {
				fmt.Printf("  ... 以及其他 %d 个依赖\n", len(dependencies)-10)
				break
			}
		}
	}

	// 显示项目URL
	urls := pkg.Info.GetProjectURLs()
	if len(urls) > 0 {
		fmt.Println("\n项目链接:")
		for name, url := range urls {
			fmt.Printf("  %s: %s\n", name, url)
		}
	}
}

// 处理version命令
func handleVersionCommand(ctx context.Context, pypiClient client.Client, args []string) {
	if len(args) < 2 {
		fmt.Println("错误: 需要包名和版本号")
		fmt.Println("用法: go run examples/combined/main.go version <package> <version>")
		os.Exit(1)
	}
	packageName := args[0]
	version := args[1]

	fmt.Printf("正在获取包 '%s' 版本 '%s' 的信息...\n", packageName, version)
	pkg, err := pypiClient.GetPackageVersion(ctx, packageName, version)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}

	// 显示版本信息
	fmt.Println("\n版本详情:")
	fmt.Println("名称:", pkg.Info.Name)
	fmt.Println("版本:", pkg.Info.Version)
	fmt.Println("发布文件:")
	for i, file := range pkg.Urls {
		if i < 3 {
			fmt.Printf("  %d. %s (%s, %d bytes)\n", i+1, file.Filename, file.PackageType, file.Size)
		} else {
			fmt.Printf("  ... 以及其他 %d 个文件\n", len(pkg.Urls)-3)
			break
		}
	}
}

// 处理releases命令
func handleReleasesCommand(ctx context.Context, pypiClient client.Client, args []string) {
	if len(args) < 1 {
		fmt.Println("错误: 缺少包名参数")
		fmt.Println("用法: go run examples/combined/main.go releases <package>")
		os.Exit(1)
	}
	packageName := args[0]

	fmt.Printf("正在获取包 '%s' 的所有版本...\n", packageName)
	versions, err := pypiClient.GetPackageReleases(ctx, packageName)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}

	// 显示版本列表
	fmt.Printf("\n找到 %d 个版本:\n", len(versions))
	limit := 15
	if len(versions) < limit {
		limit = len(versions)
	}
	for i := 0; i < limit; i++ {
		fmt.Printf("  %d. %s\n", i+1, versions[i])
	}
	if len(versions) > limit {
		fmt.Printf("  ... 以及其他 %d 个版本\n", len(versions)-limit)
	}
}

// 处理search命令
func handleSearchCommand(ctx context.Context, pypiClient client.Client, args []string) {
	if len(args) < 1 {
		fmt.Println("错误: 缺少搜索关键词")
		fmt.Println("用法: go run examples/combined/main.go search <keyword> [limit]")
		os.Exit(1)
	}
	keyword := args[0]

	// 获取可选的限制参数
	limit := 10
	if len(args) > 1 {
		if _, err := fmt.Sscanf(args[1], "%d", &limit); err != nil {
			fmt.Printf("警告: 无法解析限制参数 '%s'，使用默认值 %d\n", args[1], limit)
		}
	}

	fmt.Printf("正在搜索关键词 '%s' 的包...\n", keyword)
	startTime := time.Now()
	results, err := pypiClient.SearchPackages(ctx, keyword, limit)
	searchDuration := time.Since(startTime)

	if err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}

	// 显示搜索结果
	fmt.Printf("\n找到 %d 个匹配 '%s' 的包（搜索耗时 %.2f 秒）:\n", len(results), keyword, searchDuration.Seconds())
	for i, pkg := range results {
		fmt.Printf("%d. %s\n", i+1, pkg)
	}

	if len(results) == limit {
		fmt.Printf("\n注: 结果已限制为前 %d 个，可能有更多匹配项\n", limit)
	}
}

// 处理list命令
func handleListCommand(ctx context.Context, pypiClient client.Client, args []string) {
	// 获取可选的限制参数
	limit := 15
	if len(args) > 0 {
		if _, err := fmt.Sscanf(args[0], "%d", &limit); err != nil {
			fmt.Printf("警告: 无法解析限制参数 '%s'，使用默认值 %d\n", args[0], limit)
		}
	}

	fmt.Println("正在获取包列表...")
	startTime := time.Now()
	packages, err := pypiClient.GetAllPackages(ctx)
	listDuration := time.Since(startTime)

	if err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}

	// 显示包列表
	fmt.Printf("\nPyPI 包索引中共有 %d 个包（获取耗时 %.2f 秒）\n", len(packages), listDuration.Seconds())
	fmt.Printf("\n示例（前%d个）:\n", limit)

	// 显示前N个包
	showLimit := limit
	if len(packages) < showLimit {
		showLimit = len(packages)
	}
	for i := 0; i < showLimit; i++ {
		fmt.Printf("  %d. %s\n", i+1, packages[i])
	}
}

// 处理check命令
func handleCheckCommand(ctx context.Context, pypiClient client.Client, args []string) {
	if len(args) < 2 {
		fmt.Println("错误: 需要包名和版本号")
		fmt.Println("用法: go run examples/combined/main.go check <package> <version>")
		os.Exit(1)
	}
	packageName := args[0]
	version := args[1]

	fmt.Printf("正在检查包 '%s' 版本 '%s' 的漏洞...\n", packageName, version)
	vulnerabilities, err := pypiClient.CheckPackageVulnerabilities(ctx, packageName, version)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}

	// 显示漏洞信息
	if len(vulnerabilities) == 0 {
		fmt.Println("\n未发现任何已知漏洞")
	} else {
		fmt.Printf("\n发现 %d 个漏洞:\n", len(vulnerabilities))
		for i, vuln := range vulnerabilities {
			fmt.Printf("%d. ID: %s\n", i+1, vuln.ID)
			fmt.Printf("   摘要: %s\n", vuln.Summary)
			fmt.Printf("   详情: %s\n", vuln.Details)
			fmt.Printf("   链接: %s\n", vuln.Link)

			// 显示已修复的版本
			if len(vuln.FixedIn) > 0 {
				fmt.Printf("   修复版本: %v\n", vuln.FixedIn)
			}

			// 显示CVE信息（如果有）
			if vuln.HasCVE() {
				fmt.Printf("   CVE编号: %v\n", vuln.GetCVEs())
			}

			fmt.Println()
		}
	}
}

// 处理mirrors命令
func handleMirrorsCommand(args []string) {
	fmt.Println("支持的PyPI镜像源:")
	fmt.Printf("  %-12s: %s\n", "official", mirrors.OfficialURL)
	fmt.Printf("  %-12s: %s\n", "tsinghua", mirrors.TsinghuaURL)
	fmt.Printf("  %-12s: %s\n", "douban", mirrors.DoubanURL)
	fmt.Printf("  %-12s: %s\n", "aliyun", mirrors.AliyunURL)
	fmt.Printf("  %-12s: %s\n", "tencent", mirrors.TencentURL)
	fmt.Printf("  %-12s: %s\n", "ustc", mirrors.UstcURL)
	fmt.Printf("  %-12s: %s\n", "netease", mirrors.NeteaseURL)

	fmt.Println("\n使用示例:")
	fmt.Println("  go run examples/combined/main.go search flask --mirror tsinghua")
}
