package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/crawler-go-go-go/go-requests"
	"github.com/scagogogo/pypi-crawler/pkg/model"
)

// Repository 表示一个存放着Pypi包的仓库
// 它提供了与PyPI包索引和包信息交互的方法
type Repository struct {
	options *Options
}

// NewRepository 创建一个新的PyPI仓库实例
//
// 参数:
//   - options: 可选的仓库选项，如果不提供，将使用默认选项（官方PyPI地址）
//
// 返回值:
//   - *Repository: 创建的仓库实例
//
// 使用示例:
//
//	// 创建默认仓库（使用官方PyPI）
//	repo := repository.NewRepository()
//
//	// 创建自定义仓库
//	options := repository.NewOptions().SetServerURL("https://mirrors.aliyun.com/pypi")
//	repo := repository.NewRepository(options)
func NewRepository(options ...*Options) *Repository {
	if len(options) == 0 {
		options = append(options, NewOptions())
	}
	return &Repository{
		options: options[0],
	}
}

// DownloadIndex 下载服务器上的所有包的索引文件
//
// 该方法会请求PyPI仓库的/simple端点，获取所有可用包的列表
//
// 参数:
//   - ctx: 上下文，用于控制请求的生命周期
//
// 返回值:
//   - model.PackageIndexes: 包名称的字符串切片
//   - error: 如果下载或解析过程中发生错误，则返回错误信息
//
// 使用示例:
//
//	ctx := context.Background()
//	repo := repository.NewRepository()
//	indexes, err := repo.DownloadIndex(ctx)
//	if err != nil {
//	    log.Fatalf("下载索引失败: %v", err)
//	}
//	fmt.Printf("共有 %d 个包\n", len(indexes))
//	// 输出前10个包名
//	for i := 0; i < 10 && i < len(indexes); i++ {
//	    fmt.Println(indexes[i])
//	}
//
// 响应数据示例:
//
//	["flask", "requests", "django", ...]
func (x *Repository) DownloadIndex(ctx context.Context) (model.PackageIndexes, error) {
	indexUrl := fmt.Sprintf("%s/simple", x.options.ServerURL)
	responseBody, err := x.getBytes(ctx, indexUrl)
	if err != nil {
		return nil, fmt.Errorf("download index from %s error: %s", indexUrl, err.Error())
	}

	// 有些镜像会修改标题，这里就不做检查了
	//if !strings.Contains(responseBody, "<title>Simple Index</title>") {
	//	return nil, fmt.Errorf("download index from %s error, invalid response: %s", indexUrl, responseBody)
	//}

	//_ = os.WriteFile("a.html", []byte(responseBody), os.ModePerm)

	return x.ParseIndexPage(string(responseBody))
}

// ParseIndexPage 解析包索引页面的HTML内容，提取包名称列表
//
// 该方法使用goquery库解析HTML，从中提取所有包的名称
//
// 参数:
//   - indexPageHtml: 需要解析的HTML页面内容字符串
//
// 返回值:
//   - model.PackageIndexes: 解析出的包名称列表
//   - error: 解析过程中的错误，如HTML格式不正确等
//
// 使用示例:
//
//	htmlContent := `<html><body>
//	  <a href="/simple/flask/">flask</a>
//	  <a href="/simple/requests/">requests</a>
//	</body></html>`
//	repo := repository.NewRepository()
//	packages, err := repo.ParseIndexPage(htmlContent)
//	if err != nil {
//	    log.Fatalf("解析失败: %v", err)
//	}
//	// 输出: [flask requests]
//	fmt.Println(packages)
func (x *Repository) ParseIndexPage(indexPageHtml string) (model.PackageIndexes, error) {
	document, err := goquery.NewDocumentFromReader(strings.NewReader(indexPageHtml))
	if err != nil {
		return nil, fmt.Errorf("page index page error, body = %s, error msg = %s", indexPageHtml, err.Error())
	}

	packageIndexes := make(model.PackageIndexes, 0)
	document.Find("body a").Each(func(i int, selection *goquery.Selection) {
		packageName := strings.TrimSpace(selection.Text())
		if packageName == "" {
			return
		}
		packageIndexes = append(packageIndexes, packageName)
	})
	return packageIndexes, nil
}

// GetPackage 获取Python包的详细信息
//
// 该方法通过请求PyPI的JSON API获取指定包的详细信息
//
// 参数:
//   - ctx: 上下文，用于控制请求的生命周期
//   - packageName: 要获取信息的包名称
//
// 返回值:
//   - *model.Package: 包含包详细信息的结构体指针
//   - error: 请求或解析过程中的错误
//
// 使用示例:
//
//	ctx := context.Background()
//	repo := repository.NewRepository()
//	pkg, err := repo.GetPackage(ctx, "requests")
//	if err != nil {
//	    log.Fatalf("获取包信息失败: %v", err)
//	}
//	fmt.Printf("包名: %s\n", pkg.Information.Name)
//	fmt.Printf("版本: %s\n", pkg.Information.Version)
//	fmt.Printf("描述: %s\n", pkg.Information.Description)
//
// 响应数据示例:
//
//	{
//	  "info": {
//	    "name": "requests",
//	    "version": "2.28.1",
//	    "summary": "Python HTTP for Humans",
//	    "description": "# Requests\n\n**Requests** is...",
//	    ...
//	  },
//	  "releases": { ... },
//	  "urls": [ ... ],
//	  ...
//	}
func (x *Repository) GetPackage(ctx context.Context, packageName string) (*model.Package, error) {
	targetUrl := fmt.Sprintf("%s/pypi/%s/json", x.options.ServerURL, url.PathEscape(packageName))
	packageBytes, err := x.getBytes(ctx, targetUrl)
	if err != nil {
		return nil, fmt.Errorf("request package %s information from %s error: %s", packageName, targetUrl, err.Error())
	}
	pkg := &model.Package{}
	err = json.Unmarshal(packageBytes, &pkg)
	if err != nil {
		return nil, fmt.Errorf("unmarshal package %s information response error, body = %s, error msg = %s", packageName, string(packageBytes), err.Error())
	}
	return pkg, nil
}

// 内部使用统一的方法来请求URL并返回字节数据
//
// 该方法处理HTTP请求的细节，包括代理设置等
//
// 参数:
//   - ctx: 请求上下文
//   - targetUrl: 目标URL字符串
//
// 返回值:
//   - []byte: 响应体的字节数据
//   - error: 请求过程中的错误
func (x *Repository) getBytes(ctx context.Context, targetUrl string) ([]byte, error) {
	options := requests.NewOptions[any, []byte](targetUrl, requests.BytesResponseHandler())
	if x.options.Proxy != "" {
		options.AppendRequestSetting(requests.RequestSettingProxy(x.options.Proxy))
	}
	return requests.SendRequest[any, []byte](ctx, options)
}
