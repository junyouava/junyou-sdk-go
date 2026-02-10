package junyousdk

import (
	"net/http"
)

// APIService API 服务
type APIService struct {
	client *Client
}

// NewAPIService 创建 API 服务
func NewAPIService(client *Client) *APIService {
	return &APIService{
		client: client,
	}
}

// RegisterInfo 注册信息
type RegisterInfo struct {
	// PhoneNumber 手机号码
	PhoneNumber string `json:"phone_number"`
}

// OpenIdToken OpenId Token，用于认证相关接口
type OpenIdToken struct {
	// OpenId 用户的 OpenId
	OpenId string `json:"open_id"`
}

// EnterpriseJKSURLRequest 企业 JKS 地址请求
type EnterpriseJKSURLRequest struct {
	// JKSUrl JKS 文件访问地址
	JKSUrl string `json:"jks_url"`
}

// Register 注册
func (s *APIService) Register(registerInfo *RegisterInfo) (*Result[string], error) {
	return DoRequest[string](s.client,
		http.MethodPost,
		APIPathRegister,
		registerInfo,
		nil,
	)
}

// AuthLogin 登录认证
func (s *APIService) AuthLogin(openIdToken OpenIdToken) (*Result[string], error) {
	return DoRequest[string](s.client,
		http.MethodPost,
		APIPathAuthLogin,
		openIdToken,
		nil,
	)
}

// AuthSetPWD 设置密码认证
func (s *APIService) AuthSetPWD(openIdToken OpenIdToken) (*Result[string], error) {
	return DoRequest[string](s.client,
		http.MethodPost,
		APIPathAuthSetPWD,
		openIdToken,
		nil,
	)
}

// AuthCMT 验证认证
func (s *APIService) AuthCMT(openIdToken OpenIdToken) (*Result[string], error) {
	return DoRequest[string](s.client,
		http.MethodPost,
		APIPathAuthCMT,
		openIdToken,
		nil,
	)
}

// SetEnterpriseJKSURL 设置企业 JKS 地址
func (s *APIService) SetEnterpriseJKSURL(req EnterpriseJKSURLRequest) (*Result[map[string]any], error) {
	return DoRequest[map[string]any](s.client,
		http.MethodPost,
		APIPathEnterpriseJKSURL,
		req,
		nil,
	)
}
