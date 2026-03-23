package junyousdk

import (
	"net/http"
)

// PreGOCRewardRequest POST /api/open/v1/goc/pre_reward
type PreGOCRewardRequest struct {
	OpenId string `json:"open_id"` // 收款方 open_id
	Amount string `json:"amount"`  // 金额，>0 的十进制字符串
}

// CommitGOCRewardRequest POST /api/open/v1/goc/reward（无需 pay_password）
type CommitGOCRewardRequest struct {
	BizNo     string `json:"biz_no"`
	Message   string `json:"message"` // 对 PreRewardGOC 返回的 Data 做 json.Marshal 得到的字符串（与 Data 为同一条业务消息）
	PublicKey string `json:"public_key"`
	DerHex    string `json:"der_hex"`
}

// PreRewardGOC GOC 预提交。成功时 Data 即为链上待签/待提交的业务消息（与提交接口 message 字段为同一语义，序列化后即为 message）
func (s *APIService) PreRewardGOC(req PreGOCRewardRequest) (*Result[map[string]any], error) {
	return DoRequest[map[string]any](s.client,
		http.MethodPost,
		APIPathGOCPreReward,
		req,
		nil,
	)
}

// RewardGOC GOC 提交上链（与 PreRewardGOC 对应）
func (s *APIService) RewardGOC(req CommitGOCRewardRequest) (*Result[map[string]any], error) {
	return DoRequest[map[string]any](s.client,
		http.MethodPost,
		APIPathGOCReward,
		req,
		nil,
	)
}
