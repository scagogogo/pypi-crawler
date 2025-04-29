package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/scagogogo/pypi-crawler/pkg/pypi/mirrors"
)

func main() {
	// 创建客户端 - 使用官方镜像
	pypiClient := mirrors.NewOfficialClient()

	// 使用国内镜像源可以加速访问
	// pypiClient = mirrors.NewTsinghuaClient()
	// pypiClient = mirrors.NewAliyunClient()

	// 使用自定义选项
	// opts := client.NewOptions().
	//    WithTimeout(15 * time.Second).
	//    WithMaxRetries(5).
	//    WithProxy("http://127.0.0.1:8080")
	// pypiClient = client.NewClient(opts)

	// 创建上下文
	ctx := context.Background()

	// 示例1: 获取包信息
	packageName := "requests"
	fmt.Printf("获取包 %s 的信息...\n", packageName)
	pkg, err := pypiClient.GetPackageInfo(ctx, packageName)
	if err != nil {
		log.Fatalf("获取包信息失败: %v", err)
	}

	// 输出包的基本信息
	fmt.Printf("包名: %s\n", pkg.Info.Name)
	fmt.Printf("版本: %s\n", pkg.Info.Version)
	fmt.Printf("摘要: %s\n", pkg.Info.Summary)
	fmt.Printf("作者: %s (%s)\n", pkg.Info.Author, pkg.Info.AuthorEmail)
	fmt.Printf("许可证: %s\n", pkg.Info.License)

	// 示例2: 获取包的所有版本
	fmt.Printf("\n获取包 %s 的所有版本...\n", packageName)
	versions, err := pypiClient.GetPackageReleases(ctx, packageName)
	if err != nil {
		log.Fatalf("获取版本列表失败: %v", err)
	}

	// 最多显示10个版本
	maxVersions := 10
	if len(versions) < maxVersions {
		maxVersions = len(versions)
	}
	fmt.Printf("找到 %d 个版本，显示前 %d 个:\n", len(versions), maxVersions)
	for i := 0; i < maxVersions; i++ {
		fmt.Printf("  %d. %s\n", i+1, versions[i])
	}

	// 示例3: 获取特定版本的详细信息
	oldVersion := "2.25.0" // 请求一个旧版本
	fmt.Printf("\n获取包 %s 版本 %s 的信息...\n", packageName, oldVersion)
	versionPkg, err := pypiClient.GetPackageVersion(ctx, packageName, oldVersion)
	if err != nil {
		log.Fatalf("获取特定版本信息失败: %v", err)
	}

	fmt.Printf("包名: %s\n", versionPkg.Info.Name)
	fmt.Printf("版本: %s\n", versionPkg.Info.Version)

	// 获取发布文件信息
	if len(versionPkg.Urls) > 0 {
		fmt.Printf("发布文件: %d 个\n", len(versionPkg.Urls))
		for i, file := range versionPkg.Urls {
			if i < 3 { // 最多显示3个文件
				uploadTime, _ := file.GetUploadTimeISO()
				fmt.Printf("  %d. %s (%s, %d 字节, 上传于 %s)\n",
					i+1, file.Filename, file.PackageType, file.Size,
					uploadTime.Format(time.RFC3339))
			} else {
				fmt.Printf("  ... 以及其他 %d 个文件\n", len(versionPkg.Urls)-3)
				break
			}
		}
	}

	// 示例4: 检查包的漏洞
	fmt.Printf("\n检查包 %s 版本 %s 的漏洞...\n", packageName, oldVersion)
	vulnerabilities, err := pypiClient.CheckPackageVulnerabilities(ctx, packageName, oldVersion)
	if err != nil {
		log.Fatalf("检查漏洞失败: %v", err)
	}

	if len(vulnerabilities) > 0 {
		fmt.Printf("发现 %d 个漏洞:\n", len(vulnerabilities))
		for i, vuln := range vulnerabilities {
			fmt.Printf("  %d. [%s] %s\n", i+1, vuln.ID, vuln.Summary)
			if len(vuln.FixedIn) > 0 {
				fmt.Printf("     已在以下版本修复: %v\n", vuln.FixedIn)
			}
			if vuln.IsWithdrawn() {
				fmt.Printf("     注意: 此漏洞报告已被撤回\n")
			}
		}
	} else {
		fmt.Printf("未发现已知漏洞\n")
	}
}
