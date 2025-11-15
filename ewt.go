package junyousdk

import (
	"net/http"
)

// EWTBizNoInfo EWT 业务编号信息
type EWTBizNoInfo struct {
	// EWTBizNo EWT 业务编号
	EWTBizNo string `json:"ewt_biz_no"`
}

// ConfirmEWTReleaseByPartner 确认权证释放（合作伙伴）
func (s *APIService) ConfirmEWTReleaseByPartner(ewtBizNoInfo EWTBizNoInfo) (*Result[string], error) {
	return DoRequest[string](s.client,
		http.MethodPost,
		APIPathEWTConfirmReleaseByPartner,
		ewtBizNoInfo,
	)
}
