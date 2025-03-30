package model

// PackageIndexes 表示PyPI仓库中所有包的索引列表
//
// 它是一个字符串切片，每个元素对应一个包的名称
//
// 使用示例:
//
//	repo := repository.NewRepository()
//	indexes, err := repo.DownloadIndex(context.Background())
//	if err == nil {
//	    fmt.Printf("共找到 %d 个包\n", len(indexes))
//
//	    // 查找特定包是否存在
//	    packageName := "requests"
//	    exists := false
//	    for _, name := range indexes {
//	        if name == packageName {
//	            exists = true
//	            break
//	        }
//	    }
//	    fmt.Printf("包 %s 是否存在: %v\n", packageName, exists)
//
//	    // 显示前10个包名
//	    for i := 0; i < 10 && i < len(indexes); i++ {
//	        fmt.Printf("包 #%d: %s\n", i+1, indexes[i])
//	    }
//	}
type PackageIndexes []string
