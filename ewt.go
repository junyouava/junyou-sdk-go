package junyousdk

import (
	"fmt"
	"net/http"
	"net/url"
)

// EWTBizNoInfo EWT 业务编号信息
type EWTBizNoInfo struct {
	// EWTBizNo EWT 业务编号
	EWTBizNo string `json:"ewt_biz_no"`
}

// PreEWTReleaseByPartnerRequest 预提交权证释放请求体
// 对应接口: POST /v1/ewt/pre_ewt_rbp_commit
type PreEWTReleaseByPartnerRequest struct {
	Amount       string `json:"amount"`         // 权证数量
	Ratio        string `json:"ratio"`          // 总释放比例
	Level1OpenId string `json:"level1_open_id"` // 一级合伙人 OpenId
	Level1Ratio  string `json:"level1_ratio"`   // 一级合伙人分配比例
	Level2OpenId string `json:"level2_open_id"` // 二级合伙人 OpenId
	Level2Ratio  string `json:"level2_ratio"`   // 二级合伙人分配比例
}

// CommitEWTReleaseByPartnerRequest 提交权证释放（伙伴）请求体
// 对应接口: POST /v1/ewt/commit_ewt_rbp
type CommitEWTReleaseByPartnerRequest struct {
	BizNo     string `json:"biz_no"`     // 业务单号，例如 EWT20250101000001
	Message   string `json:"message"`    // 原始业务消息 JSON 字符串
	PublicKey string `json:"public_key"` // 公钥（未压缩十六进制）
	DerHex    string `json:"der_hex"`    // DER 编码的签名十六进制
}

// ConfirmEWTReleaseByPartner 确认权证释放（合作伙伴）
func (s *APIService) ConfirmEWTReleaseByPartner(ewtBizNoInfo EWTBizNoInfo) (*Result[string], error) {
	return DoRequest[string](s.client,
		http.MethodPost,
		APIPathEWTConfirmReleaseByPartner,
		ewtBizNoInfo,
	)
}

// PreCommitEWTReleaseByPartner 预提交权证释放（与 CommitEWTReleaseByPartner 配套）
// 对应接口: POST /v1/ewt/pre_ewt_rbp_commit
func (s *APIService) PreCommitEWTReleaseByPartner(req PreEWTReleaseByPartnerRequest) (*Result[map[string]any], error) {
	return DoRequest[map[string]any](s.client,
		http.MethodPost,
		APIPathEWTPreOpenReleaseByPartner,
		req,
	)
}

// CommitEWTReleaseByPartner 提交权证释放（伙伴）
// 对应接口: POST /v1/ewt/commit_ewt_rbp
func (s *APIService) CommitEWTReleaseByPartner(req CommitEWTReleaseByPartnerRequest) (*Result[map[string]any], error) {
	return DoRequest[map[string]any](s.client,
		http.MethodPost,
		APIPathEWTCommitReleaseByPartner,
		req,
	)
}

// GetEWTBalance 权证余额查询
// 对应接口: GET /v1/ewt/balance?page&page_size
// 目前默认按企业维度查询（不自动带 X-Open-Auth）。
func (s *APIService) GetEWTBalance(page, pageSize int) (*Result[map[string]any], error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	query := url.Values{}
	query.Set("page", fmt.Sprintf("%d", page))
	query.Set("page_size", fmt.Sprintf("%d", pageSize))

	apiPath := fmt.Sprintf("%s?%s", APIPathEWTBalance, query.Encode())
	return DoRequest[map[string]any](s.client,
		http.MethodGet,
		apiPath,
		nil,
	)
}

// GetEWTTransactionDetails 权证交易明细查询
// 对应接口: GET /v1/ewt/transaction_details?page&page_size&transaction_type&biz_type&year&month
// 目前默认按企业维度查询（不自动带 X-Open-Auth）。
func (s *APIService) GetEWTTransactionDetails(
	page, pageSize int,
	transactionType, bizType string,
	year, month int,
) (*Result[map[string]any], error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	query := url.Values{}
	query.Set("page", fmt.Sprintf("%d", page))
	query.Set("page_size", fmt.Sprintf("%d", pageSize))
	if transactionType != "" {
		query.Set("transaction_type", transactionType)
	}
	if bizType != "" {
		query.Set("biz_type", bizType)
	}
	if year > 0 {
		query.Set("year", fmt.Sprintf("%d", year))
	}
	if month > 0 {
		query.Set("month", fmt.Sprintf("%d", month))
	}

	apiPath := fmt.Sprintf("%s?%s", APIPathEWTTransactionDetails, query.Encode())
	return DoRequest[map[string]any](s.client,
		http.MethodGet,
		apiPath,
		nil,
	)
}
