package models

import "time"

// ReleaseFile 表示包的一个发布文件
// 包含了文件的详细信息，如URL、哈希值、大小等
type ReleaseFile struct {
	// Filename 文件名
	Filename string `json:"filename"`

	// URL 文件下载URL
	URL string `json:"url"`

	// PackageType 包类型，如sdist、bdist_wheel等
	PackageType string `json:"packagetype"`

	// PythonVersion Python版本要求
	PythonVersion string `json:"python_version"`

	// RequiresPython Python版本要求（字符串形式）
	RequiresPython string `json:"requires_python"`

	// Size 文件大小（字节）
	Size int64 `json:"size"`

	// UploadTime 上传时间（旧格式）
	UploadTime string `json:"upload_time"`

	// UploadTimeISO8601 上传时间（ISO 8601格式）
	UploadTimeISO8601 string `json:"upload_time_iso_8601"`

	// Digests 文件的哈希摘要信息
	Digests ReleaseDigests `json:"digests"`

	// MD5Digest MD5哈希值（旧字段，保留兼容性）
	MD5Digest string `json:"md5_digest"`

	// Yanked 此文件是否被撤回
	Yanked bool `json:"yanked"`

	// YankedReason 撤回原因
	YankedReason string `json:"yanked_reason,omitempty"`

	// CommentText 注释
	CommentText string `json:"comment_text"`
}

// ReleaseDigests 文件的哈希值信息
type ReleaseDigests struct {
	// MD5 MD5哈希值
	MD5 string `json:"md5"`

	// SHA256 SHA256哈希值
	SHA256 string `json:"sha256"`

	// Blake2b256 Blake2b-256哈希值
	Blake2b256 string `json:"blake2b_256"`
}

// GetUploadTimeISO 将上传时间解析为time.Time
// 优先使用ISO 8601格式，如果不可用则尝试旧格式
func (rf *ReleaseFile) GetUploadTimeISO() (time.Time, error) {
	if rf.UploadTimeISO8601 != "" {
		return time.Parse(time.RFC3339, rf.UploadTimeISO8601)
	}

	// 旧格式
	if rf.UploadTime != "" {
		return time.Parse("2006-01-02T15:04:05", rf.UploadTime)
	}

	return time.Time{}, nil
}

// IsYanked 检查发布文件是否被撤回
func (rf *ReleaseFile) IsYanked() bool {
	return rf.Yanked
}

// IsWheel 检查是否为wheel格式
func (rf *ReleaseFile) IsWheel() bool {
	return rf.PackageType == "bdist_wheel"
}

// IsSourceDist 检查是否为源码包格式
func (rf *ReleaseFile) IsSourceDist() bool {
	return rf.PackageType == "sdist"
}
