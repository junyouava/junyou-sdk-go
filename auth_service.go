package junyousdk

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/junyouava/junyou-sdk-go/internal"
)

// Signature 签名信息
type Signature struct {
	AccessId  string `json:"access_id"`
	Signature string `json:"signature"`
	Nonce     string `json:"nonce"`
	Timestamp string `json:"timestamp"`
}

// SignatureWithData 签名信息与数据合并结构
type SignatureWithData struct {
	AccessId  string `json:"access_id"`
	Signature string `json:"signature"`
	Nonce     string `json:"nonce"`
	Timestamp string `json:"timestamp"`
	OpenAuth  string `json:"open_auth"`
}

// AuthService 认证服务
type AuthService struct {
	client *Client
}

// NewAuthService 创建认证服务
func NewAuthService(client *Client) *AuthService {
	return &AuthService{
		client: client,
	}
}

// GenerateSignature 生成签名
func (a *AuthService) GenerateSignature(method, apiPath string) (*Signature, error) {
	// 参数验证
	if method == "" {
		return nil, fmt.Errorf("method is required")
	}
	if apiPath == "" {
		return nil, fmt.Errorf("path is required")
	}

	config := a.client.GetConfig()
	if config.AccessId == "" || config.AccessKey == "" {
		return nil, fmt.Errorf("access_id and access_key are required")
	}

	// 生成 nonce
	nonce, err := internal.GenerateNonce(4)
	if err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// 生成时间戳（当前时间加 3 分钟）
	timestamp := strconv.FormatInt(time.Now().Add(3*time.Minute).Unix(), 10)

	// 构建签名字符串（包含 accessId 作为第一个字段）
	signString := fmt.Sprintf("%s\n%s\n%s\n%s\n%s", config.AccessId, method, apiPath, nonce, timestamp)

	// 解码 AccessKey (Base64)
	accessKeyBytes, err := base64.StdEncoding.DecodeString(config.AccessKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode access_key: %w", err)
	}

	// 计算 HMAC-SHA256 签名
	signatureBytes := internal.HMACSHA256(accessKeyBytes, []byte(signString))
	signature := base64.StdEncoding.EncodeToString(signatureBytes)

	return &Signature{
		AccessId:  config.AccessId,
		Signature: signature,
		Nonce:     nonce,
		Timestamp: timestamp,
	}, nil
}

// GenerateAuthHeader 生成认证 Header
func (a *AuthService) GenerateAuthHeader(method, apiPath string) (http.Header, error) {
	signature, err := a.GenerateSignature(method, apiPath)
	if err != nil {
		return nil, fmt.Errorf("failed to generate signature: %w", err)
	}

	header := make(http.Header)
	header[HeaderAccessId] = []string{signature.AccessId}
	header[HeaderSignature] = []string{signature.Signature}
	header[HeaderNonce] = []string{signature.Nonce}
	header[HeaderTimestamp] = []string{signature.Timestamp}

	config := a.client.GetConfig()
	if config.ContentType != "" {
		header[HeaderContentType] = []string{config.ContentType}
	}

	return header, nil
}

// GenerateSignatureWithAuthCMT 生成签名并调用 AuthCMT，合并返回签名信息和OpenAuth
func (a *AuthService) GenerateSignatureWithAuthCMT(openIdToken OpenIdToken) (*SignatureWithData, error) {
	// 生成签名
	signature, err := a.GenerateSignature(http.MethodPost, APIPathAuthCMT)
	if err != nil {
		return nil, fmt.Errorf("failed to generate signature: %w", err)
	}

	// 调用 AuthCMT
	result, err := a.client.API().AuthCMT(openIdToken)
	if err != nil {
		return nil, fmt.Errorf("failed to call AuthCMT: %w", err)
	}

	// 合并签名信息和数据
	return &SignatureWithData{
		AccessId:  signature.AccessId,
		Signature: signature.Signature,
		Nonce:     signature.Nonce,
		Timestamp: signature.Timestamp,
		OpenAuth:  result.Data,
	}, nil
}
