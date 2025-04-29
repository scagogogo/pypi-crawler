package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPackageInfo_GetAllDependencies(t *testing.T) {
	t.Run("获取已定义的依赖", func(t *testing.T) {
		info := &PackageInfo{
			RequiresDist: []string{"certifi>=2017.4.17", "urllib3>=1.21.1"},
		}

		deps := info.GetAllDependencies()
		assert.Len(t, deps, 2)
		assert.Equal(t, []string{"certifi>=2017.4.17", "urllib3>=1.21.1"}, deps)
	})

	t.Run("获取空依赖", func(t *testing.T) {
		info := &PackageInfo{}

		deps := info.GetAllDependencies()
		assert.Empty(t, deps)
		assert.NotNil(t, deps)
	})

	t.Run("nil依赖字段", func(t *testing.T) {
		info := &PackageInfo{
			RequiresDist: nil,
		}

		deps := info.GetAllDependencies()
		assert.Empty(t, deps)
		assert.NotNil(t, deps)
	})
}

func TestPackageInfo_HasPythonRequirement(t *testing.T) {
	t.Run("有Python版本要求", func(t *testing.T) {
		info := &PackageInfo{
			RequiresPython: ">=3.6",
		}

		assert.True(t, info.HasPythonRequirement())
	})

	t.Run("无Python版本要求", func(t *testing.T) {
		info := &PackageInfo{
			RequiresPython: "",
		}

		assert.False(t, info.HasPythonRequirement())
	})
}

func TestPackageInfo_IsYanked(t *testing.T) {
	t.Run("包已撤回", func(t *testing.T) {
		info := &PackageInfo{
			Yanked:       true,
			YankedReason: "安全问题",
		}

		assert.True(t, info.IsYanked())
	})

	t.Run("包未撤回", func(t *testing.T) {
		info := &PackageInfo{
			Yanked: false,
		}

		assert.False(t, info.IsYanked())
	})
}

func TestPackageInfo_GetProjectURLs(t *testing.T) {
	t.Run("获取已定义的URL", func(t *testing.T) {
		info := &PackageInfo{
			ProjectURLs: map[string]string{
				"Documentation": "https://docs.example.com",
				"Source":        "https://github.com/example/pkg",
			},
		}

		urls := info.GetProjectURLs()
		assert.Len(t, urls, 2)
		assert.Equal(t, "https://docs.example.com", urls["Documentation"])
		assert.Equal(t, "https://github.com/example/pkg", urls["Source"])
	})

	t.Run("获取空URL", func(t *testing.T) {
		info := &PackageInfo{}

		urls := info.GetProjectURLs()
		assert.Empty(t, urls)
		assert.NotNil(t, urls)
	})

	t.Run("nil URL字段", func(t *testing.T) {
		info := &PackageInfo{
			ProjectURLs: nil,
		}

		urls := info.GetProjectURLs()
		assert.Empty(t, urls)
		assert.NotNil(t, urls)
	})
}
