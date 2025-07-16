import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'PyPI Crawler',
  description: 'PyPI Crawler - 强大的 PyPI 包信息获取库',
  
  // 基础配置
  base: '/pypi-crawler/',
  lang: 'zh-CN',

  // 忽略死链接检查
  ignoreDeadLinks: [
    // 忽略本地开发服务器链接
    /^http:\/\/localhost/,
    // 忽略相对路径链接
    /^\.\.?\//
  ],
  
  // 主题配置
  themeConfig: {
    // 网站标题和Logo
    siteTitle: 'PyPI Crawler',
    logo: '/logo.svg',
    
    // 导航栏
    nav: [
      { text: '首页', link: '/' },
      { text: '指南', link: '/quick-start' },
      { text: 'API', link: '/api-reference' },
      { text: '示例代码', link: '/examples' },
      { text: '最佳实践', link: '/best-practices' },
      { text: '常见问题', link: '/faq' },
      {
        text: 'GitHub',
        link: 'https://github.com/scagogogo/pypi-crawler'
      }
    ],
    
    // 侧边栏
    sidebar: [
      {
        text: '开始使用',
        items: [
          { text: '介绍', link: '/' },
          { text: '快速开始', link: '/quick-start' }
        ]
      },
      {
        text: 'API 文档',
        items: [
          { text: 'API 参考', link: '/api-reference' },
          { text: '数据模型', link: '/data-models' },
          { text: '客户端配置', link: '/client-configuration' }
        ]
      },
      {
        text: '配置指南',
        items: [
          { text: '镜像源配置', link: '/mirrors' },
          { text: '错误处理', link: '/error-handling' }
        ]
      },
      {
        text: '进阶使用',
        items: [
          { text: '示例代码', link: '/examples' },
          { text: '最佳实践', link: '/best-practices' }
        ]
      },
      {
        text: '其他',
        items: [
          { text: '常见问题', link: '/faq' },
          { text: '更新日志', link: '/changelog' }
        ]
      }
    ],
    
    // 社交链接
    socialLinks: [
      { icon: 'github', link: 'https://github.com/scagogogo/pypi-crawler' }
    ],
    
    // 页脚
    footer: {
      message: '基于 MIT 许可证发布',
      copyright: 'Copyright © 2024 scagogogo'
    },
    
    // 搜索
    search: {
      provider: 'local'
    },
    
    // 编辑链接
    editLink: {
      pattern: 'https://github.com/scagogogo/pypi-crawler/edit/main/docs/:path',
      text: '在 GitHub 上编辑此页'
    },
    
    // 最后更新时间
    lastUpdated: {
      text: '最后更新于',
      formatOptions: {
        dateStyle: 'short',
        timeStyle: 'medium'
      }
    },
    
    // 文档页脚导航
    docFooter: {
      prev: '上一页',
      next: '下一页'
    },
    
    // 大纲标题
    outline: {
      label: '页面导航'
    },
    
    // 返回顶部
    returnToTopLabel: '回到顶部',
    
    // 侧边栏菜单标签
    sidebarMenuLabel: '菜单',
    
    // 深色模式切换标签
    darkModeSwitchLabel: '主题',
    lightModeSwitchTitle: '切换到浅色模式',
    darkModeSwitchTitle: '切换到深色模式'
  },
  
  // Markdown 配置
  markdown: {
    // 代码块行号
    lineNumbers: true,
    
    // 代码块主题
    theme: {
      light: 'github-light',
      dark: 'github-dark'
    }
  },
  
  // 头部配置
  head: [
    ['link', { rel: 'icon', href: '/pypi-crawler/favicon.ico' }],
    ['meta', { name: 'theme-color', content: '#3c8772' }],
    ['meta', { property: 'og:type', content: 'website' }],
    ['meta', { property: 'og:locale', content: 'zh-CN' }],
    ['meta', { property: 'og:title', content: 'PyPI Crawler | PyPI 包信息获取库' }],
    ['meta', { property: 'og:site_name', content: 'PyPI Crawler' }],
    ['meta', { property: 'og:image', content: 'https://scagogogo.github.io/pypi-crawler/og-image.png' }],
    ['meta', { property: 'og:url', content: 'https://scagogogo.github.io/pypi-crawler/' }]
  ]
})
