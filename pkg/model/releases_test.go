package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReleases_UnmarshalJSON(t *testing.T) {
	// 准备测试数据
	jsonData := []byte(`{
		"1.0.0": [
			{
				"url": "https://example.com/package-1.0.0.whl",
				"filename": "package-1.0.0.whl",
				"packagetype": "bdist_wheel"
			}
		],
		"1.1.0": [
			{
				"url": "https://example.com/package-1.1.0.whl",
				"filename": "package-1.1.0.whl",
				"packagetype": "bdist_wheel"
			}
		]
	}`)

	// 测试正常反序列化
	releases := &Releases{}
	err := json.Unmarshal(jsonData, releases)
	assert.Nil(t, err)
	assert.NotNil(t, releases)
	assert.Equal(t, 2, len(releases.VersionOrders))
	assert.Equal(t, "1.0.0", releases.VersionOrders[0])
	assert.Equal(t, "1.1.0", releases.VersionOrders[1])
	assert.NotNil(t, releases.VersionMap)
	assert.Equal(t, 2, len(releases.VersionMap))
	assert.NotNil(t, releases.VersionMap["1.0.0"])
	assert.NotNil(t, releases.VersionMap["1.1.0"])
	assert.Equal(t, "https://example.com/package-1.0.0.whl", (*releases.VersionMap["1.0.0"])[0].URL)
	assert.Equal(t, "https://example.com/package-1.1.0.whl", (*releases.VersionMap["1.1.0"])[0].URL)

	// 测试空JSON
	emptyJSON := []byte(`{}`)
	emptyReleases := &Releases{}
	err = json.Unmarshal(emptyJSON, emptyReleases)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(emptyReleases.VersionOrders))
	assert.NotNil(t, emptyReleases.VersionMap)
	assert.Equal(t, 0, len(emptyReleases.VersionMap))

	// 测试格式错误的JSON
	invalidJSON := []byte(`{"1.0.0": "not an array"}`)
	invalidReleases := &Releases{}
	err = json.Unmarshal(invalidJSON, invalidReleases)
	assert.NotNil(t, err)
}

func TestReleases_MarshalJSON(t *testing.T) {
	// 创建一个有序的Releases对象
	releases := &Releases{
		VersionOrders: []string{"1.0.0", "1.1.0"},
		VersionMap: map[string]*ReleaseUrls{
			"1.0.0": &ReleaseUrls{
				&ReleaseURL{
					URL:         "https://example.com/package-1.0.0.whl",
					Filename:    "package-1.0.0.whl",
					Packagetype: "bdist_wheel",
				},
			},
			"1.1.0": &ReleaseUrls{
				&ReleaseURL{
					URL:         "https://example.com/package-1.1.0.whl",
					Filename:    "package-1.1.0.whl",
					Packagetype: "bdist_wheel",
				},
			},
		},
	}

	// 测试序列化
	data, err := json.Marshal(releases)
	assert.Nil(t, err)
	assert.NotNil(t, data)
	assert.True(t, len(data) > 0)

	// 反序列化验证顺序保留
	newReleases := &Releases{}
	err = json.Unmarshal(data, newReleases)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(newReleases.VersionOrders))
	assert.Equal(t, "1.0.0", newReleases.VersionOrders[0])
	assert.Equal(t, "1.1.0", newReleases.VersionOrders[1])

	// 测试空对象的序列化
	emptyReleases := &Releases{
		VersionOrders: []string{},
		VersionMap:    map[string]*ReleaseUrls{},
	}
	emptyData, err := json.Marshal(emptyReleases)
	assert.Nil(t, err)
	assert.NotNil(t, emptyData)
	assert.Equal(t, "{}", string(emptyData))
}
