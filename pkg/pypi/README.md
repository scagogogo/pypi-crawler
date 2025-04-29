# PyPI API 单元测试

本目录包含 PyPI API 的单元测试套件，覆盖了所有主要功能。测试使用模拟服务器避免实际网络请求。

## 目录结构

```
pkg/pypi/
├── api/            - API 接口定义
├── client/         - API 实现
│   ├── testdata/   - 模拟 API 响应
│   └── client_test.go - 客户端测试
├── mirrors/        - 镜像源工厂
├── models/         - 数据模型
```

## 运行测试

### 运行所有测试

```bash
# 在项目根目录运行
go test ./pkg/pypi/... -v
```

### 运行特定包的测试

```bash
# 测试客户端
go test ./pkg/pypi/client -v

# 测试模型
go test ./pkg/pypi/models -v

# 测试镜像源
go test ./pkg/pypi/mirrors -v
```

### 运行特定测试

```bash
# 运行特定测试函数
go test ./pkg/pypi/client -v -run TestGetPackageInfo

# 运行特定测试子函数
go test ./pkg/pypi/client -v -run "TestGetPackageInfo/获取存在的包"
```

## 生成测试覆盖率报告

### 命令行报告

```bash
# 运行测试并生成覆盖率信息
go test ./pkg/pypi/... -coverprofile=coverage.out

# 显示覆盖率统计
go tool cover -func=coverage.out
```

### HTML 报告

```bash
# 生成 HTML 格式的可视化覆盖率报告
go test ./pkg/pypi/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## 测试数据

测试使用 `testdata` 目录中的模拟数据。如果目录不存在或文件缺失，测试会自动创建必要的测试数据。

生成的测试数据文件已被 `.gitignore` 忽略，不会提交到代码库中。

## 注意事项

1. 所有测试都不依赖外部网络连接，使用模拟HTTP服务器
2. 测试覆盖了正常和异常情况，如不存在的包、网络错误等
3. 模型测试不依赖于客户端实现，确保接口和实现的独立性 