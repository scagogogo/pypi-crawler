package main

import (
	"context"
	"fmt"
	"log"

	"github.com/scagogogo/pypi-crawler/pkg/pypi/mirrors"
)

func main() {
	// 创建默认API客户端 (使用官方PyPI)
	pypiClient := mirrors.NewOfficialClient()

	// 也可以使用预定义的镜像源
	// pypiClient = mirrors.NewTsinghuaClient() // 使用清华大学镜像
	// pypiClient = mirrors.NewDoubanClient()   // 使用豆瓣镜像

	// 或者自定义配置
	// options := client.NewOptions().
	//    WithTimeout(15 * time.Second).
	//    WithMaxRetries(5).
	//    WithProxy("http://your-proxy:8080")
	// pypiClient = client.NewClient(options)

	// 创建上下文
	ctx := context.Background()

	// 示例1: 获取所有包列表
	fmt.Println("正在获取包索引列表...")
	packageIndexes, err := pypiClient.GetAllPackages(ctx)
	if err != nil {
		log.Fatalf("获取包索引失败: %v", err)
	}
	fmt.Printf("共找到 %d 个包\n", len(packageIndexes))

	// 显示前10个包名
	fmt.Println("\n前10个包:")
	limit := 10
	if len(packageIndexes) < 10 {
		limit = len(packageIndexes)
	}
	for i := 0; i < limit; i++ {
		fmt.Printf("  %d. %s\n", i+1, packageIndexes[i])
	}

	// 示例2: 获取特定包的信息
	packageName := "requests"
	fmt.Printf("\n正在获取包 %s 的信息...\n", packageName)
	pkg, err := pypiClient.GetPackageInfo(ctx, packageName)
	if err != nil {
		log.Fatalf("获取包信息失败: %v", err)
	}

	fmt.Printf("包名: %s\n", pkg.Info.Name)
	fmt.Printf("版本: %s\n", pkg.Info.Version)
	fmt.Printf("摘要: %s\n", pkg.Info.Summary)
	fmt.Printf("主页: %s\n", pkg.Info.HomePage)
	fmt.Printf("作者: %s (%s)\n", pkg.Info.Author, pkg.Info.AuthorEmail)

	if pkg.Info.RequiresDist != nil && len(pkg.Info.RequiresDist) > 0 {
		fmt.Printf("依赖项 (%d):\n", len(pkg.Info.RequiresDist))
		for i, dep := range pkg.Info.RequiresDist {
			if i < 5 { // 只显示前5个依赖
				fmt.Printf("  - %s\n", dep)
			} else if i == 5 {
				fmt.Printf("  - ... 以及其他 %d 个依赖\n", len(pkg.Info.RequiresDist)-5)
				break
			}
		}
	}

	// 示例3: 获取特定版本的包信息
	version := "2.28.0" // 指定较老的版本
	fmt.Printf("\n正在获取包 %s 版本 %s 的信息...\n", packageName, version)
	versionPkg, err := pypiClient.GetPackageVersion(ctx, packageName, version)
	if err != nil {
		log.Fatalf("获取特定版本信息失败: %v", err)
	}

	fmt.Printf("包名: %s\n", versionPkg.Info.Name)
	fmt.Printf("版本: %s\n", versionPkg.Info.Version)
	if len(versionPkg.Urls) > 0 {
		fmt.Printf("发布时间: %s\n", versionPkg.Urls[0].UploadTime)
	}

	// 如果有漏洞信息，显示漏洞
	if len(versionPkg.Vulnerabilities) > 0 {
		fmt.Printf("已知漏洞数量: %d\n", len(versionPkg.Vulnerabilities))
		for i, vuln := range versionPkg.Vulnerabilities {
			fmt.Printf("  漏洞 #%d: %v\n", i+1, vuln.ID)
			fmt.Printf("     摘要: %s\n", vuln.Summary)
		}
	} else {
		fmt.Println("没有已知漏洞")
	}

	fmt.Println("\n示例运行完成！")
}
