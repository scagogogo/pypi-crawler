package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReleaseFile_GetUploadTimeISO(t *testing.T) {
	t.Run("使用ISO8601时间", func(t *testing.T) {
		file := &ReleaseFile{
			UploadTimeISO8601: "2023-05-15T10:30:00Z",
			UploadTime:        "2023-05-15T10:30:00", // 旧格式，应该被忽略
		}

		uploadTime, err := file.GetUploadTimeISO()
		require.NoError(t, err)
		assert.Equal(t, 2023, uploadTime.Year())
		assert.Equal(t, time.May, uploadTime.Month())
		assert.Equal(t, 15, uploadTime.Day())
	})

	t.Run("使用旧格式时间", func(t *testing.T) {
		file := &ReleaseFile{
			UploadTimeISO8601: "", // 为空，使用旧格式
			UploadTime:        "2023-05-15T10:30:00",
		}

		uploadTime, err := file.GetUploadTimeISO()
		require.NoError(t, err)
		assert.Equal(t, 2023, uploadTime.Year())
		assert.Equal(t, time.May, uploadTime.Month())
		assert.Equal(t, 15, uploadTime.Day())
	})

	t.Run("两种格式都为空", func(t *testing.T) {
		file := &ReleaseFile{
			UploadTimeISO8601: "",
			UploadTime:        "",
		}

		uploadTime, err := file.GetUploadTimeISO()
		require.NoError(t, err)             // 不应该返回错误
		assert.True(t, uploadTime.IsZero()) // 应该返回零值
	})

	t.Run("非法格式", func(t *testing.T) {
		file := &ReleaseFile{
			UploadTimeISO8601: "not-a-date",
		}

		_, err := file.GetUploadTimeISO()
		assert.Error(t, err) // 应该返回解析错误
	})
}

func TestReleaseFile_IsYanked(t *testing.T) {
	t.Run("已撤回", func(t *testing.T) {
		file := &ReleaseFile{
			Yanked: true,
		}

		assert.True(t, file.IsYanked())
	})

	t.Run("未撤回", func(t *testing.T) {
		file := &ReleaseFile{
			Yanked: false,
		}

		assert.False(t, file.IsYanked())
	})
}

func TestReleaseFile_IsWheel(t *testing.T) {
	t.Run("是wheel格式", func(t *testing.T) {
		file := &ReleaseFile{
			PackageType: "bdist_wheel",
		}

		assert.True(t, file.IsWheel())
	})

	t.Run("不是wheel格式", func(t *testing.T) {
		file := &ReleaseFile{
			PackageType: "sdist",
		}

		assert.False(t, file.IsWheel())
	})
}

func TestReleaseFile_IsSourceDist(t *testing.T) {
	t.Run("是源码包", func(t *testing.T) {
		file := &ReleaseFile{
			PackageType: "sdist",
		}

		assert.True(t, file.IsSourceDist())
	})

	t.Run("不是源码包", func(t *testing.T) {
		file := &ReleaseFile{
			PackageType: "bdist_wheel",
		}

		assert.False(t, file.IsSourceDist())
	})
}
