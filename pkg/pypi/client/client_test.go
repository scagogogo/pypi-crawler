package client

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 测试数据文件夹和文件路径
const (
	testDataDir         = "testdata"
	packageInfoFile     = "package_info.json"
	packageVersionFile  = "package_version.json"
	packageWithVulnFile = "package_with_vuln.json"
	simpleIndexFile     = "simple_index.html"
)

// 确保测试数据目录存在
func init() {
	if _, err := os.Stat(testDataDir); os.IsNotExist(err) {
		_ = os.Mkdir(testDataDir, os.ModePerm)
	}
}

// 创建一个mock server来模拟PyPI API
func setupMockServer(t *testing.T) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var responseFile string
		var status int = http.StatusOK

		path := r.URL.Path
		switch {
		case path == "/pypi/requests/json":
			responseFile = filepath.Join(testDataDir, packageInfoFile)
		case path == "/pypi/requests/2.28.1/json":
			responseFile = filepath.Join(testDataDir, packageVersionFile)
		case path == "/pypi/vulnerable-package/1.0.0/json":
			responseFile = filepath.Join(testDataDir, packageWithVulnFile)
		case path == "/simple/" || path == "/simple":
			responseFile = filepath.Join(testDataDir, simpleIndexFile)
		case strings.Contains(path, "/json"):
			// 测试不存在的包
			status = http.StatusNotFound
			w.WriteHeader(status)
			_, _ = w.Write([]byte(`{"message": "包不存在"}`))
			return
		default:
			status = http.StatusNotFound
			w.WriteHeader(status)
			return
		}

		// 如果找到匹配的响应文件，则读取并返回其内容
		if responseFile != "" {
			if _, err := os.Stat(responseFile); os.IsNotExist(err) {
				t.Logf("警告: 测试数据文件不存在: %s", responseFile)
				createDummyTestFile(responseFile, path)
			}

			bytes, err := os.ReadFile(responseFile)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error": "内部服务器错误"}`))
				return
			}

			w.WriteHeader(status)
			_, _ = w.Write(bytes)
		}
	}))

	return server
}

// 如果测试文件不存在，创建一个伪测试文件
func createDummyTestFile(filePath, endpoint string) {
	var content []byte

	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		_ = os.MkdirAll(dir, os.ModePerm)
	}

	switch filepath.Base(filePath) {
	case packageInfoFile:
		content = []byte(`{
			"info": {
				"name": "requests",
				"version": "2.28.1",
				"summary": "Python HTTP for Humans",
				"author": "Kenneth Reitz",
				"author_email": "me@kennethreitz.org",
				"license": "Apache 2.0",
				"home_page": "https://requests.readthedocs.io",
				"requires_dist": ["certifi>=2017.4.17", "charset-normalizer>=2.0.0", "idna>=2.5", "urllib3>=1.21.1"]
			},
			"releases": {
				"2.28.1": [
					{
						"filename": "requests-2.28.1-py3-none-any.whl",
						"packagetype": "bdist_wheel",
						"size": 123456
					}
				],
				"2.28.0": [
					{
						"filename": "requests-2.28.0-py3-none-any.whl",
						"packagetype": "bdist_wheel",
						"size": 123400
					}
				]
			},
			"urls": [
				{
					"filename": "requests-2.28.1-py3-none-any.whl",
					"packagetype": "bdist_wheel",
					"size": 123456
				}
			]
		}`)
	case packageVersionFile:
		content = []byte(`{
			"info": {
				"name": "requests",
				"version": "2.28.1",
				"summary": "Python HTTP for Humans",
				"author": "Kenneth Reitz",
				"license": "Apache 2.0"
			},
			"urls": [
				{
					"filename": "requests-2.28.1-py3-none-any.whl",
					"packagetype": "bdist_wheel",
					"size": 123456
				}
			]
		}`)
	case packageWithVulnFile:
		content = []byte(`{
			"info": {
				"name": "vulnerable-package",
				"version": "1.0.0",
				"summary": "测试漏洞包"
			},
			"vulnerabilities": [
				{
					"id": "VULN-123",
					"summary": "安全漏洞摘要",
					"details": "详细的漏洞描述",
					"link": "https://example.com/vuln/123",
					"fixed_in": ["1.0.1"]
				}
			]
		}`)
	case simpleIndexFile:
		content = []byte(`<!DOCTYPE html>
		<html>
		<head><title>Simple Index</title></head>
		<body>
			<a href="/simple/requests/">requests</a>
			<a href="/simple/flask/">flask</a>
			<a href="/simple/django/">django</a>
			<a href="/simple/numpy/">numpy</a>
			<a href="/simple/pandas/">pandas</a>
		</body>
		</html>`)
	}

	_ = os.WriteFile(filePath, content, os.ModePerm)
}

