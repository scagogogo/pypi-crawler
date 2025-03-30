package repository

import (
	"context"
	"fmt"
	"os"
	"testing"

	project_root_directory "github.com/golang-infrastructure/go-project-root-directory"
	"github.com/stretchr/testify/assert"
)

func TestNewRepository(t *testing.T) {
	// 测试默认仓库创建
	r := NewRepository()
	assert.NotNil(t, r)
	assert.NotNil(t, r.options)
	assert.Equal(t, DefaultServerURL, r.options.ServerURL)
	assert.Equal(t, "", r.options.Proxy)

	// 测试带自定义选项的仓库创建
	customOptions := NewOptions().SetServerURL("https://custom-pypi.org").SetProxy("http://proxy")
	r = NewRepository(customOptions)
	assert.NotNil(t, r)
	assert.Equal(t, "https://custom-pypi.org", r.options.ServerURL)
	assert.Equal(t, "http://proxy", r.options.Proxy)
}

func TestPypiRepository_GetPackage(t *testing.T) {
	// 跳过需要网络连接的测试
	if !shouldRunNetworkTests() {
		t.Skip("Skipping network-dependent test")
		return
	}

	// 正常情况
	pkg, err := NewRepository().GetPackage(context.Background(), "requests")
	assert.Nil(t, err)
	assert.NotNil(t, pkg)
	assert.Equal(t, "requests", pkg.Information.Name)

	// 测试不存在的包
	pkg, err = NewRepository().GetPackage(context.Background(), "this-package-does-not-exist-12345678909876543210")
	assert.NotNil(t, err)
	assert.Nil(t, pkg)

	// 测试空包名
	pkg, err = NewRepository().GetPackage(context.Background(), "")
	assert.NotNil(t, err)
	assert.Nil(t, pkg)

	// 测试错误的仓库URL
	r := NewRepository(NewOptions().SetServerURL("https://invalid-url.example"))
	pkg, err = r.GetPackage(context.Background(), "requests")
	assert.NotNil(t, err)
	assert.Nil(t, pkg)
}

func TestPypiRepository_DownloadIndex(t *testing.T) {
	// 跳过需要网络连接的测试
	if !shouldRunNetworkTests() {
		t.Skip("Skipping network-dependent test")
		return
	}

	// 正常情况
	indexPageBytes, err := NewRepository().DownloadIndex(context.Background())
	assert.Nil(t, err)
	assert.True(t, len(indexPageBytes) > 0)

	// 测试错误的仓库URL
	r := NewRepository(NewOptions().SetServerURL("https://invalid-url.example"))
	indexPageBytes, err = r.DownloadIndex(context.Background())
	assert.NotNil(t, err)
	assert.Nil(t, indexPageBytes)
}

func TestPypiRepository_ParseIndexPage(t *testing.T) {
	// 直接使用内联的HTML字符串进行测试，不依赖外部文件

	// 测试包含包列表的HTML
	sampleHTML := `<!DOCTYPE html>
<html>
  <head>
    <title>Simple Index</title>
  </head>
  <body>
    <a href="/simple/package1/">package1</a>
    <a href="/simple/package2/">package2</a>
    <a href="/simple/requests/">requests</a>
    <a href="/simple/flask/">flask</a>
    <a href="/simple/django/">django</a>
  </body>
</html>`

	packageIndexes, err := NewRepository().ParseIndexPage(sampleHTML)
	assert.Nil(t, err)
	assert.NotNil(t, packageIndexes)
	assert.Equal(t, 5, len(packageIndexes))
	assert.Equal(t, "package1", packageIndexes[0])
	assert.Equal(t, "package2", packageIndexes[1])
	assert.Equal(t, "requests", packageIndexes[2])
	assert.Equal(t, "flask", packageIndexes[3])
	assert.Equal(t, "django", packageIndexes[4])

	// 可选：尝试加载外部sample.html文件，但如果不存在则不会导致测试失败
	indexFilepath, err := project_root_directory.GetRootFilePath("data/sample.html")
	if err == nil {
		indexPageBytes, err := os.ReadFile(indexFilepath)
		if err == nil {
			packageIndexes, err = NewRepository().ParseIndexPage(string(indexPageBytes))
			assert.Nil(t, err)
			assert.NotNil(t, packageIndexes)
			assert.True(t, len(packageIndexes) > 0)
		}
	}

	// 测试无包的HTML
	emptyHtml := `<html><body></body></html>`
	packageIndexes, err = NewRepository().ParseIndexPage(emptyHtml)
	assert.Nil(t, err)
	assert.NotNil(t, packageIndexes)
	assert.Equal(t, 0, len(packageIndexes))

	// 测试包含有包的HTML
	htmlWithPackages := `<html><body><a>package1</a><a>package2</a></body></html>`
	packageIndexes, err = NewRepository().ParseIndexPage(htmlWithPackages)
	assert.Nil(t, err)
	assert.NotNil(t, packageIndexes)
	assert.Equal(t, 2, len(packageIndexes))
	assert.Equal(t, "package1", packageIndexes[0])
	assert.Equal(t, "package2", packageIndexes[1])

	// 测试错误的HTML
	// 注意：goquery可能会尝试修复不完整的HTML，所以使用更加无效的内容
	invalidHtml := `<not-valid-html>`
	packageIndexes, err = NewRepository().ParseIndexPage(invalidHtml)
	// 如果goquery能修复HTML，那么err可能为nil，我们就检查是否返回了一个空列表
	if err == nil {
		assert.Equal(t, 0, len(packageIndexes))
	} else {
		assert.NotNil(t, err)
	}
}

func TestRepository_getBytes(t *testing.T) {
	// 跳过需要网络连接的测试
	if !shouldRunNetworkTests() {
		t.Skip("Skipping network-dependent test")
		return
	}

	// 这是一个内部方法，但可以通过反射或导出来测试
	// 这里我们只测试能否正常获取字节
	r := NewRepository()
	data, err := r.getBytes(context.Background(), fmt.Sprintf("%s/simple", DefaultServerURL))
	assert.Nil(t, err)
	assert.NotNil(t, data)
	assert.True(t, len(data) > 0)

	// 测试无效URL
	data, err = r.getBytes(context.Background(), "https://invalid-url.example")
	assert.NotNil(t, err)
	assert.Nil(t, data)
}

func RepositoryTest(t *testing.T, r *Repository) {
	// 不为空
	assert.NotNil(t, r)

	// 能够正常获取到索引
	index, err := r.DownloadIndex(context.Background())
	assert.Nil(t, err)
	if err != nil {
		t.Log(err.Error())
	}
	assert.True(t, len(index) > 0)

	// 能够获取到包的信息
	//packageName := index[rand.Intn(len(index))]
	packageName := "requests"
	packageInformation, err := r.GetPackage(context.Background(), packageName)
	assert.Nil(t, err)
	assert.NotNil(t, packageInformation)
	assert.Equal(t, packageName, packageInformation.Information.Name)
}
