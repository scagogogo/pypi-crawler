package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReleaseUrls_Len(t *testing.T) {
	// 测试空集合
	emptyUrls := ReleaseUrls{}
	assert.Equal(t, 0, emptyUrls.Len())

	// 测试有一个元素的集合
	singleItemUrls := ReleaseUrls{
		&ReleaseURL{URL: "https://example.com/package1.whl"},
	}
	assert.Equal(t, 1, singleItemUrls.Len())

	// 测试有多个元素的集合
	multiItemUrls := ReleaseUrls{
		&ReleaseURL{URL: "https://example.com/package1.whl"},
		&ReleaseURL{URL: "https://example.com/package2.whl"},
		&ReleaseURL{URL: "https://example.com/package3.whl"},
	}
	assert.Equal(t, 3, multiItemUrls.Len())
}
