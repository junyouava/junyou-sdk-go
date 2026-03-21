# Junyou AVA SDK Go

## 功能特性

- 🔐 **安全认证**：支持 HMAC-SHA256 签名算法，自动生成认证 Header（签名仅使用 path，不含 query）
- 📝 **用户注册**：提供用户注册接口，支持手机号注册
- 🔑 **多种认证方式**：支持登录认证、设置密码认证、验证认证等多种令牌获取方式
- 🎫 **权证（EWT）**：确认释放、预提交/提交合伙人释放、余额与交易明细查询；支持可选用户身份（`X-Open-Auth`）
- ⚙️ **灵活配置**：支持自定义配置，包括 API 地址、版本、内容类型等
- 🔧 **自定义 HTTP 客户端**：支持使用自定义 HTTP 客户端，方便集成到现有项目
- 📦 **类型安全**：使用 Go 泛型，提供类型安全的 API 响应处理
- 🛡️ **完善的错误处理**：区分网络错误和业务错误，提供详细的错误信息

## 要求

- **Go 版本**：>= 1.21

本项目使用 Go 1.21 作为最低版本要求，主要使用了以下特性：

- Go 泛型（Generics）
- 标准库的现代化 API

## 安装

```bash
go get github.com/junyouava/junyou-sdk-go
```

## 快速开始

### 初始化客户端

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "time"

    junyousdk "github.com/junyouava/junyou-sdk-go"
)

func main() {
    // 方式1: 使用默认配置
    config := junyousdk.DefaultConfig().
        WithAccessId("your-access-id").
        WithAccessKey("your-access-key")

    client, err := junyousdk.NewClient(config)
    if err != nil {
        log.Fatal(err)
    }

    // 方式2: 直接创建配置
    client, err = junyousdk.NewClient(&junyousdk.Config{
        AccessId:  "your-access-id",
        AccessKey: "your-access-key",
    })
    if err != nil {
        log.Fatal(err)
    }

    // 方式3: 使用自定义 HTTP 客户端
    httpClient := &http.Client{
        Timeout: 30 * time.Second,
    }
    client, err = junyousdk.NewClientWithHTTPClient(config, httpClient)
    if err != nil {
        log.Fatal(err)
    }
}
```

### 注册

```go
registerInfo := &junyousdk.RegisterInfo{
    PhoneNumber: "13800138000",
}

result, err := client.API().Register(registerInfo)
if err != nil {
    fmt.Printf("注册失败: %v\n", err)
    return
}

if !result.Success {
    fmt.Printf("注册失败: %s\n", result.Message)
    return
}

fmt.Printf("注册成功: %s\n", result.Data)
```

### 获取登录令牌（Open Token / openAuth）

调用 `POST /api/open/v1/auth/login`，用用户的 `open_id` 换取 **Open Token**。成功时 **`result.Data` 即为后续 EWT 等接口中的 `openAuth`**：SDK 会将其放在请求头 **`X-Open-Auth`**（常量 `junyousdk.HeaderOpenAuth`）上。

```go
openIdToken := junyousdk.OpenIdToken{
    OpenId: "user-open-id",
}

result, err := client.API().AuthLogin(openIdToken)
if err != nil {
    fmt.Printf("获取令牌失败: %v\n", err)
    return
}

if !result.Success {
    fmt.Printf("获取令牌失败: %s\n", result.Message)
    return
}

openAuth := result.Data // 传给 GetEWTBalance、GetEWTTransactionDetails、PreCommitEWTReleaseByPartner 等
fmt.Printf("Open Token: %s\n", openAuth)
```

### 获取设置密码令牌

```go
openIdToken := junyousdk.OpenIdToken{
    OpenId: "user-open-id",
}

result, err := client.API().AuthSetPWD(openIdToken)
if err != nil {
    fmt.Printf("获取令牌失败: %v\n", err)
    return
}

if !result.Success {
    fmt.Printf("获取令牌失败: %s\n", result.Message)
    return
}

accessToken := result.Data
fmt.Printf("Access Token: %s\n", accessToken)
```

### 获取验证令牌

```go
openIdToken := junyousdk.OpenIdToken{
    OpenId: "user-open-id",
}

result, err := client.API().AuthCMT(openIdToken)
if err != nil {
    fmt.Printf("获取令牌失败: %v\n", err)
    return
}

