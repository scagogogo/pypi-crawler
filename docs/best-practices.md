# 最佳实践

本文档提供使用 PyPI Crawler 的最佳实践建议，帮助您构建高效、稳定的应用程序。

## 📋 目录

- [性能优化](#性能优化)
- [错误处理策略](#错误处理策略)
- [并发控制](#并发控制)
- [缓存策略](#缓存策略)
- [监控和日志](#监控和日志)
- [安全考虑](#安全考虑)
- [生产环境部署](#生产环境部署)

## 性能优化

### 1. 选择合适的镜像源

```go
// 根据地理位置选择镜像源
func selectOptimalMirror() api.PyPIClient {
    // 中国大陆用户
    if isInChina() {
        return mirrors.NewTsinghuaClient()
    }
    
    // 其他地区用户
    return mirrors.NewOfficialClient()
}

func isInChina() bool {
    // 简单的地区检测逻辑
    // 实际应用中可以使用更精确的方法
    return strings.Contains(os.Getenv("TZ"), "Asia/Shanghai")
}
```

### 2. 复用客户端实例

```go
// ❌ 错误做法：每次创建新客户端
func badExample() {
    for _, pkg := range packages {
        client := mirrors.NewTsinghuaClient() // 每次都创建新客户端
        info, _ := client.GetPackageInfo(ctx, pkg)
        // ...
    }
}

// ✅ 正确做法：复用客户端
func goodExample() {
    client := mirrors.NewTsinghuaClient() // 创建一次
    
    for _, pkg := range packages {
        info, _ := client.GetPackageInfo(ctx, pkg)
        // ...
    }
}
```

### 3. 合理设置超时时间

```go
func createOptimizedClient() api.PyPIClient {
    options := client.NewOptions()
    
    // 根据操作类型设置不同的超时时间
    switch operationType {
    case "single_query":
        options.WithTimeout(10 * time.Second)
    case "batch_operation":
        options.WithTimeout(5 * time.Minute)
    case "full_index":
        options.WithTimeout(10 * time.Minute)
    }
    
    return mirrors.NewTsinghuaClient(options)
}
```

### 4. 使用连接池

```go
type ClientPool struct {
    clients chan api.PyPIClient
    factory func() api.PyPIClient
}

func NewClientPool(size int, factory func() api.PyPIClient) *ClientPool {
    pool := &ClientPool{
        clients: make(chan api.PyPIClient, size),
        factory: factory,
    }
    
    // 预填充连接池
    for i := 0; i < size; i++ {
        pool.clients <- factory()
    }
    
    return pool
}

func (p *ClientPool) Get() api.PyPIClient {
    select {
    case client := <-p.clients:
        return client
    default:
        return p.factory()
    }
}

func (p *ClientPool) Put(client api.PyPIClient) {
    select {
    case p.clients <- client:
    default:
        // 池已满，丢弃客户端
    }
}
```

## 错误处理策略

### 1. 分层错误处理

```go
type PyPIError struct {
    Type    string
    Message string
    Cause   error
    Retry   bool
}

func (e *PyPIError) Error() string {
    return fmt.Sprintf("[%s] %s", e.Type, e.Message)
}

func classifyError(err error) *PyPIError {
    errStr := err.Error()
    
    switch {
    case strings.Contains(errStr, "404"):
        return &PyPIError{
            Type:    "NOT_FOUND",
            Message: "包不存在",
            Cause:   err,
            Retry:   false,
        }
    case strings.Contains(errStr, "timeout"):
        return &PyPIError{
            Type:    "TIMEOUT",
            Message: "请求超时",
            Cause:   err,
            Retry:   true,
        }
    case strings.Contains(errStr, "5"):
        return &PyPIError{
            Type:    "SERVER_ERROR",
            Message: "服务器错误",
            Cause:   err,
            Retry:   true,
        }
    default:
        return &PyPIError{
            Type:    "UNKNOWN",
            Message: "未知错误",
            Cause:   err,
            Retry:   false,
        }
    }
}
```

### 2. 智能重试机制

```go
type RetryConfig struct {
    MaxAttempts int
    BaseDelay   time.Duration
    MaxDelay    time.Duration
    Multiplier  float64
}

func retryWithBackoff(fn func() error, config RetryConfig) error {
    var lastErr error
    delay := config.BaseDelay
    
    for attempt := 0; attempt < config.MaxAttempts; attempt++ {
        if attempt > 0 {
            time.Sleep(delay)
            
            // 指数退避
            delay = time.Duration(float64(delay) * config.Multiplier)
            if delay > config.MaxDelay {
                delay = config.MaxDelay
            }
        }
        
        if err := fn(); err == nil {
            return nil
        } else {
            lastErr = err
            
            // 检查是否应该重试
            if pypiErr := classifyError(err); !pypiErr.Retry {
                break
            }
        }
    }
    
    return fmt.Errorf("重试 %d 次后失败: %w", config.MaxAttempts, lastErr)
}

// 使用示例
func getPackageWithRetry(client api.PyPIClient, packageName string) (*models.Package, error) {
    var result *models.Package
    
    err := retryWithBackoff(func() error {
        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()
        
        pkg, err := client.GetPackageInfo(ctx, packageName)
        if err != nil {
            return err
        }
        
        result = pkg
        return nil
    }, RetryConfig{
        MaxAttempts: 3,
        BaseDelay:   1 * time.Second,
        MaxDelay:    10 * time.Second,
        Multiplier:  2.0,
    })
    
    return result, err
}
```

## 并发控制

### 1. 限制并发数

```go
type ConcurrencyLimiter struct {
    semaphore chan struct{}
}

func NewConcurrencyLimiter(limit int) *ConcurrencyLimiter {
    return &ConcurrencyLimiter{
        semaphore: make(chan struct{}, limit),
    }
}

func (cl *ConcurrencyLimiter) Acquire() {
    cl.semaphore <- struct{}{}
}

func (cl *ConcurrencyLimiter) Release() {
    <-cl.semaphore
}

func (cl *ConcurrencyLimiter) Do(fn func()) {
    cl.Acquire()
    defer cl.Release()
    fn()
}

// 使用示例
func batchProcessPackages(client api.PyPIClient, packages []string) {
    limiter := NewConcurrencyLimiter(5) // 最多5个并发
    var wg sync.WaitGroup
    
    for _, pkg := range packages {
        wg.Add(1)
        go func(packageName string) {
            defer wg.Done()
            
            limiter.Do(func() {
                processPackage(client, packageName)
            })
        }(pkg)
    }
    
    wg.Wait()
}
```

### 2. 速率限制

```go
type RateLimiter struct {
    ticker   *time.Ticker
    requests chan struct{}
}

func NewRateLimiter(requestsPerSecond int) *RateLimiter {
    rl := &RateLimiter{
        ticker:   time.NewTicker(time.Second / time.Duration(requestsPerSecond)),
        requests: make(chan struct{}, requestsPerSecond),
    }
    
    // 填充初始令牌
    for i := 0; i < requestsPerSecond; i++ {
        rl.requests <- struct{}{}
    }
    
    // 定期补充令牌
    go func() {
        for range rl.ticker.C {
            select {
            case rl.requests <- struct{}{}:
            default:
                // 令牌桶已满
            }
        }
    }()
    
    return rl
}

func (rl *RateLimiter) Wait() {
    <-rl.requests
}

func (rl *RateLimiter) Close() {
    rl.ticker.Stop()
}
```

## 缓存策略

### 1. 内存缓存

```go
type MemoryCache struct {
    data   map[string]*CacheItem
    mutex  sync.RWMutex
    maxAge time.Duration
}

type CacheItem struct {
    Value     interface{}
    ExpiresAt time.Time
}

func NewMemoryCache(maxAge time.Duration) *MemoryCache {
    cache := &MemoryCache{
        data:   make(map[string]*CacheItem),
        maxAge: maxAge,
    }
    
    // 定期清理过期项
    go cache.cleanup()
    
    return cache
}

func (mc *MemoryCache) Get(key string) (interface{}, bool) {
    mc.mutex.RLock()
    defer mc.mutex.RUnlock()
    
    item, exists := mc.data[key]
    if !exists || time.Now().After(item.ExpiresAt) {
        return nil, false
    }
    
    return item.Value, true
}

func (mc *MemoryCache) Set(key string, value interface{}) {
    mc.mutex.Lock()
    defer mc.mutex.Unlock()
    
    mc.data[key] = &CacheItem{
        Value:     value,
        ExpiresAt: time.Now().Add(mc.maxAge),
    }
}

func (mc *MemoryCache) cleanup() {
    ticker := time.NewTicker(time.Minute)
    defer ticker.Stop()
    
    for range ticker.C {
        mc.mutex.Lock()
        now := time.Now()
        
        for key, item := range mc.data {
            if now.After(item.ExpiresAt) {
                delete(mc.data, key)
            }
        }
        
        mc.mutex.Unlock()
    }
}
```

### 2. 缓存装饰器

```go
type CachedClient struct {
    client api.PyPIClient
    cache  *MemoryCache
}

func NewCachedClient(client api.PyPIClient, cacheMaxAge time.Duration) *CachedClient {
    return &CachedClient{
        client: client,
        cache:  NewMemoryCache(cacheMaxAge),
    }
}

func (cc *CachedClient) GetPackageInfo(ctx context.Context, packageName string) (*models.Package, error) {
    cacheKey := fmt.Sprintf("package:%s", packageName)
    
    // 尝试从缓存获取
    if cached, found := cc.cache.Get(cacheKey); found {
        return cached.(*models.Package), nil
    }
    
    // 缓存未命中，从API获取
    pkg, err := cc.client.GetPackageInfo(ctx, packageName)
    if err != nil {
        return nil, err
    }
    
    // 存入缓存
    cc.cache.Set(cacheKey, pkg)
    
    return pkg, nil
}
```

## 监控和日志

### 1. 结构化日志

```go
import "log/slog"

type LoggedClient struct {
    client api.PyPIClient
    logger *slog.Logger
}

func NewLoggedClient(client api.PyPIClient, logger *slog.Logger) *LoggedClient {
    return &LoggedClient{
        client: client,
        logger: logger,
    }
}

func (lc *LoggedClient) GetPackageInfo(ctx context.Context, packageName string) (*models.Package, error) {
    start := time.Now()
    
    lc.logger.Info("开始获取包信息",
        slog.String("package", packageName),
        slog.String("operation", "GetPackageInfo"))
    
    pkg, err := lc.client.GetPackageInfo(ctx, packageName)
    
    duration := time.Since(start)
    
    if err != nil {
        lc.logger.Error("获取包信息失败",
            slog.String("package", packageName),
            slog.String("error", err.Error()),
            slog.Duration("duration", duration))
        return nil, err
    }
    
    lc.logger.Info("成功获取包信息",
        slog.String("package", packageName),
        slog.String("version", pkg.Info.Version),
        slog.Duration("duration", duration))
    
    return pkg, nil
}
```

### 2. 性能指标收集

```go
type Metrics struct {
    RequestCount    int64
    ErrorCount      int64
    TotalDuration   time.Duration
    mutex          sync.RWMutex
}

func (m *Metrics) RecordRequest(duration time.Duration, err error) {
    m.mutex.Lock()
    defer m.mutex.Unlock()
    
    m.RequestCount++
    m.TotalDuration += duration
    
    if err != nil {
        m.ErrorCount++
    }
}

func (m *Metrics) GetStats() (requests int64, errors int64, avgDuration time.Duration) {
    m.mutex.RLock()
    defer m.mutex.RUnlock()
    
    requests = m.RequestCount
    errors = m.ErrorCount
    
    if requests > 0 {
        avgDuration = m.TotalDuration / time.Duration(requests)
    }
    
    return
}

type MetricsClient struct {
    client  api.PyPIClient
    metrics *Metrics
}

func NewMetricsClient(client api.PyPIClient) *MetricsClient {
    return &MetricsClient{
        client:  client,
        metrics: &Metrics{},
    }
}

func (mc *MetricsClient) GetPackageInfo(ctx context.Context, packageName string) (*models.Package, error) {
    start := time.Now()
    
    pkg, err := mc.client.GetPackageInfo(ctx, packageName)
    
    duration := time.Since(start)
    mc.metrics.RecordRequest(duration, err)
    
    return pkg, err
}
```

## 安全考虑

### 1. 输入验证

```go
func validatePackageName(name string) error {
    if name == "" {
        return fmt.Errorf("包名不能为空")
    }
    
    if len(name) > 214 {
        return fmt.Errorf("包名过长")
    }
    
    // PyPI包名规则验证
    matched, _ := regexp.MatchString(`^[a-zA-Z0-9]([a-zA-Z0-9._-]*[a-zA-Z0-9])?$`, name)
    if !matched {
        return fmt.Errorf("包名格式无效")
    }
    
    return nil
}

func safeGetPackageInfo(client api.PyPIClient, ctx context.Context, packageName string) (*models.Package, error) {
    if err := validatePackageName(packageName); err != nil {
        return nil, fmt.Errorf("输入验证失败: %w", err)
    }
    
    return client.GetPackageInfo(ctx, packageName)
}
```

### 2. 敏感信息处理

```go
func sanitizeUserAgent(userAgent string) string {
    // 移除可能的敏感信息
    re := regexp.MustCompile(`\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b`) // IP地址
    userAgent = re.ReplaceAllString(userAgent, "[IP]")
    
    re = regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`) // 邮箱
    userAgent = re.ReplaceAllString(userAgent, "[EMAIL]")
    
    return userAgent
}
```

## 生产环境部署

### 1. 配置管理

```go
type Config struct {
    Mirror      string        `env:"PYPI_MIRROR" default:"tsinghua"`
    Timeout     time.Duration `env:"PYPI_TIMEOUT" default:"30s"`
    MaxRetries  int           `env:"PYPI_MAX_RETRIES" default:"3"`
    Concurrency int           `env:"PYPI_CONCURRENCY" default:"5"`
    CacheMaxAge time.Duration `env:"PYPI_CACHE_MAX_AGE" default:"1h"`
    LogLevel    string        `env:"LOG_LEVEL" default:"info"`
}

func LoadConfig() (*Config, error) {
    var config Config
    
    // 使用环境变量加载配置
    if err := env.Parse(&config); err != nil {
        return nil, err
    }
    
    return &config, nil
}

func CreateClientFromConfig(config *Config) api.PyPIClient {
    options := client.NewOptions().
        WithTimeout(config.Timeout).
        WithMaxRetries(config.MaxRetries)
    
    switch config.Mirror {
    case "tsinghua":
        return mirrors.NewTsinghuaClient(options)
    case "aliyun":
        return mirrors.NewAliyunClient(options)
    default:
        return mirrors.NewOfficialClient(options)
    }
}
```

### 2. 健康检查

```go
func HealthCheck(client api.PyPIClient) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    // 尝试获取一个知名包的信息
    _, err := client.GetPackageInfo(ctx, "requests")
    return err
}

func StartHealthCheckServer(client api.PyPIClient, port int) {
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        if err := HealthCheck(client); err != nil {
            w.WriteHeader(http.StatusServiceUnavailable)
            fmt.Fprintf(w, "Health check failed: %v", err)
            return
        }
        
        w.WriteHeader(http.StatusOK)
        fmt.Fprint(w, "OK")
    })
    
    log.Printf("Health check server starting on port %d", port)
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
```

### 3. 优雅关闭

```go
func GracefulShutdown(cancel context.CancelFunc) {
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
    
    go func() {
        <-sigChan
        log.Println("收到关闭信号，开始优雅关闭...")
        cancel()
    }()
}
```

---

**下一步**: 查看 [常见问题](./faq.md) 了解常见问题的解决方案。
