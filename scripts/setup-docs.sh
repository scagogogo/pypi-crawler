#!/bin/bash

# PyPI Crawler 文档站点设置脚本
# 此脚本帮助您快速设置 VitePress 文档站点

set -e

echo "🚀 PyPI Crawler 文档站点设置"
echo "================================"

# 检查 Node.js 是否安装
if ! command -v node &> /dev/null; then
    echo "❌ 错误: 未找到 Node.js"
    echo "请先安装 Node.js (版本 16 或更高): https://nodejs.org/"
    exit 1
fi

# 检查 Node.js 版本
NODE_VERSION=$(node -v | cut -d'v' -f2 | cut -d'.' -f1)
if [ "$NODE_VERSION" -lt 16 ]; then
    echo "❌ 错误: Node.js 版本过低 (当前: $(node -v))"
    echo "请升级到 Node.js 16 或更高版本"
    exit 1
fi

echo "✅ Node.js 版本: $(node -v)"

# 检查 npm 是否可用
if ! command -v npm &> /dev/null; then
    echo "❌ 错误: 未找到 npm"
    exit 1
fi

echo "✅ npm 版本: $(npm -v)"

# 安装依赖
echo ""
echo "📦 安装依赖..."
npm install

echo ""
echo "✅ 依赖安装完成!"

# 提供使用说明
echo ""
echo "🎉 设置完成! 现在您可以:"
echo ""
echo "  开发模式 (热重载):"
echo "    npm run docs:dev"
echo ""
echo "  构建生产版本:"
echo "    npm run docs:build"
echo ""
echo "  预览生产版本:"
echo "    npm run docs:preview"
echo ""
echo "📚 文档将在以下地址可用:"
echo "  开发: http://localhost:5173"
echo "  预览: http://localhost:4173"
echo ""
echo "🔗 GitHub Pages 部署:"
echo "  推送到 main 分支后，GitHub Actions 会自动部署到:"
echo "  https://scagogogo.github.io/pypi-crawler/"
echo ""
echo "💡 提示: 确保在 GitHub 仓库设置中启用 GitHub Pages"
