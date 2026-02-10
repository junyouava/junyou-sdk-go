package junyousdk

// API 路径常量
const (
	// 注册相关
	APIPathRegister = "/api/open/v1/register"

	// 认证相关
	APIPathAuthLogin  = "/api/open/v1/auth/login"
	APIPathAuthSetPWD = "/api/open/v1/auth/set_pwd"
	APIPathAuthCMT    = "/api/open/v1/auth/cmt"

	// EWT 相关
	APIPathEWTConfirmReleaseByPartner = "/api/open/v1/ewt/confirm_ewt_rbp"
	APIPathEWTCommitReleaseByPartner  = "/api/open/v1/ewt/commit_ewt_rbp"
	APIPathEWTPreOpenReleaseByPartner = "/api/open/v1/ewt/pre_ewt_rbp_open"
	APIPathEWTBalance                 = "/api/open/v1/ewt/balance"
	APIPathEWTTransactionDetails      = "/api/open/v1/ewt/transaction_details"

	// 企业相关
	APIPathEnterpriseJKSURL = "/api/open/v1/enterprise/jks_url"
)

// 默认配置常量
const (
	DefaultAddress     = "https://open-api.junyouchain.com"
	DefaultVersion     = "v1"
	DefaultContentType = "application/json"
)

// 认证 Header 常量
const (
	HeaderAccessId    = "X-Access-ID"
	HeaderSignature   = "X-Signature"
	HeaderNonce       = "X-Signature-Nonce"
	HeaderTimestamp   = "X-Timestamp"
	HeaderContentType = "Content-Type"
	// HeaderOpenAuth 用户 Open Token，用于标识“当前用户”（如预提交权证释放的接收方）。未携带时服务端 userId 为 0，可能返回校验失败。
	HeaderOpenAuth = "X-Open-Auth"
)
