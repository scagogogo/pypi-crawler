#!/bin/bash

# PyPI Crawler 文档测试脚本
# 此脚本用于测试文档构建和预览功能

set -e

echo "🧪 PyPI Crawler 文档测试"
echo "========================"

# 检查依赖是否已安装
if [ ! -d "node_modules" ]; then
    echo "❌ 错误: 未找到 node_modules 目录"
    echo "请先运行: npm install 或 ./scripts/setup-docs.sh"
    exit 1
fi

echo "✅ 依赖检查通过"

# 测试构建
echo ""
echo "🔨 测试文档构建..."
if npm run docs:build; then
    echo "✅ 构建测试通过"
else
    echo "❌ 构建测试失败"
    exit 1
fi

# 检查构建输出
if [ -d "docs/.vitepress/dist" ]; then
    echo "✅ 构建输出目录存在"
    
    # 检查关键文件
    if [ -f "docs/.vitepress/dist/index.html" ]; then
        echo "✅ 首页文件存在"
    else
        echo "❌ 首页文件缺失"
        exit 1
    fi
    
    if [ -f "docs/.vitepress/dist/api-reference.html" ]; then
        echo "✅ API 文档文件存在"
    else
        echo "❌ API 文档文件缺失"
        exit 1
    fi
    
    # 统计生成的文件
    html_count=$(find docs/.vitepress/dist -name "*.html" | wc -l)
    echo "✅ 生成了 $html_count 个 HTML 页面"
    
else
    echo "❌ 构建输出目录不存在"
    exit 1
fi

# 清理构建输出（可选）
echo ""
echo "🧹 清理构建输出..."
rm -rf docs/.vitepress/dist
echo "✅ 清理完成"

echo ""
echo "🎉 所有测试通过!"
echo ""
echo "💡 下一步:"
echo "  开发模式: npm run docs:dev"
echo "  构建文档: npm run docs:build"
echo "  预览文档: npm run docs:preview"
echo ""
echo "🌐 GitHub Pages 部署:"
echo "  推送到 main 分支后自动部署到:"
echo "  https://scagogogo.github.io/pypi-crawler/"
