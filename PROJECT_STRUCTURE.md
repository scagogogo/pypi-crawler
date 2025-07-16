# 项目结构说明

本文档说明了 PyPI Crawler 项目的目录结构和组织方式。

## 📁 项目结构

```
pypi-crawler/
├── 📁 pkg/                    # Go 包源码
│   └── pypi/                  # PyPI 爬虫核心库
│       ├── api/               # API 接口定义
│       ├── client/            # 客户端实现
│       ├── mirrors/           # 镜像源工厂
│       └── models/            # 数据模型
├── 📁 examples/               # 示例代码
├── 📁 docs/                   # 文档站点 (独立的前端项目)
│   ├── 📄 package.json        # Node.js 项目配置
│   ├── 📄 package-lock.json   # 依赖锁定文件
│   ├── 📁 node_modules/       # Node.js 依赖
│   ├── 📁 scripts/            # 文档相关脚本
│   ├── 📁 .vitepress/         # VitePress 配置
│   ├── 📁 public/             # 静态资源
│   ├── 📄 *.md                # 文档页面
│   └── 📄 index.md            # 文档首页
├── 📁 .github/                # GitHub 配置
│   └── workflows/             # GitHub Actions
│       ├── 📄 test.yml        # Go 测试工作流
│       └── 📄 docs.yml        # 文档部署工作流
├── 📄 go.mod                  # Go 模块定义
├── 📄 go.sum                  # Go 依赖校验
├── 📄 README.md               # 项目说明
└── 📄 LICENSE                 # 许可证
```

## 🎯 设计原则

### 1. **清晰分离**
- **Go 项目**: 根目录专注于 Go 代码和配置
- **文档站点**: `docs/` 目录包含所有前端相关文件
- **避免混乱**: 前端依赖不污染 Go 项目根目录

### 2. **独立管理**
- **Go 依赖**: 通过 `go.mod` 管理
- **Node.js 依赖**: 通过 `docs/package.json` 管理
- **构建流程**: 分别在各自目录中进行

### 3. **自动化**
- **Go 测试**: 自动运行单元测试和覆盖率检查
- **文档部署**: 自动构建和部署到 GitHub Pages
- **多版本支持**: 测试多个 Go 版本的兼容性

## 🚀 开发工作流

### Go 开发
```bash
# 在项目根目录
go test ./pkg/pypi/...           # 运行测试
go build ./examples/...          # 构建示例
go mod tidy                      # 整理依赖
```

### 文档开发
```bash
# 快速设置
./docs/scripts/setup-docs.sh

# 开发模式
cd docs && npm run docs:dev

# 构建文档
cd docs && npm run docs:build

# 测试文档
./docs/scripts/test-docs.sh
```

## 🔄 CI/CD 流程

### Go 测试工作流 (`.github/workflows/test.yml`)
- **触发条件**: 推送到 main 分支 (排除 docs/ 目录)
- **测试矩阵**: Go 1.19, 1.20, 1.21
- **功能**: 
  - 运行单元测试
  - 生成覆盖率报告
  - 上传到 Codecov

### 文档部署工作流 (`.github/workflows/docs.yml`)
- **触发条件**: docs/ 目录有变更
- **功能**:
  - 构建 VitePress 文档
  - 部署到 GitHub Pages
  - 自动更新在线文档

## 📚 文档站点

### 技术栈
- **框架**: VitePress (基于 Vue.js 和 Vite)
- **主题**: 默认主题 + 自定义配置
- **部署**: GitHub Pages
- **域名**: https://scagogogo.github.io/pypi-crawler/

### 特性
- 🎨 美观的界面设计
- 🔍 本地搜索功能
- 📱 响应式布局
- 🌙 深色/浅色主题
- 🚀 快速加载速度
- 📖 中文界面

## 🛠️ 维护指南

### 添加新的 Go 功能
1. 在 `pkg/pypi/` 下添加代码
2. 编写单元测试
3. 更新相关文档
4. 提交代码，自动触发测试

### 更新文档
1. 编辑 `docs/` 下的 Markdown 文件
2. 本地预览: `cd docs && npm run docs:dev`
3. 提交代码，自动部署到 GitHub Pages

### 发布新版本
1. 更新 `docs/changelog.md`
2. 创建 Git 标签
3. 发布 GitHub Release

## 🔧 故障排除

### Go 测试失败
- 检查 GitHub Actions 日志
- 本地运行: `go test ./pkg/pypi/... -v`
- 确保所有依赖都是最新的

### 文档构建失败
- 检查 Node.js 版本 (需要 16+)
- 运行: `./docs/scripts/test-docs.sh`
- 检查 Markdown 语法和链接

### GitHub Pages 不更新
- 确保 GitHub Pages 设置为 "GitHub Actions"
- 检查工作流权限设置
- 查看 Actions 运行日志

## 📞 获取帮助

- **GitHub Issues**: 报告问题和建议
- **GitHub Discussions**: 讨论和交流
- **文档**: 查看在线文档获取详细信息

---

这种项目结构确保了 Go 项目的纯净性，同时提供了强大的文档站点功能。
