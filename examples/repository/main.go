package main

import (
	"fmt"

	"github.com/scagogogo/pypi-crawler/pkg/pypi/client"
	"github.com/scagogogo/pypi-crawler/pkg/pypi/mirrors"
)

func main() {
	// 使用不同的镜像源创建PyPI客户端

	// 创建默认客户端（使用官方PyPI）
	officialClient := mirrors.NewOfficialClient()
	fmt.Println("官方PyPI客户端:", officialClient)

	// 使用国内镜像源

	// 清华大学镜像
	tsinghuaClient := mirrors.NewTsinghuaClient()
	fmt.Println("清华大学镜像客户端:", tsinghuaClient)

	// 豆瓣镜像
	doubanClient := mirrors.NewDoubanClient()
	fmt.Println("豆瓣镜像客户端:", doubanClient)

	// 阿里云镜像
	aliyunClient := mirrors.NewAliyunClient()
	fmt.Println("阿里云镜像客户端:", aliyunClient)

	// 腾讯云镜像
	tencentClient := mirrors.NewTencentClient()
	fmt.Println("腾讯云镜像客户端:", tencentClient)

	// 中国科技大学镜像
	ustcClient := mirrors.NewUstcClient()
	fmt.Println("中国科技大学镜像客户端:", ustcClient)

	// 网易镜像
	neteaseClient := mirrors.NewNeteaseClient()
	fmt.Println("网易镜像客户端:", neteaseClient)

	// 使用自定义选项
	fmt.Println("\n创建带自定义选项的客户端:")

	customOptions := client.NewOptions().
		WithUserAgent("PyPI-Crawler/1.0").
		WithTimeout(30).
		WithMaxRetries(3)

	customClient := mirrors.NewOfficialClient(customOptions)
	fmt.Println("自定义客户端:", customClient)

	// 显示所有支持的镜像源URL
	fmt.Println("\n支持的镜像源URL:")
	fmt.Println("官方PyPI:", mirrors.OfficialURL)
	fmt.Println("清华大学:", mirrors.TsinghuaURL)
	fmt.Println("豆瓣:", mirrors.DoubanURL)
	fmt.Println("阿里云:", mirrors.AliyunURL)
	fmt.Println("腾讯云:", mirrors.TencentURL)
	fmt.Println("中国科技大学:", mirrors.UstcURL)
	fmt.Println("网易:", mirrors.NeteaseURL)
}
