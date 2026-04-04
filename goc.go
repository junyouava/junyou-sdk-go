package junyousdk

import (
	"net/http"
)

// PreGOCRewardRequest 对应 POST /api/open/v1/goc/pre_reward。
// Body 仅 amount；收款方由 X-Open-Auth 解析。成功时 result.data 为待上链的 GOC 转账消息 JSON（通常含 from、to、amount、biz_no、biz_type、biz_desc 等），键名以响应为准。
type PreGOCRewardRequest struct {
	Amount string `json:"amount"` // 金额，>0 的十进制字符串
}

// CommitGOCRewardRequest 对应 POST /api/open/v1/goc/reward。
// Body 为 biz_no、message、public_key、der_hex。本请求不加 X-Open-Auth；不要求 pay_password。
type CommitGOCRewardRequest struct {
	BizNo     string `json:"biz_no"`     // 预提交返回的 biz_no
	Message   string `json:"message"`    // 预提交 data 的同内容 JSON 字符串；本地/密盾对该字符串签名，原样提交
	PublicKey string `json:"public_key"` // 链上验签公钥
	DerHex    string `json:"der_hex"`    // 签名 DER 十六进制
}

// PreRewardGOC GOC 预提交。成功时 Data 即待签名/待提交的链上业务消息。
// openAuth 必填：收款方 Open Token（X-Open-Auth），须先对该用户 open_id 调用 AuthLogin；服务端要求每次预提交使用新的 Token。
func (s *APIService) PreRewardGOC(req PreGOCRewardRequest, openAuth string) (*Result[map[string]any], error) {
	return DoRequest[map[string]any](s.client,
		http.MethodPost,
		APIPathGOCPreReward,
		req,
		openAuthExtraHeaders(openAuth),
	)
}

// RewardGOC GOC 提交上链（与 PreRewardGOC 对应）。不携带 X-Open-Auth。
func (s *APIService) RewardGOC(req CommitGOCRewardRequest) (*Result[map[string]any], error) {
	return DoRequest[map[string]any](s.client,
		http.MethodPost,
		APIPathGOCReward,
		req,
		nil,
	)
}
