package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/scagogogo/pypi-crawler/pkg/pypi/api"
	"github.com/scagogogo/pypi-crawler/pkg/pypi/models"
)

// Client 实现了与PyPI JSON API交互的客户端
// 使用纯API调用方式，不使用爬虫方式
type Client struct {
	options *Options
	client  *http.Client
}

// NewClient 创建一个新的PyPI客户端实例
//
// 参数:
//   - options: 可选配置，如API基址、超时、代理等
//
// 返回值:
//   - api.PyPIClient: 客户端接口实例
//
// 使用示例:
//
//	// 使用默认配置创建客户端
//	client := client.NewClient()
//
//	// 使用自定义配置创建客户端
//	options := client.NewOptions().WithBaseURL("https://pypi.org").WithTimeout(10 * time.Second)
//	client := client.NewClient(options)
func NewClient(options ...*Options) api.PyPIClient {
	// 使用默认选项或传入的选项
	var clientOptions *Options
	if len(options) == 0 {
		clientOptions = NewOptions()
	} else {
		clientOptions = options[0]
	}

	// 创建HTTP客户端
	transport := http.DefaultTransport.(*http.Transport).Clone()

	// 设置代理
	if clientOptions.Proxy != "" {
		proxyURL, _ := url.Parse(clientOptions.Proxy)
		transport.Proxy = http.ProxyURL(proxyURL)
	}

	httpClient := &http.Client{
		Transport: transport,
		Timeout:   clientOptions.Timeout,
	}

	return &Client{
		options: clientOptions,
		client:  httpClient,
	}
}

// NewClientWithOptions 使用指定选项创建客户端（供内部和兼容层使用）
//
// 参数:
//   - options: 配置选项，如为nil则使用默认值
//
// 返回值:
//   - api.PyPIClient: 客户端接口实例
func NewClientWithOptions(options *Options) api.PyPIClient {
	if options == nil {
		options = NewOptions()
	}

	// 创建HTTP客户端
	transport := http.DefaultTransport.(*http.Transport).Clone()

	// 设置代理
	if options.Proxy != "" {
		proxyURL, _ := url.Parse(options.Proxy)
		transport.Proxy = http.ProxyURL(proxyURL)
	}

	httpClient := &http.Client{
		Transport: transport,
		Timeout:   options.Timeout,
	}

	return &Client{
		options: options,
		client:  httpClient,
	}
}

// GetPackageInfo 获取指定包的最新版本信息
//
// 参数:
//   - ctx: 上下文，用于控制请求的生命周期
//   - packageName: 要获取信息的包名
//
// 返回值:
//   - *models.Package: 包含包详细信息的结构体指针
//   - error: 如有错误则返回，否则为nil
func (c *Client) GetPackageInfo(ctx context.Context, packageName string) (*models.Package, error) {
	// 构建API URL
	apiURL := fmt.Sprintf("%s/pypi/%s/json", c.options.BaseURL, url.PathEscape(packageName))

	// 发送请求
	responseBody, err := c.sendRequest(ctx, apiURL)
	if err != nil {
		return nil, fmt.Errorf("获取包 %s 信息失败: %w", packageName, err)
	}

	// 解析JSON响应
	var pkg models.Package
	if err := json.Unmarshal(responseBody, &pkg); err != nil {
		return nil, fmt.Errorf("解析包 %s 信息失败: %w", packageName, err)
	}

	return &pkg, nil
}

// GetPackageVersion 获取指定包的特定版本信息
//
// 参数:
//   - ctx: 上下文，用于控制请求的生命周期
//   - packageName: 包名
//   - version: 版本号
//
// 返回值:
//   - *models.Package: 包含特定版本详细信息的结构体指针
//   - error: 如有错误则返回，否则为nil
func (c *Client) GetPackageVersion(ctx context.Context, packageName string, version string) (*models.Package, error) {
	// 构建API URL
	apiURL := fmt.Sprintf("%s/pypi/%s/%s/json", c.options.BaseURL, url.PathEscape(packageName), url.PathEscape(version))

	// 发送请求
	responseBody, err := c.sendRequest(ctx, apiURL)
	if err != nil {
		return nil, fmt.Errorf("获取包 %s 版本 %s 信息失败: %w", packageName, version, err)
	}

	// 解析JSON响应
	var pkg models.Package
	if err := json.Unmarshal(responseBody, &pkg); err != nil {
		return nil, fmt.Errorf("解析包 %s 版本 %s 信息失败: %w", packageName, version, err)
	}

	return &pkg, nil
}

// GetPackageReleases 获取指定包的所有发布版本
//
// 参数:
//   - ctx: 上下文，用于控制请求的生命周期
//   - packageName: 包名
//
// 返回值:
//   - []string: 包含版本号的字符串切片
//   - error: 如有错误则返回，否则为nil
func (c *Client) GetPackageReleases(ctx context.Context, packageName string) ([]string, error) {
	// 通过获取包信息来获取所有版本
	pkg, err := c.GetPackageInfo(ctx, packageName)
	if err != nil {
		return nil, fmt.Errorf("获取包 %s 版本列表失败: %w", packageName, err)
	}

	// 提取版本号
	versions := make([]string, 0, len(pkg.Releases))
	for version := range pkg.Releases {
		versions = append(versions, version)
	}

	return versions, nil
}

