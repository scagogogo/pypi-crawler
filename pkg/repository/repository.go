package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/crawler-go-go-go/go-requests"
	"github.com/scagogogo/pypi-crawler/pkg/model"
	"net/url"
	"strings"
)

// Repository 表示一个存放着Pypi包的仓库
type Repository struct {
	options *Options
}

func NewRepository(options ...*Options) *Repository {
	if len(options) == 0 {
		options = append(options, NewOptions())
	}
	return &Repository{
		options: options[0],
	}
}

// DownloadIndex 下载服务器上的所有的包的索引文件
func (x *Repository) DownloadIndex(ctx context.Context) (model.PackageIndexes, error) {
	indexUrl := fmt.Sprintf("%s/simple", x.options.ServerURL)
	responseBody, err := requests.GetString(ctx, indexUrl)
	if err != nil {
		return nil, fmt.Errorf("download index from %s error: %s", indexUrl, err.Error())
	}

	// 有些镜像会修改标题，这里就不做检查了
	//if !strings.Contains(responseBody, "<title>Simple Index</title>") {
	//	return nil, fmt.Errorf("download index from %s error, invalid response: %s", indexUrl, responseBody)
	//}

	//_ = os.WriteFile("a.html", []byte(responseBody), os.ModePerm)

	return x.ParseIndexPage(responseBody)
}

// ParseIndexPage 解析包索引页面
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

// GetPackage 获取Python中的包的信息
func (x *Repository) GetPackage(ctx context.Context, packageName string) (*model.Package, error) {
	targetUrl := fmt.Sprintf("%s/pypi/%s/json", x.options.ServerURL, url.PathEscape(packageName))
	packageBytes, err := requests.GetBytes(ctx, targetUrl)
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
