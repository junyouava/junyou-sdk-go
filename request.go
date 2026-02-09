package junyousdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"reflect"
)

// wrappedResponse API 响应结构（带 result 包装）
type wrappedResponse[T any] struct {
	Result *Result[T] `json:"result"`
}

// parseErrorResponse 解析错误响应
func parseErrorResponse[T any](data []byte, statusCode int, apiPath string) (*Result[T], error) {
	// 尝试解析响应
	apiResponse, err := parseResponse[T](data)
	if err == nil {
		return buildErrorResult(apiResponse, statusCode, apiPath, "http")
	}

	// 无法解析 JSON，返回原始响应
	message := string(data)
	if message == "" {
		message = fmt.Sprintf("HTTP %d", statusCode)
	}
	result := NewSysErrorResult[T](message)
	result.Code = statusCode
	return result, fmt.Errorf("http status %d on %s: %s", statusCode, apiPath, message)
}

// isZeroValue 检查值是否为零值
func isZeroValue[T any](v T) bool {
	var zero T
	return reflect.DeepEqual(v, zero)
}

// parseResponse 解析响应（支持带 result 包装和不带包装两种格式）
func parseResponse[T any](data []byte) (*Result[T], error) {
	// 先尝试解析带 result 包装的响应
	var wrappedResp wrappedResponse[T]
	if err := json.Unmarshal(data, &wrappedResp); err == nil && wrappedResp.Result != nil {
		return wrappedResp.Result, nil
	}

	// 尝试解析不带包装的响应
	var directResp Result[T]
	if err := json.Unmarshal(data, &directResp); err != nil {
		return nil, fmt.Errorf("failed to parse response JSON: %w", err)
	}
	return &directResp, nil
}

// buildErrorResult 构建错误结果和错误消息
func buildErrorResult[T any](apiResponse *Result[T], statusCode int, apiPath string, errorType string) (*Result[T], error) {
	result := NewParamErrorResult[T](apiResponse.Message)
	if errorType == "http" {
		result = NewSysErrorResult[T](apiResponse.Message)
	}
	result.Code = statusCode
	result.ErrCode = apiResponse.ErrCode
	result.Data = apiResponse.Data

	// 构建错误消息
	errorMsg := fmt.Sprintf("%s error %d on %s: %s", errorType, statusCode, apiPath, apiResponse.Message)
	if !isZeroValue(apiResponse.Data) {
		dataBytes, _ := json.Marshal(apiResponse.Data)
		errorMsg = fmt.Sprintf("%s (data: %s)", errorMsg, string(dataBytes))
	}

	return result, errors.New(errorMsg)
}

// DoRequest 执行请求
func DoRequest[T any](c *Client, method, apiPath string, body any) (*Result[T], error) {
	// 生成认证 Header
	header, err := c.auth.GenerateAuthHeader(method, apiPath)
	if err != nil {
		return NewSysErrorResult[T]("failed to generate auth header"), fmt.Errorf("failed to generate auth header: %w", err)
	}

	// 构建请求 URL
	baseURL, err := url.Parse(c.config.Address)
	if err != nil {
		return NewSysErrorResult[T]("invalid base URL"), fmt.Errorf("invalid base URL: %w", err)
	}
	// 使用 ResolveReference 或手动拼接路径，避免 path.Join 的问题
	apiURL, err := url.Parse(apiPath)
	if err != nil {
		return NewSysErrorResult[T]("invalid API path"), fmt.Errorf("invalid API path: %w", err)
	}
	reqURL := baseURL.ResolveReference(apiURL).String()

	// 序列化请求体
	var bodyBytes []byte
	if body != nil {
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return NewSysErrorResult[T]("failed to marshal request body"), fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	// 创建 HTTP 请求
	httpReq, err := http.NewRequest(method, reqURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return NewSysErrorResult[T]("failed to create request"), fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// 设置 Header
	httpReq.Header = header.Clone()

	// 请求日志
	reqBodyLog := string(bodyBytes)
	if reqBodyLog == "" {
		reqBodyLog = "(empty)"
	}
	log.Printf("[SDK Request] %s %s", method, reqURL)
	log.Printf("[SDK Request] Headers: %v", httpReq.Header)
	log.Printf("[SDK Request] Body: %s", reqBodyLog)

	// 发送请求
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return NewSysErrorResult[T]("request failed"), fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return NewSysErrorResult[T]("failed to read response"), fmt.Errorf("failed to read response body: %w", err)
	}

	// 响应日志
	respBodyLog := string(data)
	if respBodyLog == "" {
		respBodyLog = "(empty)"
	}
	log.Printf("[SDK Response] %s %s | Status: %d", method, reqURL, resp.StatusCode)
	log.Printf("[SDK Response] Body: %s", respBodyLog)

	// 检查 HTTP 状态码（在解析 JSON 之前）
	if resp.StatusCode != http.StatusOK {
		return parseErrorResponse[T](data, resp.StatusCode, apiPath)
	}

	// 检查响应体是否为空
	if len(data) == 0 {
		var zeroValue T
		return NewSuccessResult("success", zeroValue), nil
	}

	// 解析响应
	apiResponse, err := parseResponse[T](data)
	if err != nil {
		return NewSysErrorResult[T]("failed to parse response"), fmt.Errorf("failed to parse response JSON from %s: %w", apiPath, err)
	}

	// 检查业务状态码
	if apiResponse.Code != http.StatusOK {
		return buildErrorResult(apiResponse, apiResponse.Code, apiPath, "business")
	}

	// 返回成功结果
	result := NewSuccessResult("success", apiResponse.Data)
	result.ErrCode = apiResponse.ErrCode
	return result, nil
}
