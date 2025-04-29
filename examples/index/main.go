package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/scagogogo/pypi-crawler/pkg/pypi/mirrors"
)

func main() {
	// 创建一个带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// 选择要使用的镜像源
	// 默认使用官方源，但可以注释掉下面这行并取消注释另一个镜像源
	client := mirrors.NewOfficialClient()
	// client := mirrors.NewTsinghuaClient() // 使用清华大学镜像可能更快

	fmt.Println("正在从PyPI获取包索引...")
	fmt.Println("这可能需要几分钟时间，请耐心等待...")

	// 获取所有包列表
	indexList, err := client.GetAllPackages(ctx)
	if err != nil {
		fmt.Printf("获取索引列表失败: %v\n", err)
		os.Exit(1)
	}

	// 显示包的总数
	fmt.Printf("\nPyPI索引中共有 %d 个包\n", len(indexList))

	// 对包名进行排序，使输出更加一致
	sort.Strings(indexList)

	// 显示前10个包名（如果有的话）
	fmt.Println("\n按字母顺序排列的前10个包示例:")
	limit := 10
	if len(indexList) < limit {
		limit = len(indexList)
	}

	for i := 0; i < limit; i++ {
		fmt.Printf("  %d. %s\n", i+1, indexList[i])
	}

	// 也可以使用map形式获取索引
	fmt.Println("\n使用GetPackageList获取索引(Map形式):")
	packageMap, err := client.GetPackageList(ctx)
	if err != nil {
		fmt.Printf("获取索引列表(Map形式)失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Map形式的索引也包含 %d 个包\n", len(packageMap))

	// 检查一些常见包是否存在
	commonPackages := []string{"requests", "django", "flask", "numpy", "pandas"}
	fmt.Println("\n检查一些常见包是否存在:")
	for _, pkg := range commonPackages {
		_, exists := packageMap[pkg]
		fmt.Printf("  %s: %v\n", pkg, exists)
	}
}
