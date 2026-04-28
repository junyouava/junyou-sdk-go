package junyousdk

import (
	"fmt"
	"net/http"
)

// RefundAsset 开放平台退款业务类型（GOC / EWT 路径与链上业务类型不同，请求体预申请与上链提交一致）。
type RefundAsset string

const (
	// RefundAssetGOC 通用积分退款：POST .../goc/pre_refund、.../goc/refund；须企业 GOC ESACL。
	RefundAssetGOC RefundAsset = "goc"
	// RefundAssetEWT 权证退款：POST .../ewt/pre_refund、.../ewt/refund；须企业权证 ESACL。
	RefundAssetEWT RefundAsset = "ewt"
)

// OpenPreRefundRequest 对应 POST /goc/pre_refund、/ewt/pre_refund 请求体。
// 出账为企业（X-Access-ID），收款人为 X-Open-Auth 解析用户。
type OpenPreRefundRequest struct {
	Amount         string `json:"amount"`                    // 退款金额，十进制字符串，须大于 0
	BizNo          string `json:"biz_no,omitempty"`          // 业务单号，可省略由服务端生成
	BizDescription string `json:"biz_description,omitempty"` // 业务描述，可省略
}

// RefundCommitRequest 与 CommitGOCRewardRequest 及开放文档 GOCRewardCommitRequest 一致，
// GOC/EWT 退款上链共用（biz_no、message、public_key、der_hex）。
type RefundCommitRequest = CommitGOCRewardRequest

func (a RefundAsset) preRefundPath() (string, error) {
	switch a {
	case RefundAssetGOC:
		return APIPathGOCPreRefund, nil
	case RefundAssetEWT:
		return APIPathEWTPreRefund, nil
	default:
		return "", fmt.Errorf("junyousdk: unknown RefundAsset %q", string(a))
	}
}

func (a RefundAsset) refundCommitPath() (string, error) {
	switch a {
	case RefundAssetGOC:
		return APIPathGOCRefund, nil
	case RefundAssetEWT:
		return APIPathEWTRefund, nil
	default:
		return "", fmt.Errorf("junyousdk: unknown RefundAsset %q", string(a))
	}
}

// PreRefund 企业退款预申请。openAuth 为收款用户 Open Token（X-Open-Auth），须与当前企业绑定用户一致。
// 成功时 result.data 一般为 SDKResult（内层待签名消息等），与 GOC pre_reward / pre_refund 文档一致。
func (s *APIService) PreRefund(asset RefundAsset, req OpenPreRefundRequest, openAuth string) (*Result[map[string]any], error) {
	path, err := asset.preRefundPath()
	if err != nil {
		return nil, err
	}
	return DoRequest[map[string]any](s.client,
		http.MethodPost,
		path,
		req,
		openAuthExtraHeaders(openAuth),
	)
}

// RefundCommit 退款上链提交（对 PreRefund 返回的 message 本地签名后提交）。无需 X-Open-Auth。
func (s *APIService) RefundCommit(asset RefundAsset, req RefundCommitRequest) (*Result[map[string]any], error) {
	path, err := asset.refundCommitPath()
	if err != nil {
		return nil, err
	}
	return DoRequest[map[string]any](s.client,
		http.MethodPost,
		path,
		req,
		nil,
	)
}