if !result.Success {
    fmt.Printf("获取令牌失败: %s\n", result.Message)
    return
}

accessToken := result.Data
fmt.Printf("Access Token: %s\n", accessToken)
```

### 权证：确认释放

```go
ewtBizNoInfo := junyousdk.EWTBizNoInfo{
    EWTBizNo: "ewt-biz-no",
}

result, err := client.API().ConfirmEWTReleaseByPartner(ewtBizNoInfo)
if err != nil {
    fmt.Printf("释放权证失败: %v\n", err)
    return
}

if !result.Success {
    fmt.Printf("释放权证失败: %s\n", result.Message)
    return
}

fmt.Println("释放权证成功")
```

### 权证：合伙人释放（预提交 + 提交）

预提交接口需要 **用户身份**，须先 `AuthLogin` 得到 `openAuth`，再传入 `PreCommitEWTReleaseByPartner`。提交阶段使用预提交返回的 `biz_no` 及本地签名结果调用 `CommitEWTReleaseByPartner`。

```go
// 1. 换取 Open Token（接收方用户）
loginResult, err := client.API().AuthLogin(junyousdk.OpenIdToken{OpenId: "receiver-open-id"})
if err != nil || !loginResult.Success {
    // 处理错误
    return
}
openAuth := loginResult.Data

// 2. 预提交
preReq := junyousdk.PreEWTReleaseByPartnerRequest{
    Amount: "100", Ratio: "10",
    // Level1OpenId / Level1Ratio / Level2OpenId / Level2Ratio 按业务填写
}
preResult, err := client.API().PreCommitEWTReleaseByPartner(preReq, openAuth)
// 根据 preResult.Data 解析待签名内容、biz_no，本地签名后：

commitReq := junyousdk.CommitEWTReleaseByPartnerRequest{
    BizNo:     "从预提交返回解析",
    Message:   "与签名一致的 JSON 字符串",
    PublicKey: "公钥",
    DerHex:    "DER 签名十六进制",
}
commitResult, err := client.API().CommitEWTReleaseByPartner(commitReq)
// 处理 commitResult
```

### 权证：余额查询

- **`openAuth` 为空字符串**：不按用户过滤，一般为 **企业维度**。
- **`openAuth` 非空**：带 `X-Open-Auth`，按 **该用户** 查询（通常即 `AuthLogin` 返回的 `result.Data`）。

```go
// 企业维度
result, err := client.API().GetEWTBalance(1, 10, "")

// 用户维度（需先 AuthLogin）
result, err = client.API().GetEWTBalance(1, 10, openAuth)
```

### 权证：交易明细查询

参数与余额类似，最后增加 **`openAuth`**；`transaction_type`、`biz_type` 可传空串，`year`/`month` 为 0 表示不按年月筛选。

```go
result, err := client.API().GetEWTTransactionDetails(
    1, 10,
    "", "", // transaction_type, biz_type
    0, 0,   // year, month
    "",     // openAuth：企业维度
)

result, err = client.API().GetEWTTransactionDetails(
    1, 10, "in", "EWT1005", 2025, 3,
    openAuth, // 用户维度
)
```

### 企业 JKS 访问链接上报

```go
result, err := client.API().SetEnterpriseJKSURL(junyousdk.EnterpriseJKSURLRequest{
    JKSUrl: "https://your-vault.example.com/#code=xxx",
})
```

### 生成签名

```go
signature, err := client.Auth().GenerateSignature("POST", "/api/open/v1/register")
if err != nil {
    fmt.Printf("生成签名失败: %v\n", err)
    return
}

fmt.Printf("AccessId: %s\n", signature.AccessId)
fmt.Printf("Signature: %s\n", signature.Signature)
fmt.Printf("Nonce: %s\n", signature.Nonce)
fmt.Printf("Timestamp: %s\n", signature.Timestamp)
```

### 生成认证 Header

```go
header, err := client.Auth().GenerateAuthHeader("POST", "/api/open/v1/register")
if err != nil {
    fmt.Printf("生成认证 Header 失败: %v\n", err)
    return
}

