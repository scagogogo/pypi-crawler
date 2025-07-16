# 错误处理

本文档介绍 PyPI Crawler 中可能遇到的各种错误类型以及相应的处理策略。

## 📋 目录

- [错误类型概览](#错误类型概览)
- [网络错误](#网络错误)
- [HTTP 错误](#http-错误)
- [解析错误](#解析错误)
- [上下文错误](#上下文错误)
- [错误处理最佳实践](#错误处理最佳实践)
- [重试策略](#重试策略)
- [日志记录](#日志记录)

## 错误类型概览

PyPI Crawler 中的错误主要分为以下几类：

| 错误类型 | 描述 | 是否可重试 | 常见原因 |
|----------|------|------------|----------|
| 网络错误 | 连接失败、超时等 | ✅ | 网络不稳定、DNS 解析失败 |
| HTTP 4xx | 客户端错误 | ❌ | 包不存在、请求格式错误 |
| HTTP 5xx | 服务器错误 | ✅ | 服务器临时故障 |
| 解析错误 | JSON 解析失败 | ❌ | API 响应格式异常 |
| 上下文错误 | 超时、取消 | ❌ | 操作被取消或超时 |

## 网络错误

### 连接失败

```go
pkg, err := client.GetPackageInfo(ctx, "requests")
if err != nil {
    if strings.Contains(err.Error(), "connection refused") {
        fmt.Println("服务器拒绝连接，请检查网络或更换镜像源")
        return
    }
    
    if strings.Contains(err.Error(), "no such host") {
        fmt.Println("DNS 解析失败，请检查网络设置")
        return
    }
}
```

### 超时处理

```go
import (
    "context"
    "time"
    "net"
)

func handleTimeoutError(err error) {
    if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
        fmt.Println("请求超时，建议：")
        fmt.Println("1. 增加超时时间")
        fmt.Println("2. 检查网络连接")
        fmt.Println("3. 更换镜像源")
        return
    }
}

// 使用示例
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

pkg, err := client.GetPackageInfo(ctx, "requests")
if err != nil {
    handleTimeoutError(err)
}
```

### 网络错误重试

```go
func getPackageWithRetry(client api.PyPIClient, packageName string, maxRetries int) (*models.Package, error) {
    var lastErr error
    
    for attempt := 0; attempt < maxRetries; attempt++ {
        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        
        pkg, err := client.GetPackageInfo(ctx, packageName)
        cancel()
        
        if err == nil {
            return pkg, nil
        }
        
        lastErr = err
        
        // 检查是否为可重试的网络错误
        if !isRetryableNetworkError(err) {
            break
        }
        
        // 指数退避
        delay := time.Duration(attempt+1) * time.Second
        time.Sleep(delay)
    }
    
    return nil, fmt.Errorf("重试 %d 次后仍然失败: %w", maxRetries, lastErr)
}

func isRetryableNetworkError(err error) bool {
    if netErr, ok := err.(net.Error); ok {
        return netErr.Timeout() || netErr.Temporary()
    }
    
    // 检查其他可重试的错误
    errStr := err.Error()
    retryableErrors := []string{
        "connection refused",
        "connection reset",
        "no route to host",
        "network is unreachable",
    }
    
    for _, retryableErr := range retryableErrors {
        if strings.Contains(errStr, retryableErr) {
            return true
        }
    }
    
    return false
}
```

## HTTP 错误

### 404 错误 - 包不存在

```go
func handlePackageNotFound(err error, packageName string) {
    if strings.Contains(err.Error(), "404") {
        fmt.Printf("包 '%s' 不存在，可能的原因：\n", packageName)
        fmt.Println("1. 包名拼写错误")
        fmt.Println("2. 包已被删除")
        fmt.Println("3. 包名大小写不匹配")
        
        // 建议相似的包名
        suggestSimilarPackages(packageName)
    }
}

func suggestSimilarPackages(packageName string) {
    // 简单的相似包名建议逻辑
    suggestions := []string{
        strings.ToLower(packageName),
        strings.ReplaceAll(packageName, "_", "-"),
        strings.ReplaceAll(packageName, "-", "_"),
    }
    
    fmt.Println("建议尝试以下包名：")
    for _, suggestion := range suggestions {
        if suggestion != packageName {
            fmt.Printf("  - %s\n", suggestion)
        }
    }
}
```

### 403 错误 - 访问被拒绝

```go
func handle403Error(err error) {
    if strings.Contains(err.Error(), "403") {
        fmt.Println("访问被拒绝，可能的原因：")
        fmt.Println("1. IP 被限制")
        fmt.Println("2. User-Agent 被屏蔽")
        fmt.Println("3. 请求频率过高")
        fmt.Println("建议：")
        fmt.Println("1. 更换镜像源")
        fmt.Println("2. 设置合适的 User-Agent")
        fmt.Println("3. 降低请求频率")
    }
}
```

### 429 错误 - 请求频率限制

```go
import "strconv"

func handle429Error(err error, resp *http.Response) error {
    if strings.Contains(err.Error(), "429") {
        fmt.Println("请求频率过高，触发限流")
        
        // 尝试从响应头获取重试时间
        if resp != nil {
            if retryAfter := resp.Header.Get("Retry-After"); retryAfter != "" {
                if seconds, parseErr := strconv.Atoi(retryAfter); parseErr == nil {
                    fmt.Printf("建议等待 %d 秒后重试\n", seconds)
                    time.Sleep(time.Duration(seconds) * time.Second)
                    return nil
                }
            }
        }
        
        // 默认等待时间
        fmt.Println("等待 60 秒后重试...")
        time.Sleep(60 * time.Second)
    }
    
    return err
}
```

### 5xx 服务器错误

```go
func handle5xxError(err error) bool {
    errStr := err.Error()
    
    serverErrors := []string{"500", "502", "503", "504"}
    for _, code := range serverErrors {
        if strings.Contains(errStr, code) {
            fmt.Printf("服务器错误 (%s)，这通常是临时性问题\n", code)
            fmt.Println("建议：")
            fmt.Println("1. 稍后重试")
            fmt.Println("2. 更换镜像源")
            return true // 表示这是服务器错误
        }
    }
    
    return false
}
```

## 解析错误

### JSON 解析失败

```go
import "encoding/json"

func handleJSONError(err error, responseBody []byte) {
    if jsonErr, ok := err.(*json.SyntaxError); ok {
        fmt.Printf("JSON 解析错误，位置: %d\n", jsonErr.Offset)
        
        // 显示错误附近的内容
        start := max(0, int(jsonErr.Offset)-50)
        end := min(len(responseBody), int(jsonErr.Offset)+50)
        
        fmt.Printf("错误附近的内容: %s\n", string(responseBody[start:end]))
        
        fmt.Println("可能的原因：")
        fmt.Println("1. API 响应格式发生变化")
        fmt.Println("2. 网络传输过程中数据损坏")
        fmt.Println("3. 镜像源返回了非标准响应")
    }
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}
```

### 数据验证错误

```go
func validatePackageData(pkg *models.Package) error {
    if pkg == nil {
        return fmt.Errorf("包数据为空")
    }
    
    if pkg.Info == nil {
        return fmt.Errorf("包信息为空")
    }
    
    if pkg.Info.Name == "" {
        return fmt.Errorf("包名为空")
    }
    
    if pkg.Info.Version == "" {
        return fmt.Errorf("版本号为空")
    }
    
    return nil
}

// 使用示例
pkg, err := client.GetPackageInfo(ctx, "requests")
if err != nil {
    return err
}

if validationErr := validatePackageData(pkg); validationErr != nil {
    return fmt.Errorf("数据验证失败: %w", validationErr)
}
```

## 上下文错误

### 超时处理

```go
func handleContextError(err error) {
    if err == context.DeadlineExceeded {
        fmt.Println("操作超时，建议：")
        fmt.Println("1. 增加超时时间")
        fmt.Println("2. 检查网络连接")
        fmt.Println("3. 分批处理大量数据")
        return
    }
    
    if err == context.Canceled {
        fmt.Println("操作被取消")
        return
    }
}
```

---

**下一步**: 查看 [示例代码](./examples.md) 获取更多实用的代码示例。
