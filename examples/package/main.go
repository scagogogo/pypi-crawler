package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/scagogogo/pypi-crawler/pkg/pypi/mirrors"
)

func main() {
	// 创建一个带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 获取要查询的包名，默认为"requests"
	packageName := "requests"
	if len(os.Args) > 1 {
		packageName = os.Args[1]
	}

	// 创建PyPI客户端（使用官方源）
	client := mirrors.NewOfficialClient()

	// 获取名为"requests"的包的信息
	fmt.Printf("正在获取包 '%s' 的信息...\n", packageName)
	pkg, err := client.GetPackageInfo(ctx, packageName)
	if err != nil {
		fmt.Printf("获取包信息失败: %v\n", err)
		os.Exit(1)
	}

	// 显示包的基本信息
	fmt.Println("\n包的基本信息:")
	fmt.Println("包名:", pkg.Info.Name)
	fmt.Println("版本:", pkg.Info.Version)
	fmt.Println("摘要:", pkg.Info.Summary)
	fmt.Println("作者:", pkg.Info.Author)
	fmt.Println("作者邮箱:", pkg.Info.AuthorEmail)
	fmt.Println("主页:", pkg.Info.HomePage)
	fmt.Println("许可证:", pkg.Info.License)

	// 显示Python版本要求（如果有）
	if pkg.Info.HasPythonRequirement() {
		fmt.Println("Python版本要求:", pkg.Info.RequiresPython)
	}

	// 显示项目URL
	projectURLs := pkg.Info.GetProjectURLs()
	if len(projectURLs) > 0 {
		fmt.Println("\n项目链接:")
		for name, url := range projectURLs {
			fmt.Printf("  %s: %s\n", name, url)
		}
	}

	// 显示依赖项
	dependencies := pkg.Info.GetAllDependencies()
	if len(dependencies) > 0 {
		fmt.Printf("\n依赖项 (%d):\n", len(dependencies))
		for i, dep := range dependencies {
			if i < 15 {
				fmt.Printf("  %d. %s\n", i+1, dep)
			} else {
				fmt.Printf("  ...以及其他 %d 个依赖\n", len(dependencies)-15)
				break
			}
		}
	}

	// 获取包的所有版本
	fmt.Printf("\n正在获取 '%s' 的所有版本...\n", packageName)
	versions, err := client.GetPackageReleases(ctx, packageName)
	if err != nil {
		fmt.Printf("获取版本列表失败: %v\n", err)
	} else {
		fmt.Printf("共有 %d 个版本\n", len(versions))
		fmt.Println("最近的10个版本(如果有):")

		// 显示最近的几个版本
		limit := 10
		if len(versions) < limit {
			limit = len(versions)
		}

		for i := 0; i < limit; i++ {
			fmt.Printf("  %d. %s\n", i+1, versions[i])
		}
	}

	// 检查当前版本的漏洞
	fmt.Printf("\n检查 '%s' 版本 '%s' 的漏洞...\n", packageName, pkg.Info.Version)
	vulns, err := client.CheckPackageVulnerabilities(ctx, packageName, pkg.Info.Version)
	if err != nil {
		fmt.Printf("检查漏洞失败: %v\n", err)
	} else if len(vulns) == 0 {
		fmt.Println("未发现已知漏洞")
	} else {
		fmt.Printf("发现 %d 个漏洞:\n", len(vulns))
		for i, vuln := range vulns {
			fmt.Printf("  %d. ID: %s\n", i+1, vuln.ID)
			fmt.Printf("     摘要: %s\n", vuln.Summary)
			fmt.Printf("     详情: %s\n", vuln.Details)
			fmt.Printf("     链接: %s\n", vuln.Link)

			// 显示已修复的版本
			if len(vuln.FixedIn) > 0 {
				fmt.Printf("     已修复版本: %v\n", vuln.FixedIn)
			}

			// 检查是否有CVE编号
			if vuln.HasCVE() {
				fmt.Printf("     CVE编号: %v\n", vuln.GetCVEs())
			}

			fmt.Println()
		}
	}
}