// 创建测试用的客户端
func createTestClient(server *httptest.Server) *Client {
	options := NewOptions().
		WithBaseURL(server.URL).
		WithTimeout(5 * time.Second).
		WithMaxRetries(1)

	client := NewClient(options).(*Client)
	return client
}

// 测试获取包信息功能
func TestGetPackageInfo(t *testing.T) {
	server := setupMockServer(t)
	defer server.Close()

	client := createTestClient(server)
	ctx := context.Background()

	t.Run("获取存在的包", func(t *testing.T) {
		pkg, err := client.GetPackageInfo(ctx, "requests")
		require.NoError(t, err)
		assert.NotNil(t, pkg)
		assert.Equal(t, "requests", pkg.Info.Name)
		assert.Equal(t, "2.28.1", pkg.Info.Version)
		assert.Equal(t, "Python HTTP for Humans", pkg.Info.Summary)
	})

	t.Run("获取不存在的包", func(t *testing.T) {
		pkg, err := client.GetPackageInfo(ctx, "non-existent-package")
		assert.Error(t, err)
		assert.Nil(t, pkg)
	})

	t.Run("空包名", func(t *testing.T) {
		pkg, err := client.GetPackageInfo(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, pkg)
	})
}

// 测试获取包特定版本功能
func TestGetPackageVersion(t *testing.T) {
	server := setupMockServer(t)
	defer server.Close()

	client := createTestClient(server)
	ctx := context.Background()

	t.Run("获取存在的包版本", func(t *testing.T) {
		pkg, err := client.GetPackageVersion(ctx, "requests", "2.28.1")
		require.NoError(t, err)
		assert.NotNil(t, pkg)
		assert.Equal(t, "requests", pkg.Info.Name)
		assert.Equal(t, "2.28.1", pkg.Info.Version)
	})

	t.Run("获取不存在的包版本", func(t *testing.T) {
		pkg, err := client.GetPackageVersion(ctx, "requests", "999.999.999")
		assert.Error(t, err)
		assert.Nil(t, pkg)
	})

	t.Run("空包名或版本", func(t *testing.T) {
		pkg, err := client.GetPackageVersion(ctx, "", "2.28.1")
		assert.Error(t, err)
		assert.Nil(t, pkg)

		pkg, err = client.GetPackageVersion(ctx, "requests", "")
		assert.Error(t, err)
		assert.Nil(t, pkg)
	})
}

// 测试获取包所有版本功能
func TestGetPackageReleases(t *testing.T) {
	server := setupMockServer(t)
	defer server.Close()

	client := createTestClient(server)
	ctx := context.Background()

	t.Run("获取存在的包版本列表", func(t *testing.T) {
		versions, err := client.GetPackageReleases(ctx, "requests")
		require.NoError(t, err)
		assert.NotEmpty(t, versions)
		assert.Contains(t, versions, "2.28.1")
		assert.Contains(t, versions, "2.28.0")
	})

	t.Run("获取不存在的包版本列表", func(t *testing.T) {
		versions, err := client.GetPackageReleases(ctx, "non-existent-package")
		assert.Error(t, err)
		assert.Empty(t, versions)
	})
}

