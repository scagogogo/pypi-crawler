# 数据模型

本文档详细介绍了 PyPI Crawler 中使用的所有数据结构。

## 📋 目录

- [Package - 包信息](#package---包信息)
- [PackageInfo - 包元数据](#packageinfo---包元数据)
- [ReleaseFile - 发布文件](#releasefile---发布文件)
- [ReleaseDigests - 文件哈希](#releasedigests---文件哈希)
- [Vulnerability - 安全漏洞](#vulnerability---安全漏洞)

## Package - 包信息

`Package` 是从 PyPI API 获取的完整包信息的根结构。

```go
type Package struct {
    Info            *PackageInfo             `json:"info"`
    LastSerial      int                      `json:"last_serial"`
    Releases        map[string][]*ReleaseFile `json:"releases"`
    Urls            []*ReleaseFile           `json:"urls"`
    Vulnerabilities []Vulnerability          `json:"vulnerabilities"`
}
```

### 字段说明

| 字段 | 类型 | 描述 |
|------|------|------|
| `Info` | `*PackageInfo` | 包的基本元数据信息 |
| `LastSerial` | `int` | 包的最后序列号，用于增量更新 |
| `Releases` | `map[string][]*ReleaseFile` | 所有版本的发布文件，键为版本号 |
| `Urls` | `[]*ReleaseFile` | 最新版本的发布文件列表 |
| `Vulnerabilities` | `[]Vulnerability` | 已知的安全漏洞信息 |

### 使用示例

```go
pkg, err := client.GetPackageInfo(ctx, "requests")
if err != nil {
    log.Fatal(err)
}

// 访问基本信息
fmt.Printf("包名: %s\n", pkg.Info.Name)
fmt.Printf("当前版本: %s\n", pkg.Info.Version)

// 访问所有版本
fmt.Printf("总版本数: %d\n", len(pkg.Releases))

// 访问最新版本的文件
fmt.Printf("最新版本文件数: %d\n", len(pkg.Urls))

// 检查漏洞
if len(pkg.Vulnerabilities) > 0 {
    fmt.Printf("发现 %d 个漏洞\n", len(pkg.Vulnerabilities))
}
```

## PackageInfo - 包元数据

`PackageInfo` 包含包的详细元数据信息。

```go
type PackageInfo struct {
    Name                   string            `json:"name"`
    Version                string            `json:"version"`
    Summary                string            `json:"summary"`
    Description            string            `json:"description"`
    DescriptionContentType string            `json:"description_content_type"`
    Author                 string            `json:"author"`
    AuthorEmail            string            `json:"author_email"`
    Maintainer             string            `json:"maintainer"`
    MaintainerEmail        string            `json:"maintainer_email"`
    License                string            `json:"license"`
    Keywords               string            `json:"keywords"`
    ClassifiersArray       []string          `json:"classifiers"`
    ProjectURL             string            `json:"project_url"`
    ProjectURLs            map[string]string `json:"project_urls"`
    RequiresDist           []string          `json:"requires_dist"`
    RequiresPython         string            `json:"requires_python"`
    HomePage               string            `json:"home_page"`
    DocsURL                string            `json:"docs_url"`
    DownloadURL            string            `json:"download_url"`
    Yanked                 bool              `json:"yanked"`
    YankedReason           string            `json:"yanked_reason,omitempty"`
}
```

### 主要字段说明

| 字段 | 类型 | 描述 |
|------|------|------|
| `Name` | `string` | 包名 |
| `Version` | `string` | 当前版本号 |
| `Summary` | `string` | 包的简短描述 |
| `Description` | `string` | 包的详细描述（通常是 README 内容） |
| `Author` | `string` | 作者姓名 |
| `AuthorEmail` | `string` | 作者邮箱 |
| `License` | `string` | 许可证信息 |
| `RequiresDist` | `[]string` | 依赖的其他包列表 |
| `RequiresPython` | `string` | Python 版本要求 |
| `ProjectURLs` | `map[string]string` | 项目相关链接 |
| `Yanked` | `bool` | 是否被撤回 |

### 便捷方法

```go
// 获取所有依赖
dependencies := pkg.Info.GetAllDependencies()

// 检查是否有 Python 版本要求
if pkg.Info.HasPythonRequirement() {
    fmt.Printf("需要 Python: %s\n", pkg.Info.RequiresPython)
}

// 检查是否被撤回
if pkg.Info.IsYanked() {
    fmt.Printf("包已被撤回: %s\n", pkg.Info.YankedReason)
}

// 获取项目链接
urls := pkg.Info.GetProjectURLs()
for name, url := range urls {
    fmt.Printf("%s: %s\n", name, url)
}
```

## ReleaseFile - 发布文件

`ReleaseFile` 表示包的一个发布文件（如 wheel 文件或源码包）。

```go
type ReleaseFile struct {
    Filename          string         `json:"filename"`
    URL               string         `json:"url"`
    PackageType       string         `json:"packagetype"`
    PythonVersion     string         `json:"python_version"`
    RequiresPython    string         `json:"requires_python"`
    Size              int64          `json:"size"`
    UploadTime        string         `json:"upload_time"`
    UploadTimeISO8601 string         `json:"upload_time_iso_8601"`
    Digests           ReleaseDigests `json:"digests"`
    MD5Digest         string         `json:"md5_digest"`
    Yanked            bool           `json:"yanked"`
    YankedReason      string         `json:"yanked_reason,omitempty"`
    CommentText       string         `json:"comment_text"`
}
```

### 字段说明

| 字段 | 类型 | 描述 |
|------|------|------|
| `Filename` | `string` | 文件名 |
| `URL` | `string` | 下载链接 |
| `PackageType` | `string` | 包类型（sdist, bdist_wheel 等） |
| `Size` | `int64` | 文件大小（字节） |
| `UploadTime` | `string` | 上传时间（旧格式） |
| `UploadTimeISO8601` | `string` | 上传时间（ISO 8601 格式） |
| `Digests` | `ReleaseDigests` | 文件哈希值 |
| `Yanked` | `bool` | 是否被撤回 |

### 便捷方法

```go
// 解析上传时间
uploadTime, err := file.GetUploadTimeISO()
if err == nil {
    fmt.Printf("上传时间: %s\n", uploadTime.Format("2006-01-02 15:04:05"))
}

// 检查文件类型
if file.IsWheel() {
    fmt.Println("这是一个 wheel 文件")
}

if file.IsSourceDist() {
    fmt.Println("这是一个源码包")
}

// 检查是否被撤回
if file.IsYanked() {
    fmt.Printf("文件已被撤回: %s\n", file.YankedReason)
}
```

### 包类型说明

| 类型 | 描述 |
|------|------|
| `sdist` | 源码分发包（通常是 .tar.gz 文件） |
| `bdist_wheel` | 二进制 wheel 包（.whl 文件） |
| `bdist_egg` | 旧式 egg 包（已废弃） |

## ReleaseDigests - 文件哈希

`ReleaseDigests` 包含文件的各种哈希值，用于验证文件完整性。

```go
type ReleaseDigests struct {
    MD5        string `json:"md5"`
    SHA256     string `json:"sha256"`
    Blake2b256 string `json:"blake2b_256"`
}
```

### 使用示例

```go
for _, file := range pkg.Urls {
    fmt.Printf("文件: %s\n", file.Filename)
    fmt.Printf("MD5: %s\n", file.Digests.MD5)
    fmt.Printf("SHA256: %s\n", file.Digests.SHA256)
    fmt.Printf("Blake2b-256: %s\n", file.Digests.Blake2b256)
}
```

## Vulnerability - 安全漏洞

`Vulnerability` 表示包的一个安全漏洞信息。

```go
type Vulnerability struct {
    ID        string   `json:"id"`
    Aliases   []string `json:"aliases"`
    Summary   string   `json:"summary"`
    Details   string   `json:"details"`
    FixedIn   []string `json:"fixed_in"`
    Source    string   `json:"source"`
    Link      string   `json:"link"`
    Withdrawn string   `json:"withdrawn"`
}
```

### 字段说明

| 字段 | 类型 | 描述 |
|------|------|------|
| `ID` | `string` | 漏洞唯一标识符 |
| `Aliases` | `[]string` | 漏洞别名（如 CVE 编号） |
| `Summary` | `string` | 漏洞摘要 |
| `Details` | `string` | 漏洞详细描述 |
| `FixedIn` | `[]string` | 已修复的版本列表 |
| `Source` | `string` | 漏洞信息来源 |
| `Link` | `string` | 漏洞详情链接 |
| `Withdrawn` | `string` | 撤回时间（如果已撤回） |

### 便捷方法

```go
for _, vuln := range pkg.Vulnerabilities {
    fmt.Printf("漏洞 ID: %s\n", vuln.ID)
    fmt.Printf("摘要: %s\n", vuln.Summary)
    
    // 检查是否有 CVE 编号
    if vuln.HasCVE() {
        cves := vuln.GetCVEs()
        fmt.Printf("CVE 编号: %v\n", cves)
    }
    
    // 检查特定版本是否已修复
    if vuln.IsFixed("2.28.0") {
        fmt.Println("版本 2.28.0 已修复此漏洞")
    }
    
    // 检查是否已撤回
    if vuln.IsWithdrawn() {
        withdrawnTime, _ := vuln.GetWithdrawnTime()
        fmt.Printf("漏洞已于 %s 撤回\n", withdrawnTime.Format("2006-01-02"))
    }
}
```

## JSON 示例

### Package 结构示例

```json
{
  "info": {
    "name": "requests",
    "version": "2.31.0",
    "summary": "Python HTTP for Humans.",
    "author": "Kenneth Reitz",
    "license": "Apache 2.0",
    "requires_dist": [
      "charset-normalizer (<4,>=2)",
      "idna (<4,>=2.5)",
      "urllib3 (<3,>=1.21.1)",
      "certifi (>=2017.4.17)"
    ]
  },
  "urls": [
    {
      "filename": "requests-2.31.0-py3-none-any.whl",
      "url": "https://files.pythonhosted.org/packages/.../requests-2.31.0-py3-none-any.whl",
      "packagetype": "bdist_wheel",
      "size": 62574,
      "digests": {
        "md5": "...",
        "sha256": "..."
      }
    }
  ],
  "vulnerabilities": []
}
```

---

**下一步**: 查看 [客户端配置](./client-configuration.md) 了解如何配置客户端选项。
