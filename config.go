package junyousdk

// Config SDK 配置结构
type Config struct {
	// AccessId 访问 ID
	AccessId string
	// AccessKey 访问密钥（Base64 编码）
	AccessKey string
	// Version API 版本（可选，默认 v1）
	Version string
	// Address API 服务器地址（可选，默认 https://open-api.junyouchain.com）
	Address string
	// ContentType 请求内容类型（可选，默认 application/json）
	ContentType string
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Version:     DefaultVersion,
		Address:     DefaultAddress,
		ContentType: DefaultContentType,
	}
}

// WithAccessId 设置 AccessId
func (c *Config) WithAccessId(accessId string) *Config {
	c.AccessId = accessId
	return c
}

// WithAccessKey 设置 AccessKey
func (c *Config) WithAccessKey(accessKey string) *Config {
	c.AccessKey = accessKey
	return c
}

// WithVersion 设置版本
func (c *Config) WithVersion(version string) *Config {
	c.Version = version
	return c
}

// WithAddress 设置服务器地址
func (c *Config) WithAddress(address string) *Config {
	c.Address = address
	return c
}

// WithContentType 设置内容类型
func (c *Config) WithContentType(contentType string) *Config {
	c.ContentType = contentType
	return c
}
