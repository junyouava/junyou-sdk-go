package main

import (
	"fmt"
	"log"

	junyousdk "github.com/junyouava/junyou-sdk-go"
)

func main() {
	// 初始化客户端
	config := junyousdk.DefaultConfig().
		WithAccessId("084f95e5e2bf3f79d5f2fd069f4f5e7c").
		WithAccessKey("5L/1P8XJ2dIWIMGEHkrZ6gE0HGKvyd/4MKcyQ04oEfE=").
		WithAddress("https://open-api.test.junyouchain.com")

	client, err := junyousdk.NewClient(config)
	if err != nil {
		log.Fatalf("创建客户端失败: %v\n", err)
	}

	// 示例1: 注册
	// registerExample(client)

	// 示例2: 获取登录令牌
	// loginExample(client)

	// 示例3: 释放权证
	// ewtExample(client)

	// 示例5: 预提交权证释放（伙伴）
	// 对应接口: POST /api/open/v1/ewt/pre_ewt_rbp_commit
	ewtPreReleaseExample(client)

	// 示例6: 提交权证释放（伙伴）
	// 对应接口: POST /api/open/v1/ewt/commit_ewt_rbp
	// ewtCommitReleaseExample(client)

	// 示例7: 权证余额查询
	// 对应接口: GET /api/open/v1/ewt/balance
	// ewtBalanceExample(client)

	// 示例8: 权证交易明细查询
	// 对应接口: GET /api/open/v1/ewt/transaction_details
	// ewtTransactionDetailsExample(client)
}

func registerExample(client *junyousdk.Client) {
	fmt.Println("=== 注册示例 ===")

	registerInfo := &junyousdk.RegisterInfo{
		PhoneNumber: "13800138000",
	}

	result, err := client.API().Register(registerInfo)
	if err != nil {
		log.Printf("注册失败: %v\n", err)
		return
	}

	if !result.Success {
		log.Printf("注册失败: %s\n", result.Message)
		return
	}

	fmt.Printf("注册成功: %s\n", result.Data)
}

func loginExample(client *junyousdk.Client) {
	fmt.Println("\n=== 登录示例 ===")

	openIdToken := junyousdk.OpenIdToken{
		OpenId: "user-open-id",
	}

	result, err := client.API().AuthLogin(openIdToken)
	if err != nil {
		log.Printf("获取令牌失败: %v\n", err)
		return
	}

	if !result.Success {
		log.Printf("获取令牌失败: %s\n", result.Message)
		return
	}

	fmt.Printf("Access Token: %s\n", result.Data)
}

func ewtExample(client *junyousdk.Client) {
	fmt.Println("\n=== 释放权证示例 ===")

	ewtBizNoInfo := junyousdk.EWTBizNoInfo{
		EWTBizNo: "ewt-biz-no",
	}

	result, err := client.API().ConfirmEWTReleaseByPartner(ewtBizNoInfo)
	if err != nil {
		log.Printf("释放权证失败: %v\n", err)
		return
	}

	if !result.Success {
		log.Printf("释放权证失败: %s\n", result.Message)
		return
	}

	fmt.Println("释放权证成功")
}

