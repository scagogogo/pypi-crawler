# PyPI Crawler 示例

本目录包含了多个演示如何使用 PyPI Crawler 库的示例代码。这些示例展示了库的主要功能和使用方法。

## 示例列表

1. **index/** - 演示如何获取PyPI索引中的所有包
2. **package/** - 演示如何获取特定包的详细信息
3. **repository/** - 演示如何创建和配置不同镜像源的客户端实例
4. **search/** - 演示如何搜索PyPI包
5. **combined/** - 综合示例，演示所有API功能的使用

## 运行示例

### 索引示例 (index)

获取PyPI索引中的所有包，并显示包的总数和前10个包名：

```bash
go run examples/index/main.go
```

### 包信息示例 (package)

获取"requests"包的详细信息，包括版本、作者、摘要等：

```bash
go run examples/package/main.go
```

### 仓库示例 (repository)

演示如何创建使用不同镜像源的PyPI客户端：

```bash
go run examples/repository/main.go
```

### 搜索示例 (search)

根据关键词搜索PyPI包：

```bash
# 使用默认关键词 "python"
go run examples/search/main.go

# 使用自定义关键词
go run examples/search/main.go django
```

### 综合示例 (combined)

一个命令行工具，支持多种操作：

```bash
# 获取帮助
go run examples/combined/main.go

# 获取包信息
go run examples/combined/main.go info requests

# 获取特定版本信息
go run examples/combined/main.go version requests 2.28.1

# 列出包的所有版本
go run examples/combined/main.go releases requests

# 搜索包
go run examples/combined/main.go search django

# 列出所有包
go run examples/combined/main.go list

# 检查包漏洞
go run examples/combined/main.go check requests 2.28.1
```

## 编译并运行为可执行文件

可以将这些示例编译为独立的可执行文件：

```bash
# 构建所有示例
go build -o bin/index examples/index/main.go
go build -o bin/package examples/package/main.go
go build -o bin/repository examples/repository/main.go
go build -o bin/search examples/search/main.go
go build -o bin/pypi-tool examples/combined/main.go

# 运行编译后的可执行文件
./bin/pypi-tool info requests
```

## API 设计说明

库采用了清晰的分层设计：

1. **api 包**：定义了 `PyPIClient` 接口，为所有客户端实现提供统一约束
2. **client 包**：实现了 `PyPIClient` 接口的具体客户端
3. **mirrors 包**：提供了便捷的工厂函数，创建连接到不同镜像源的客户端
4. **models 包**：定义了数据模型，如 Package、ReleaseFile、Vulnerability 等

主要API功能包括：

- **GetPackageInfo**：获取包的最新信息
- **GetPackageVersion**：获取包的特定版本信息
- **GetPackageReleases**：获取包的所有版本列表
- **CheckPackageVulnerabilities**：检查包版本是否存在已知漏洞
- **GetAllPackages**：获取所有包的列表
- **GetPackageList**：获取所有包的集合（以 map 形式返回）
- **SearchPackages**：根据关键词搜索包 