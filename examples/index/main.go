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