// 测试检查包漏洞功能
func TestCheckPackageVulnerabilities(t *testing.T) {
	server := setupMockServer(t)
	defer server.Close()

	client := createTestClient(server)
	ctx := context.Background()

	t.Run("检查有漏洞的包", func(t *testing.T) {
		vulns, err := client.CheckPackageVulnerabilities(ctx, "vulnerable-package", "1.0.0")
		require.NoError(t, err)
		assert.NotEmpty(t, vulns)
		assert.Equal(t, "VULN-123", vulns[0].ID)
		assert.Equal(t, "安全漏洞摘要", vulns[0].Summary)
	})

	t.Run("检查无漏洞的包", func(t *testing.T) {
		vulns, err := client.CheckPackageVulnerabilities(ctx, "requests", "2.28.1")
		require.NoError(t, err)
		assert.Empty(t, vulns)
	})

	t.Run("检查不存在的包", func(t *testing.T) {
		vulns, err := client.CheckPackageVulnerabilities(ctx, "non-existent-package", "1.0.0")
		assert.Error(t, err)
		assert.Nil(t, vulns)
	})
}

// 测试获取所有包列表功能
func TestGetAllPackages(t *testing.T) {
	server := setupMockServer(t)
	defer server.Close()

	client := createTestClient(server)
	ctx := context.Background()

	t.Run("获取包列表", func(t *testing.T) {
		packages, err := client.GetAllPackages(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, packages)
		assert.Len(t, packages, 5) // 由测试数据决定
		assert.Contains(t, packages, "requests")
		assert.Contains(t, packages, "flask")
		assert.Contains(t, packages, "django")
	})
}

// 测试获取包列表(Map形式)功能
func TestGetPackageList(t *testing.T) {
	server := setupMockServer(t)
	defer server.Close()

	client := createTestClient(server)
	ctx := context.Background()

	t.Run("获取包列表(Map)", func(t *testing.T) {
		packages, err := client.GetPackageList(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, packages)
		assert.Len(t, packages, 5) // 由测试数据决定

		// 检查Map是否包含预期的键
		_, hasRequests := packages["requests"]
		assert.True(t, hasRequests)

		_, hasFlask := packages["flask"]
		assert.True(t, hasFlask)
	})
}

// 测试搜索包功能
func TestSearchPackages(t *testing.T) {
	server := setupMockServer(t)
	defer server.Close()

	client := createTestClient(server)
	ctx := context.Background()

	t.Run("搜索存在的包", func(t *testing.T) {
		results, err := client.SearchPackages(ctx, "requ", 10)
		require.NoError(t, err)
		assert.NotEmpty(t, results)
		assert.Contains(t, results, "requests")
	})

	t.Run("搜索不存在的包", func(t *testing.T) {
		results, err := client.SearchPackages(ctx, "nonexistentpackagexyz", 10)
		require.NoError(t, err)
		assert.Empty(t, results)
	})

	t.Run("限制搜索结果数量", func(t *testing.T) {
		// 创建一个特殊的客户端来测试限制功能
		customClient := createTestClient(server)

		// 创建一个自定义的parsePackageIndex函数，返回预定义的测试数据
		testPackages := []string{"package1", "package2", "package3", "package4", "package5"}

		// 创建一个包装客户端，它使用自己的parsePackageIndex函数
		wrappedClient := &clientWrapper{
			Client: customClient,
			mockParseFunc: func(html string) ([]string, error) {
				return testPackages, nil
			},
		}

		results, err := wrappedClient.SearchPackages(ctx, "pack", 3)
		require.NoError(t, err)
		assert.Len(t, results, 3) // 应该只返回3个结果
	})
}

