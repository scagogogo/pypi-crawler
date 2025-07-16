# PyPI Crawler API 文档

欢迎使用 PyPI Crawler API 文档！本文档详细介绍了如何使用这个强大的 PyPI 包信息获取库。

## 📚 文档目录

- [快速开始](./quick-start.md) - 5分钟快速上手指南
- [API 参考](./api-reference.md) - 完整的API接口文档
- [数据模型](./data-models.md) - 详细的数据结构说明
- [客户端配置](./client-configuration.md) - 客户端配置选项详解
- [镜像源配置](./mirrors.md) - 支持的镜像源列表和使用方法
- [错误处理](./error-handling.md) - 错误处理最佳实践
- [示例代码](./examples.md) - 丰富的使用示例
- [最佳实践](./best-practices.md) - 使用建议和性能优化
- [常见问题](./faq.md) - 常见问题解答
- [更新日志](./changelog.md) - 版本更新记录

## 🚀 核心功能

PyPI Crawler 提供以下核心功能：

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
- 官方 PyPI 源
- 国内镜像源（清华、豆瓣、阿里云等）
- 自定义镜像源

## 🏗️ 架构设计

```
pkg/pypi/
├── api/            # API 接口定义
├── client/         # 客户端实现
├── mirrors/        # 镜像源工厂
└── models/         # 数据模型
```

### 设计原则

1. **接口分离**: 通过 `api.PyPIClient` 接口定义统一的API规范
2. **工厂模式**: 通过 `mirrors` 包提供便捷的客户端创建方法
3. **配置灵活**: 支持链式配置和自定义选项
4. **错误处理**: 完善的错误处理和重试机制

## 📋 系统要求

- Go 1.19 或更高版本
- 网络连接（用于访问 PyPI API）

## 📦 安装

```bash
go get -u github.com/scagogogo/pypi-crawler
```

## 🔗 相关链接

- [GitHub 仓库](https://github.com/scagogogo/pypi-crawler)
- [PyPI 官方 API 文档](https://warehouse.pypi.org/api-reference/)
- [示例代码](../examples/)

## 📄 许可证

本项目使用 MIT 许可证。详情请查看 [LICENSE](../LICENSE) 文件。

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！请参考 [贡献指南](../README.md#六贡献指南) 了解详细信息。

---

**下一步**: 查看 [快速开始](./quick-start.md) 开始使用 PyPI Crawler！
