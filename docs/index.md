---
layout: home

hero:
  name: "PyPI Crawler"
  text: "强大的 PyPI 包信息获取库"
  tagline: "使用 Go 语言构建，支持多镜像源，提供完整的 PyPI API 访问能力"
  image:
    src: /logo.svg
    alt: PyPI Crawler
  actions:
    - theme: brand
      text: 快速开始
      link: /quick-start
    - theme: alt
      text: API 文档
      link: /api-reference
    - theme: alt
      text: GitHub
      link: https://github.com/scagogogo/pypi-crawler

features:
  - icon: 🚀
    title: 简单易用
    details: 提供简洁的 API 接口，几行代码即可获取 PyPI 包信息，支持链式配置。
  - icon: 🌐
    title: 多镜像源支持
    details: 内置支持官方源和7个国内镜像源，自动选择最快的访问路径。
  - icon: 🛡️
    title: 安全可靠
    details: 完善的错误处理机制，自动重试，支持安全漏洞检查功能。
  - icon: ⚡
    title: 高性能
    details: 并发安全，支持批量操作，内置缓存机制，响应速度快。
  - icon: 🔧
    title: 灵活配置
    details: 丰富的配置选项，支持代理、超时、重试等自定义设置。
  - icon: 📚
    title: 完整文档
    details: 提供详细的 API 文档、示例代码和最佳实践指南。
---

## 🎯 核心功能

### 📦 包信息获取
- 获取包的最新版本信息
- 获取包的特定版本信息  
- 获取包的所有发布版本列表

### 🔍 搜索功能
- 根据关键词搜索包
- 获取所有包的索引列表

### 🛡️ 安全功能
- 检查包的已知安全漏洞
- 获取漏洞详细信息和修复版本

### 🌐 多镜像源支持
- PyPI 官方源
- 清华大学镜像源
- 阿里云镜像源
- 豆瓣镜像源
- 腾讯云镜像源
- 中国科技大学镜像源
- 网易镜像源

## 🚀 快速体验

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/scagogogo/pypi-crawler/pkg/pypi/mirrors"
)

func main() {
    // 创建客户端（使用清华镜像源）
    client := mirrors.NewTsinghuaClient()
    
    // 获取包信息
    pkg, err := client.GetPackageInfo(context.Background(), "requests")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("包名: %s\n", pkg.Info.Name)
    fmt.Printf("版本: %s\n", pkg.Info.Version)
    fmt.Printf("摘要: %s\n", pkg.Info.Summary)
}
```

## 📊 项目统计

- **🎯 API 方法**: 8 个核心方法
- **🌐 镜像源**: 支持 7+ 个镜像源
- **📝 文档页面**: 11 页详细文档
- **💡 代码示例**: 50+ 个实用示例
- **🧪 测试覆盖率**: 80.9%
- **⭐ Go 版本**: 支持 Go 1.19+

## 🏗️ 架构设计

```
pkg/pypi/
├── api/            # API 接口定义
├── client/         # 客户端实现
├── mirrors/        # 镜像源工厂
└── models/         # 数据模型
```

## 🤝 社区

- **GitHub**: [scagogogo/pypi-crawler](https://github.com/scagogogo/pypi-crawler)
- **Issues**: [提交问题](https://github.com/scagogogo/pypi-crawler/issues)
- **Discussions**: [参与讨论](https://github.com/scagogogo/pypi-crawler/discussions)

## 📄 许可证

本项目使用 [MIT 许可证](https://github.com/scagogogo/pypi-crawler/blob/main/LICENSE)。
