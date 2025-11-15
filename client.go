package junyousdk

import (
	"fmt"
	"net/http"
	"time"
)

// Client SDK 客户端
type Client struct {
	config     *Config
	httpClient *http.Client
	auth       *AuthService
	api        *APIService
}

// applyDefaultConfig 应用默认配置值
func applyDefaultConfig(config *Config) {
	if config.Address == "" {
		config.Address = DefaultAddress
	}
	if config.Version == "" {
		config.Version = DefaultVersion
	}
	if config.ContentType == "" {
		config.ContentType = DefaultContentType
	}
}

// validateConfig 验证配置
func validateConfig(config *Config) error {
	if config.AccessId == "" {
		return fmt.Errorf("access_id is required")
	}
	if config.AccessKey == "" {
		return fmt.Errorf("access_key is required")
	}
	return nil
}

// defaultHTTPClient 返回默认的 HTTP 客户端
func defaultHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
	}
}

// NewClient 创建新的 SDK 客户端
func NewClient(config *Config) (*Client, error) {
	if config == nil {
		config = DefaultConfig()
	}
	applyDefaultConfig(config)

	if err := validateConfig(config); err != nil {
		return nil, err
	}

	client := &Client{
		config:     config,
		httpClient: defaultHTTPClient(),
	}

	// 初始化服务
	client.auth = NewAuthService(client)
	client.api = NewAPIService(client)

	return client, nil
}

// NewClientWithHTTPClient 使用自定义 HTTP 客户端创建 SDK 客户端
func NewClientWithHTTPClient(config *Config, httpClient *http.Client) (*Client, error) {
	if config == nil {
		config = DefaultConfig()
	}
	applyDefaultConfig(config)

	if err := validateConfig(config); err != nil {
		return nil, err
	}

	// 如果 httpClient 为 nil，使用默认的
	if httpClient == nil {
		httpClient = defaultHTTPClient()
	}

	client := &Client{
		config:     config,
		httpClient: httpClient,
	}

	// 初始化服务
	client.auth = NewAuthService(client)
	client.api = NewAPIService(client)

	return client, nil
}

// GetConfig 获取配置
func (c *Client) GetConfig() *Config {
	return c.config
}

// GetHTTPClient 获取 HTTP 客户端
func (c *Client) GetHTTPClient() *http.Client {
	return c.httpClient
}

// Auth 返回认证服务
func (c *Client) Auth() *AuthService {
	return c.auth
}

// API 返回 API 服务
func (c *Client) API() *APIService {
	return c.api
}