// header 可以直接用于 HTTP 请求
```

## API 文档

### Client

SDK 主客户端，提供所有服务访问入口。

#### 方法

- `NewClient(config *Config) (*Client, error)` - 创建新客户端（会验证配置）
- `NewClientWithHTTPClient(config *Config, httpClient *http.Client) (*Client, error)` - 使用自定义 HTTP 客户端创建客户端（会验证配置）
- `GetConfig() *Config` - 获取配置
- `GetHTTPClient() *http.Client` - 获取 HTTP 客户端
- `Auth() *AuthService` - 获取认证服务
- `API() *APIService` - 获取 API 服务

### AuthService

认证服务，提供签名和认证 Header 生成功能。

#### 方法

- `GenerateSignature(method, path string) (*Signature, error)` - 生成签名（path 可含 query，参与签名的为 `?` 前的 path）
- `GenerateAuthHeader(method, path string) (http.Header, error)` - 生成认证 Header
- `GenerateSignatureWithOpenAuth(method, path string, openIdToken OpenIdToken) (*SignatureWithOpenAuth, error)` - 生成签名并调用 `AuthCMT`，合并返回 OpenAuth 等信息

### APIService

API 服务，提供所有业务 API 调用。

#### 方法

| 方法 | 说明 |
|------|------|
| `Register(registerInfo *RegisterInfo) (*Result[string], error)` | 用户注册 |
| `AuthLogin(openIdToken OpenIdToken) (*Result[string], error)` | 登录；`Data` 为 Open Token（作 `openAuth`） |
| `AuthSetPWD(openIdToken OpenIdToken) (*Result[string], error)` | 设置密码相关令牌 |
| `AuthCMT(openIdToken OpenIdToken) (*Result[string], error)` | 验证认证令牌 |
| `SetEnterpriseJKSURL(req EnterpriseJKSURLRequest) (*Result[map[string]any], error)` | 设置企业 JKS 访问地址 |
| `ConfirmEWTReleaseByPartner(ewtBizNoInfo EWTBizNoInfo) (*Result[string], error)` | 确认权证释放 |
| `PreCommitEWTReleaseByPartner(req PreEWTReleaseByPartnerRequest, openAuth string) (*Result[map[string]any], error)` | 预提交合伙人释放；`openAuth` 一般来自 `AuthLogin` |
| `CommitEWTReleaseByPartner(req CommitEWTReleaseByPartnerRequest) (*Result[map[string]any], error)` | 提交合伙人释放 |
| `GetEWTBalance(page, pageSize int, openAuth string) (*Result[map[string]any], error)` | 权证余额；`openAuth==""` 企业维度，否则用户维度 |
| `GetEWTTransactionDetails(page, pageSize int, transactionType, bizType string, year, month int, openAuth string) (*Result[map[string]any], error)` | 权证交易明细；`openAuth` 语义同余额 |

## 配置选项

### Config

```go
type Config struct {
    AccessId    string // 访问 ID（必需）
    AccessKey   string // 访问密钥（必需，Base64 编码）
    Version     string // API 版本（可选，默认 "v1"）
    Address     string // API 根地址（可选，默认 "https://open-api.junyouchain.com"）
    ContentType string // 请求内容类型（可选，默认 "application/json"）
}
```

### 配置方法

- `DefaultConfig() *Config` - 创建默认配置
- `WithAccessId(accessId string) *Config` - 设置 AccessId
- `WithAccessKey(accessKey string) *Config` - 设置 AccessKey
- `WithVersion(version string) *Config` - 设置版本
- `WithAddress(address string) *Config` - 设置服务器地址
- `WithContentType(contentType string) *Config` - 设置内容类型

## 错误处理

所有 API 方法都返回 `*Result[T]` 和 `error`。

- `Result.Success` - 表示请求是否成功
- `Result.Code` - HTTP 状态码或业务状态码
- `Result.ErrCode` - 业务错误代码（字符串）
- `Result.Message` - 错误或成功消息
- `Result.Data` - 响应数据

示例：

```go
import (
    "fmt"
    "log"
)

result, err := client.API().Register(registerInfo)
if err != nil {
    // 网络错误或其他系统错误
    log.Fatal(err)
}

if !result.Success {
    // 业务错误
    if result.ErrCode != "" {
        fmt.Printf("错误: %s (错误代码: %s, 状态码: %d)\n", result.Message, result.ErrCode, result.Code)
    } else {
        fmt.Printf("错误: %s (状态码: %d)\n", result.Message, result.Code)
    }
    return
}

// 成功
fmt.Printf("成功: %s\n", result.Data)
```

## 许可证

MIT License