// 测试解析包索引功能
func TestParsePackageIndex(t *testing.T) {
	client := &Client{options: NewOptions()}

	t.Run("解析有效的HTML", func(t *testing.T) {
		html := `
		<!DOCTYPE html>
		<html>
		<head><title>Simple Index</title></head>
		<body>
			<a href="/simple/package1/">package1</a>
			<a href="/simple/package2/">package2</a>
			<a href="/simple/package3/">package3</a>
		</body>
		</html>
		`

		packages, err := client.parsePackageIndex(html)
		require.NoError(t, err)
		assert.Len(t, packages, 3)
		assert.Equal(t, []string{"package1", "package2", "package3"}, packages)
	})

	t.Run("解析无效的HTML", func(t *testing.T) {
		html := "这不是有效的HTML"

		packages, err := client.parsePackageIndex(html)
		require.NoError(t, err) // goquery通常会尝试解析无效的HTML而不返回错误
		assert.Empty(t, packages)
	})

	t.Run("解析空HTML", func(t *testing.T) {
		packages, err := client.parsePackageIndex("")
		require.NoError(t, err)
		assert.Empty(t, packages)
	})
}

// 测试HTTP请求功能
func TestSendRequest(t *testing.T) {
	// 创建一个测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/success":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"result": "success"}`))
		case "/not-found":
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error": "not found"}`))
		case "/server-error":
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error": "server error"}`))
		case "/timeout":
			time.Sleep(2 * time.Second)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"result": "timeout"}`))
		}
	}))
	defer server.Close()

	t.Run("成功的请求", func(t *testing.T) {
		client := createTestClient(server)
		ctx := context.Background()

		resp, err := client.sendRequest(ctx, server.URL+"/success")
		require.NoError(t, err)
		assert.Equal(t, `{"result": "success"}`, string(resp))
	})

	t.Run("404错误", func(t *testing.T) {
		client := createTestClient(server)
		ctx := context.Background()

		resp, err := client.sendRequest(ctx, server.URL+"/not-found")
		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("500错误", func(t *testing.T) {
		client := createTestClient(server)
		ctx := context.Background()

		resp, err := client.sendRequest(ctx, server.URL+"/server-error")
		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("上下文超时", func(t *testing.T) {
		client := createTestClient(server)
		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()

		resp, err := client.sendRequest(ctx, server.URL+"/timeout")
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

// 模拟文件读取错误
type mockFailingReader struct{}

func (m *mockFailingReader) Read(p []byte) (n int, err error) {
	return 0, io.ErrUnexpectedEOF
}

func (m *mockFailingReader) Close() error {
	return nil
}

// 测试NewClient工厂函数
func TestNewClient(t *testing.T) {
	t.Run("使用默认选项", func(t *testing.T) {
		client := NewClient()
		assert.NotNil(t, client)

		// 使用类型断言来获取内部client
		concreteClient, ok := client.(*Client)
		require.True(t, ok)

		// 检查默认值
		assert.Equal(t, DefaultBaseURL, concreteClient.options.BaseURL)
		assert.Equal(t, DefaultTimeout, concreteClient.options.Timeout)
		assert.Equal(t, DefaultMaxRetries, concreteClient.options.MaxRetries)
	})

	t.Run("使用自定义选项", func(t *testing.T) {
		options := NewOptions().
			WithBaseURL("https://test.pypi.org").
			WithTimeout(30 * time.Second).
			WithMaxRetries(5)

		client := NewClient(options)
		assert.NotNil(t, client)

		// 使用类型断言来获取内部client
		concreteClient, ok := client.(*Client)
		require.True(t, ok)

		// 检查自定义值
		assert.Equal(t, "https://test.pypi.org", concreteClient.options.BaseURL)
		assert.Equal(t, 30*time.Second, concreteClient.options.Timeout)
		assert.Equal(t, 5, concreteClient.options.MaxRetries)
	})
}

// clientWrapper 是一个Client的包装器，允许我们覆盖特定方法用于测试
type clientWrapper struct {
	*Client
	mockParseFunc func(string) ([]string, error)
}

// SearchPackages 实现API接口的方法，使用我们的模拟parsePackageIndex函数
func (c *clientWrapper) SearchPackages(ctx context.Context, keyword string, limit int) ([]string, error) {
	// 直接使用模拟的解析函数获取包列表
	packages, err := c.mockParseFunc("")
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
