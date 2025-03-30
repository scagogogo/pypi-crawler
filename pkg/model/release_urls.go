package model

// ReleaseUrls 表示一个包版本的多个发布文件集合
//
// 这是一个ReleaseURL指针的切片，每个元素代表同一版本的不同发布文件
// 通常包含不同的文件格式（如wheel、源码包等）或针对不同平台的构建版本
//
// 使用示例:
//
//	pkg, err := repo.GetPackage(context.Background(), "requests")
//	if err == nil {
//	    urls := pkg.Urls
//	    fmt.Printf("最新版本共有 %d 个发布文件\n", urls.Len())
//
//	    // 遍历所有发布文件
//	    for i, releaseURL := range urls {
//	        fmt.Printf("文件 #%d: %s (%s)\n", i+1, releaseURL.Filename, releaseURL.Packagetype)
//	    }
//
//	    // 查找wheel格式的包
//	    for _, releaseURL := range urls {
//	        if releaseURL.Packagetype == "bdist_wheel" {
//	            fmt.Printf("找到wheel包: %s\n", releaseURL.Filename)
//	            break
//	        }
//	    }
//	}
type ReleaseUrls []*ReleaseURL

// Len 返回发布文件集合中的文件数量
//
// 返回值:
//   - int: 集合中发布文件的数量
//
// 使用示例:
//
//	urls := pkg.Urls
//	fmt.Printf("共有 %d 个发布文件\n", urls.Len())
func (x ReleaseUrls) Len() int {
	return len(x)
}
