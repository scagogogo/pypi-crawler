# PyPi Crawler

# 一、这是什么？

这是一个pypi的爬虫库，能够让你获取pypi上的包的信息。

# 二、安装依赖

```bash
go get -u github.com/scagogogo/pypi-crawler
```

# 三、API示例

## 3.1 创建仓库

首先应该创建一个仓库，表示是从哪个仓库中获取包的信息，仓库又分为官方的`https://pypi.org`或其它镜像仓库：

```go
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
```

## 2.2 下载索引

```go
package main

import (
	"context"
	"fmt"
	"github.com/scagogogo/pypi-crawler/pkg/repository"
)

func main() {

	// 下载仓库中所有的包的索引目录
	r := repository.NewUstcRepository()
	index, err := r.DownloadIndex(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println(fmt.Sprintf("仓库中有 %d 个包", len(index)))

	// Output:
	// 仓库中有 515608 个包

}
```

## 2.3  获取包信息

 ```go
 package main
 
 import (
 	"context"
 	"fmt"
 	"github.com/scagogogo/pypi-crawler/pkg/repository"
 )
 
 func main() {
 
 	// 获取名为requests的包的信息
 	r := repository.NewRepository()
 	packageInformation, err := r.GetPackage(context.Background(), "requests")
 	if err != nil {
 		panic(err)
 	}
 	fmt.Println(packageInformation.Information.Name)
 	fmt.Println(packageInformation.Information.Version)
 	fmt.Println(packageInformation.Information.Description)
 
 	// Output:
 	// requests
 	//2.31.0
 	//# Requests
 	//
 	//**Requests** is a simple, yet elegant, HTTP library.
 	// ...
 
 }
 ```