// CheckPackageVulnerabilities 检查指定包和版本是否存在已知漏洞
//
// 参数:
//   - ctx: 上下文，用于控制请求的生命周期
//   - packageName: 包名
//   - version: 版本号
//
// 返回值:
//   - []models.Vulnerability: 包含漏洞信息的切片
//   - error: 如有错误则返回，否则为nil
func (c *Client) CheckPackageVulnerabilities(ctx context.Context, packageName string, version string) ([]models.Vulnerability, error) {
	// 获取特定版本的包信息，其中包含漏洞信息
	pkg, err := c.GetPackageVersion(ctx, packageName, version)
	if err != nil {
		return nil, fmt.Errorf("检查包 %s 版本 %s 漏洞失败: %w", packageName, version, err)
	}

	return pkg.Vulnerabilities, nil
}

// GetAllPackages 获取PyPI仓库中所有包的列表
//
// 该方法通过请求PyPI的Simple API获取所有可用包的索引列表
//
// 参数:
//   - ctx: 上下文，用于控制请求的生命周期
//
// 返回值:
//   - []string: 包含所有包名的切片
//   - error: 如有错误则返回，否则为nil
func (c *Client) GetAllPackages(ctx context.Context) ([]string, error) {
	// 构建Simple API URL
	simpleURL := fmt.Sprintf("%s/simple/", c.options.BaseURL)

	// 发送请求
	responseBody, err := c.sendRequest(ctx, simpleURL)
	if err != nil {
		return nil, fmt.Errorf("获取包索引失败: %w", err)
	}

	// 解析HTML响应
	return c.parsePackageIndex(string(responseBody))
}

// GetPackageList 获取PyPI仓库中所有包的列表（以map形式返回）
//
// 该方法是GetAllPackages的变种，返回map格式便于查询和遍历
//
// 参数:
//   - ctx: 上下文，用于控制请求的生命周期
//
// 返回值:
//   - map[string]struct{}: 包含所有包名的map，值为空结构体
//   - error: 如有错误则返回，否则为nil
func (c *Client) GetPackageList(ctx context.Context) (map[string]struct{}, error) {
	// 获取所有包的切片
	packages, err := c.GetAllPackages(ctx)
	if err != nil {
		return nil, err
	}

	// 转换为map
	packageMap := make(map[string]struct{}, len(packages))
	for _, pkg := range packages {
		packageMap[pkg] = struct{}{}
	}

	return packageMap, nil
}

// SearchPackages 根据关键词搜索包
//
// 该方法通过关键词搜索PyPI仓库中的包
//
// 参数:
//   - ctx: 上下文，用于控制请求的生命周期
//   - keyword: 搜索关键词
//   - limit: 最大返回结果数，如果为0则使用默认值(100)
//
// 返回值:
//   - []string: 匹配的包名列表
//   - error: 如有错误则返回，否则为nil
func (c *Client) SearchPackages(ctx context.Context, keyword string, limit int) ([]string, error) {
	// 先获取所有包
	packages, err := c.GetAllPackages(ctx)
	if err != nil {
		return nil, err
	}

	// 设置默认限制
	if limit <= 0 {
		limit = 100
	}

	// 搜索匹配的包
	var results []string
	keyword = strings.ToLower(keyword)

	for _, pkg := range packages {
		if strings.Contains(strings.ToLower(pkg), keyword) {
			results = append(results, pkg)

			// 达到限制时停止
			if len(results) >= limit {
				break
			}
		}
	}

	return results, nil
}

// parsePackageIndex 解析包索引页面的HTML内容
// 提取所有包名
//
// 参数:
//   - indexPageHTML: Simple API返回的HTML内容
//
// 返回值:
//   - []string: 提取的包名列表
//   - error: 解析错误，若无错误则为nil
func (c *Client) parsePackageIndex(indexPageHTML string) ([]string, error) {
	document, err := goquery.NewDocumentFromReader(strings.NewReader(indexPageHTML))
	if err != nil {
		return nil, fmt.Errorf("解析包索引页面失败: %w", err)
	}

	packageIndexes := make([]string, 0)
	document.Find("body a").Each(func(i int, selection *goquery.Selection) {
		packageName := strings.TrimSpace(selection.Text())
		if packageName == "" {
			return
		}
		packageIndexes = append(packageIndexes, packageName)
	})
	return packageIndexes, nil
}

// sendRequest 发送HTTP请求并返回响应体
//
// 内部使用函数，用于发送HTTP请求，处理重试和错误
//
// 参数:
//   - ctx: 上下文，用于控制请求的生命周期
//   - url: 目标URL
//
// 返回值:
//   - []byte: 响应体内容的字节数组
//   - error: 如有错误则返回，否则为nil
func (c *Client) sendRequest(ctx context.Context, requestURL string) ([]byte, error) {
	// 创建请求
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置用户代理头部
	req.Header.Set("User-Agent", c.options.UserAgent)
	req.Header.Set("Accept", "application/json")

	// 重试逻辑
	var resp *http.Response
	var lastErr error

	for attempt := 0; attempt < c.options.MaxRetries; attempt++ {
		// 如果不是第一次尝试，等待一段时间后重试
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(c.options.RetryDelay):
				// 继续重试
			}
		}

		// 发送请求
		resp, err = c.client.Do(req)
		if err == nil && resp.StatusCode < 500 {
			break // 成功或客户端错误（不重试）
		}

		// 记录最后一个错误
		if err != nil {
			lastErr = err
		} else {
			lastErr = fmt.Errorf("服务器错误: HTTP %d", resp.StatusCode)
			resp.Body.Close() // 关闭响应体以避免泄漏
		}
	}

	if resp == nil {
		if lastErr != nil {
			return nil, fmt.Errorf("请求失败，已重试 %d 次: %w", c.options.MaxRetries, lastErr)
		}
		return nil, fmt.Errorf("请求失败，已重试 %d 次", c.options.MaxRetries)
	}

	// 确保响应体最终会被关闭
	defer resp.Body.Close()

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP请求失败: %d %s", resp.StatusCode, resp.Status)
	}

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}

	return body, nil
}
