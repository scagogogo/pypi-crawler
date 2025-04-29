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

	// 获取搜索关键词，如果没有提供，则使用默认值"python"
	keyword := "python"
	if len(os.Args) > 1 {
		keyword = os.Args[1]
	}

	// 获取搜索结果数量限制，默认为20
	limit := 20
	if len(os.Args) > 2 {
		if _, err := fmt.Sscanf(os.Args[2], "%d", &limit); err != nil {
			fmt.Printf("警告: 无法解析限制参数 '%s'，使用默认值 %d\n", os.Args[2], limit)
		}
	}

	// 创建PyPI客户端（使用官方源）
	client := mirrors.NewOfficialClient()

	// 搜索包
	fmt.Printf("正在搜索关键词 '%s' 的包（限制 %d 个结果）...\n", keyword, limit)
	fmt.Println("这可能需要一些时间，请耐心等待...")

	// 执行搜索
	startTime := time.Now()
	results, err := client.SearchPackages(ctx, keyword, limit)
	searchDuration := time.Since(startTime)

	if err != nil {
		fmt.Printf("搜索失败: %v\n", err)
		os.Exit(1)
	}

	// 显示搜索结果
	fmt.Printf("\n找到 %d 个匹配 '%s' 的包（搜索耗时 %.2f 秒）:\n", len(results), keyword, searchDuration.Seconds())
	for i, pkg := range results {
		fmt.Printf("%d. %s\n", i+1, pkg)
	}

	// 如果没有找到结果
	if len(results) == 0 {
		fmt.Println("未找到任何匹配的包，请尝试其他关键词。")
	} else if len(results) == limit {
		fmt.Printf("\n注：结果被限制为 %d 个，可能还有更多匹配项。\n", limit)
		fmt.Println("若要显示更多结果，可以指定更大的限制值：")
		fmt.Printf("  go run examples/search/main.go %s 50\n", keyword)
	}

	// 选择性地获取第一个结果的详细信息
	if len(results) > 0 && len(results) <= 5 {
		fmt.Printf("\n获取排名第一的包 '%s' 的详细信息...\n", results[0])
		pkg, err := client.GetPackageInfo(ctx, results[0])
		if err != nil {
			fmt.Printf("获取包详细信息失败: %v\n", err)
		} else {
			fmt.Println("\n包的基本信息:")
			fmt.Println("名称:", pkg.Info.Name)
			fmt.Println("版本:", pkg.Info.Version)
			fmt.Println("摘要:", pkg.Info.Summary)
			fmt.Println("作者:", pkg.Info.Author)
			fmt.Println("主页:", pkg.Info.HomePage)
		}
	}
}
