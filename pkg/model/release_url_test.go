package model

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 环境变量，用于控制是否运行需要网络连接的测试
// 如果未设置或设置为true，则运行连接测试
// 如果设置为其他值，则跳过连接测试
const ENV_RUN_NETWORK_TESTS = "RUN_NETWORK_TESTS"

// 检查是否应该运行网络连接测试
func shouldRunNetworkTests() bool {
	val := os.Getenv(ENV_RUN_NETWORK_TESTS)
	return val == "" || val == "true" || val == "1"
}

func TestReleaseURL_Download(t *testing.T) {
	// 跳过需要网络连接的测试
	if !shouldRunNetworkTests() {
		t.Skip("Skipping network-dependent test")
		return
	}

	// 测试正常的下载URL
	releaseURL := &ReleaseURL{
		URL: "https://pypi.org/simple",
	}
	data, err := releaseURL.Download(context.Background())
	assert.Nil(t, err)
	assert.NotNil(t, data)
	assert.True(t, len(data) > 0)

	// 测试空URL
	emptyURL := &ReleaseURL{
		URL: "",
	}
	data, err = emptyURL.Download(context.Background())
	assert.NotNil(t, err)
	assert.Nil(t, data)
	assert.Contains(t, err.Error(), "do not has release file url")

	// 测试无效URL
	invalidURL := &ReleaseURL{
		URL: "https://invalid-url.example",
	}
	data, err = invalidURL.Download(context.Background())
	assert.NotNil(t, err)
	assert.Nil(t, data)
}
