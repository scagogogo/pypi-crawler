# 文档站点部署指南

本指南将帮助您设置和部署 PyPI Crawler 的文档站点到 GitHub Pages。

## 🚀 自动部署设置

### 1. 启用 GitHub Pages

1. 进入 GitHub 仓库页面
2. 点击 **Settings** 标签
3. 在左侧菜单中找到 **Pages**
4. 在 **Source** 部分选择 **Deploy from a branch**
5. 选择 **gh-pages** 分支和 **/ (root)** 文件夹

### 2. 配置权限

确保 GitHub Actions 有足够的权限：

1. 在仓库设置中，进入 **Actions** > **General**
2. 在 **Workflow permissions** 部分选择：
   - ✅ **Read and write permissions**
   - ✅ **Allow GitHub Actions to create and approve pull requests**

**注意**: 这是必需的，因为工作流需要推送到 `gh-pages` 分支。

### 3. 触发部署

部署会在以下情况自动触发：
- 推送到 `main` 分支且 `docs/` 目录有变更
- 手动触发（在 Actions 页面）

### 4. 部署流程

1. **构建**: GitHub Actions 在 `docs/` 目录中运行 `npm run docs:build`
2. **推送**: 将构建产物推送到 `gh-pages` 分支
3. **部署**: GitHub Pages 从 `gh-pages` 分支自动部署网站

## 🔧 本地开发

### 环境要求

- Node.js 16+ 
- npm 或 yarn

### 快速开始

```bash
# 1. 运行设置脚本
./scripts/setup-docs.sh

# 2. 启动开发服务器
npm run docs:dev
```

### 可用命令

```bash
# 开发模式（热重载）
npm run docs:dev

# 构建生产版本
npm run docs:build

# 预览生产版本
npm run docs:preview
```

## 📁 文件结构

```
docs/
├── .vitepress/
│   ├── config.js          # VitePress 配置
│   └── dist/              # 构建输出（自动生成）
├── public/
│   ├── logo.svg           # 网站 Logo
│   └── favicon.ico        # 网站图标
├── *.md                   # 文档页面
└── index.md               # 首页
```

## 🎨 自定义配置

### 修改网站信息

编辑 `docs/.vitepress/config.js`：

```js
export default defineConfig({
  title: 'PyPI Crawler',              // 网站标题
  description: '强大的 PyPI 包信息获取库', // 网站描述
  base: '/pypi-crawler/',             // 基础路径
  // ... 其他配置
})
```

### 修改主题

VitePress 支持自定义主题，您可以：

1. 修改配置文件中的 `themeConfig`
2. 添加自定义 CSS
3. 使用 Vue 组件

### 添加新页面

1. 在 `docs/` 目录下创建新的 `.md` 文件
2. 在 `config.js` 的 `sidebar` 中添加链接
3. 推送到 GitHub，自动部署

## 🌐 访问地址

- **开发环境**: http://localhost:5173
- **预览环境**: http://localhost:4173  
- **生产环境**: https://scagogogo.github.io/pypi-crawler/

## 🔍 故障排除

### 部署失败

1. 检查 GitHub Actions 日志
2. 确认 Node.js 版本兼容性
3. 验证 `package.json` 配置

### 页面显示异常

1. 检查 Markdown 语法
2. 验证内部链接路径
3. 确认图片资源路径

### 样式问题

1. 清除浏览器缓存
2. 检查 CSS 自定义
3. 验证主题配置

## 📞 获取帮助

如果遇到问题：

1. 查看 [VitePress 官方文档](https://vitepress.dev/)
2. 检查 [GitHub Actions 文档](https://docs.github.com/en/actions)
3. 在项目中提交 Issue

## 🎯 最佳实践

1. **定期更新依赖**: 保持 VitePress 版本最新
2. **优化图片**: 使用适当大小的图片资源
3. **测试链接**: 确保所有内部链接正常工作
4. **移动端适配**: 测试移动设备上的显示效果
5. **SEO 优化**: 设置合适的页面标题和描述

---

**恭喜！** 您的文档站点现在已经配置完成，可以自动部署到 GitHub Pages 了！
