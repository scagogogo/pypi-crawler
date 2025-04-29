package mirrors

import (
	"testing"

	"github.com/scagogogo/pypi-crawler/pkg/pypi/api"
	"github.com/scagogogo/pypi-crawler/pkg/pypi/client"
	"github.com/stretchr/testify/assert"
)

func TestNewOfficialClient(t *testing.T) {
	t.Run("默认选项", func(t *testing.T) {
		c := NewOfficialClient()
		assert.NotNil(t, c)

		// 检查是否返回了不为nil的接口实现
		assert.NotNil(t, c)
		assert.Implements(t, (*api.PyPIClient)(nil), c)
	})

	t.Run("自定义选项", func(t *testing.T) {
		options := client.NewOptions().WithUserAgent("TestAgent")
		c := NewOfficialClient(options)
		assert.NotNil(t, c)
		assert.Implements(t, (*api.PyPIClient)(nil), c)
	})
}

func TestNewTsinghuaClient(t *testing.T) {
	c := NewTsinghuaClient()
	assert.NotNil(t, c)
	assert.Implements(t, (*api.PyPIClient)(nil), c)
}

func TestNewDoubanClient(t *testing.T) {
	c := NewDoubanClient()
	assert.NotNil(t, c)
	assert.Implements(t, (*api.PyPIClient)(nil), c)
}

func TestNewAliyunClient(t *testing.T) {
	c := NewAliyunClient()
	assert.NotNil(t, c)
	assert.Implements(t, (*api.PyPIClient)(nil), c)
}

func TestNewTencentClient(t *testing.T) {
	c := NewTencentClient()
	assert.NotNil(t, c)
	assert.Implements(t, (*api.PyPIClient)(nil), c)
}

func TestNewUstcClient(t *testing.T) {
	c := NewUstcClient()
	assert.NotNil(t, c)
	assert.Implements(t, (*api.PyPIClient)(nil), c)
}

func TestNewNeteaseClient(t *testing.T) {
	c := NewNeteaseClient()
	assert.NotNil(t, c)
	assert.Implements(t, (*api.PyPIClient)(nil), c)
}

// 测试所有镜像源的URL是否正确
func TestAllMirrorURLs(t *testing.T) {
	assert.Equal(t, "https://pypi.org", OfficialURL)
	assert.Equal(t, "https://pypi.tuna.tsinghua.edu.cn", TsinghuaURL)
	assert.Equal(t, "https://pypi.doubanio.com", DoubanURL)
	assert.Equal(t, "https://mirrors.aliyun.com/pypi", AliyunURL)
	assert.Equal(t, "https://mirrors.cloud.tencent.com/pypi", TencentURL)
	assert.Equal(t, "https://pypi.mirrors.ustc.edu.cn", UstcURL)
	assert.Equal(t, "https://mirrors.163.com/pypi", NeteaseURL)
}
