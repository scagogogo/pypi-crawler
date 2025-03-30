package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewOptions(t *testing.T) {
	// 测试默认选项
	options := NewOptions()
	assert.NotNil(t, options)
	assert.Equal(t, DefaultServerURL, options.ServerURL)
	assert.Equal(t, "", options.Proxy)
}

func TestOptions_SetServerURL(t *testing.T) {
	// 测试设置服务器URL
	options := NewOptions()
	testURL := "https://test-pypi.org"

	// 测试方法链式调用并返回自身
	returnedOptions := options.SetServerURL(testURL)
	assert.Equal(t, options, returnedOptions)

	// 测试值已正确设置
	assert.Equal(t, testURL, options.ServerURL)

	// 测试设置空URL
	options.SetServerURL("")
	assert.Equal(t, "", options.ServerURL)
}

func TestOptions_SetProxy(t *testing.T) {
	// 测试设置代理
	options := NewOptions()
	testProxy := "http://127.0.0.1:8080"

	// 测试方法链式调用并返回自身
	returnedOptions := options.SetProxy(testProxy)
	assert.Equal(t, options, returnedOptions)

	// 测试值已正确设置
	assert.Equal(t, testProxy, options.Proxy)

	// 测试设置空代理
	options.SetProxy("")
	assert.Equal(t, "", options.Proxy)
}
