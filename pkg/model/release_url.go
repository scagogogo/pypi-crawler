package model

import (
	"context"
	"fmt"
	"time"

	"github.com/crawler-go-go-go/go-requests"
)

// ReleaseURL 表示Python包的一个特定发布版本的下载信息
//
// 包含了该版本包的文件名、下载地址、大小、校验和等信息
//
// 使用示例:
//
//	pkg, err := repo.GetPackage(context.Background(), "requests")
//	if err == nil && len(pkg.Urls) > 0 {
//	    // 获取最新版本的第一个发布文件
//	    releaseURL := pkg.Urls[0]
//	    fmt.Printf("文件名: %s\n", releaseURL.Filename)
//	    fmt.Printf("下载地址: %s\n", releaseURL.URL)
//	    fmt.Printf("文件大小: %d 字节\n", releaseURL.Size)
//	    fmt.Printf("上传时间: %s\n", releaseURL.UploadTimeIso8601)
//	    fmt.Printf("包类型: %s\n", releaseURL.Packagetype)
//
//	    // 可以直接下载文件
//	    data, err := releaseURL.Download(context.Background())
//	    if err == nil {
//	        fmt.Printf("下载成功，文件大小: %d 字节\n", len(data))
//	    }
//	}
type ReleaseURL struct {
	// CommentText 发布说明或注释
	CommentText string `json:"comment_text"`

	// Digests 文件的各种哈希摘要
	// 包含不同算法的哈希值，用于校验文件完整性
	Digests *Digests `json:"digests"`

	// Downloads 此文件的下载次数
	Downloads int `json:"downloads"`

	// Filename 文件名
	// 通常包含包名、版本和格式信息
	Filename string `json:"filename"`

	// HasSig 是否有GPG签名文件
	HasSig bool `json:"has_sig"`

	// Md5Digest 文件的MD5哈希值
	// 用于校验下载文件的完整性
	Md5Digest string `json:"md5_digest"`

	// Packagetype 包的类型
	// 常见值如 "bdist_wheel", "sdist" 等
	Packagetype string `json:"packagetype"`

	// PythonVersion 适用的Python版本
	// 如 "py3", "py2.py3" 等
	PythonVersion string `json:"python_version"`

	// RequiresPython 要求的Python版本范围
	// 如 ">=3.6,<4.0" 等
	RequiresPython string `json:"requires_python"`

	// Size 文件大小，单位为字节
	Size int `json:"size"`

	// UploadTime 上传时间（字符串格式）
	UploadTime string `json:"upload_time"`

	// UploadTimeIso8601 上传时间（ISO 8601格式的时间对象）
	UploadTimeIso8601 time.Time `json:"upload_time_iso_8601"`

	// URL 文件的下载地址
	URL string `json:"url"`

	// Yanked 是否已被撤回
	// 被撤回的版本通常是因为发现严重bug或安全问题
	Yanked bool `json:"yanked"`

	// YankedReason 版本被撤回的原因
	// 如果Yanked为true，这里可能包含撤回原因
	YankedReason any `json:"yanked_reason"`
}

// Download 下载此发布版本的文件
//
// 此方法会使用HTTP GET请求下载URL指向的文件，并返回文件内容的字节数组
//
// 参数:
//   - ctx: 上下文，用于控制请求的生命周期
//
// 返回值:
//   - []byte: 下载文件的字节内容
//   - error: 下载过程中的错误，包括URL为空、网络错误等
//
// 使用示例:
//
//	releaseURL := &model.ReleaseURL{URL: "https://files.pythonhosted.org/packages/..."}
//	data, err := releaseURL.Download(context.Background())
//	if err != nil {
//	    log.Fatalf("下载失败: %v", err)
//	}
//	fmt.Printf("下载成功，文件大小: %d 字节\n", len(data))
//
//	// 也可以将内容保存到文件
//	if err := os.WriteFile("downloaded_package.whl", data, 0644); err != nil {
//	    log.Fatalf("保存文件失败: %v", err)
//	}
func (x *ReleaseURL) Download(ctx context.Context) ([]byte, error) {
	if x.URL == "" {
		return nil, fmt.Errorf("do not has release file url")
	}
	return requests.GetBytes(ctx, x.URL)
}

// Digests 包含文件的多种哈希摘要值
//
// 用于校验下载文件的完整性和真实性
type Digests struct {
	// Blake2b_256 BLAKE2b-256哈希值
	Blake2b_256 string `json:"blake2b_256"`

	// MD5 MD5哈希值（已被认为不安全，但仍广泛使用）
	MD5 string `json:"md5"`

	// Sha256 SHA-256哈希值（推荐使用的更安全的哈希算法）
	Sha256 string `json:"sha256"`
}
