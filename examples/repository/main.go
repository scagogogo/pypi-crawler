package main

import (
	"fmt"
	"github.com/scagogogo/pypi-crawler/pkg/repository"
)

func main() {

	// 创建一个仓库
	r := repository.NewRepository()
	fmt.Println(r)

	// 创建仓库的时候可以指定一些选项
	// 仓库地址如果不指定的话使用的是官方的 https://pypi.org 这个地址
	// 代理地址是表示请求仓库的包时挂着代理IP去请求，这样如果速率快一些的话能够避免被封禁
	r = repository.NewRepository(repository.NewOptions().SetServerURL("https://pypi.org").SetProxy(""))

	// 内置了一些国内常见的pypi的镜像仓库，可以开箱即用不用再指定仓库地址
	// 比如连接到豆瓣的镜像源
	r = repository.NewDouBanRepository()

	// 其他的一些可选的镜像源
	r = repository.NewTencentCloudRepository()
	r = repository.NewUstcRepository()
	r = repository.NewNetEaseRepository()
	r = repository.NewTSingHuaRepository()

}