// ewtPreReleaseExample 预提交权证释放示例
// 对应接口: POST /api/open/v1/ewt/pre_ewt_rbp_open
// 预提交需要“用户身份”：先为该用户换取 Open Token，再在请求头中携带 X-Open-Auth，否则服务端会返回「校验失败：缺少用户身份」。
func ewtPreReleaseExample(client *junyousdk.Client) {
	fmt.Println("\n=== 预提交权证释放示例 ===")

	// 1. 为“接收权证释放”的用户换取 Open Token（示例用 openId，实际由业务侧提供）
	openId := "04a7bb30587780d34fd7916664b13651ee4a05dc8079c34a69e9cea2cc59faf7"
	loginResult, err := client.API().AuthLogin(junyousdk.OpenIdToken{OpenId: openId})
	if err != nil {
		log.Printf("获取 Open Token 失败: %v\n", err)
		return
	}
	if !loginResult.Success {
		log.Printf("获取 Open Token 失败: %s\n", loginResult.Message)
		return
	}
	openAuth := loginResult.Data

	// 2. 带 X-Open-Auth 调用预提交
	preReq := junyousdk.PreEWTReleaseByPartnerRequest{
		Amount:       "100",
		Ratio:        "1",
		Level1OpenId: "04a7bb30587780d34fd7916664b13651ee4a05dc8079c34a69e9cea2cc59faf7",
		Level1Ratio:  "0.7",
		Level2OpenId: "d92067abdbb2e2b68a4ad31597e45c1944389c0b26324233a4498a9066037369",
		Level2Ratio:  "0.3",
	}

	preResult, err := client.API().PreCommitEWTReleaseByPartner(preReq, openAuth)
	if err != nil {
		log.Printf("预提交权证释放失败: %v\n", err)
		return
	}
	if !preResult.Success {
		log.Printf("预提交权证释放失败: %s\n", preResult.Message)
		return
	}

	fmt.Printf("预提交成功，返回数据: %#v\n", preResult.Data)
}

// ewtCommitReleaseExample 提交权证释放示例
// 对应接口: POST /api/open/v1/ewt/commit_ewt_rbp
func ewtCommitReleaseExample(client *junyousdk.Client) {
	fmt.Println("\n=== 提交权证释放示例 ===")

	// TODO: 这里的 biz_no / message / 公钥 / 签名 等字段，请根据你业务系统真实数据填充
	bizNo := "EWT20250101000001"
	message := `{"from":"0x1234567890abcdef1234567890abcdef12345678","to":"0xabcdef1234567890abcdef1234567890abcdef12","amount":"100","ratio":"0.5","biz_no":"EWT20250101000001","ent_id":123,"biz_type":"EWT1005","biz_desc":"合伙人权证释放"}`

	commitReq := junyousdk.CommitEWTReleaseByPartnerRequest{
		BizNo:     bizNo,
		Message:   message,
		PublicKey: "04a1b2c3d4e5f678901234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		DerHex:    "3045022100a1b2c3d4e5f678901234567890abcdef1234567890abcdef1234567890abcdef1234567890022001234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
	}

	commitResult, err := client.API().CommitEWTReleaseByPartner(commitReq)
	if err != nil {
		log.Printf("提交权证释放失败: %v\n", err)
		return
	}
	if !commitResult.Success {
		log.Printf("提交权证释放失败: %s\n", commitResult.Message)
		return
	}

	fmt.Printf("提交权证释放成功，返回数据: %#v\n", commitResult.Data)
}

// ewtBalanceExample 权证余额查询示例
func ewtBalanceExample(client *junyousdk.Client) {
	fmt.Println("\n=== 权证余额查询示例 ===")

	result, err := client.API().GetEWTBalance(1, 20)
	if err != nil {
		log.Printf("查询权证余额失败: %v\n", err)
		return
	}
	if !result.Success {
		log.Printf("查询权证余额失败: %s\n", result.Message)
		return
	}

	fmt.Printf("权证余额查询结果: %#v\n", result.Data)
}

// ewtTransactionDetailsExample 权证交易明细查询示例
func ewtTransactionDetailsExample(client *junyousdk.Client) {
	fmt.Println("\n=== 权证交易明细查询示例 ===")

	// 这里使用演示参数，你可以根据业务需要调整
	page := 1
	pageSize := 20
	transactionType := "" // 比如: "in" / "out"，根据后端约定填写
	bizType := ""         // 比如: "EWT1005"
	year := 0             // 0 表示不筛选年份
	month := 0            // 0 表示不筛选月份

	result, err := client.API().GetEWTTransactionDetails(
		page, pageSize,
		transactionType, bizType,
		year, month,
	)
	if err != nil {
		log.Printf("查询权证交易明细失败: %v\n", err)
		return
	}
	if !result.Success {
		log.Printf("查询权证交易明细失败: %s\n", result.Message)
		return
	}

	fmt.Printf("权证交易明细查询结果: %#v\n", result.Data)
}
