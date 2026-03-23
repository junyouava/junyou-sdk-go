package main

import (
	"encoding/json"
	"fmt"
	"log"

	junyousdk "github.com/junyouava/junyou-sdk-go"
)

func main() {
	// 初始化客户端
	config := junyousdk.DefaultConfig().
		WithAccessId("f6262a470f329acfa688ef77ccf9c24d").
		WithAccessKey("69klpXN0+k2ZC7JdIyUJ+G/SMbQ/krUFu1tj2Mz2Mjs=").
		WithAddress("http://127.0.0.1")

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

	// 示例5: 权证合伙人释放（预提交 + 提交，与 GOC 示例同一套路）
	// 对应接口: POST /api/open/v1/ewt/pre_ewt_rbp_open 、 POST /api/open/v1/ewt/commit_ewt_rbp
	// ewtReleaseByPartnerExample(client)

	// 示例5b: GOC 预提交 + 提交（PreRewardGOC → json.Marshal(Data) → 签名 → RewardGOC）
	// 对应接口: POST /api/open/v1/goc/pre_reward 、 POST /api/open/v1/goc/reward
	gocRewardExample(client)

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

// ewtReleaseByPartnerExample 权证合伙人释放：AuthLogin → 预提交 → 取待签名对象 → 提交
// 对应接口: POST /api/open/v1/ewt/pre_ewt_rbp_open 、 POST /api/open/v1/ewt/commit_ewt_rbp
// 预提交须带 X-Open-Auth（先 AuthLogin）。message = string(json.Marshal(preResult.Data))，与 GOC 示例一致。
// TODO: public_key、der_hex 请换为密盾对 message 的真实签名结果。
func ewtReleaseByPartnerExample(client *junyousdk.Client) {
	fmt.Println("\n=== 权证合伙人释放：预提交 + 提交示例 ===")

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

	msgBytes, err := json.Marshal(preResult.Data)
	if err != nil {
		log.Printf("序列化 message 失败: %v\n", err)
		return
	}
	message := string(msgBytes)

	bizNo, _ := preResult.Data["biz_no"].(string)
	if bizNo == "" {
		log.Printf("预提交 data 缺少 biz_no\n")
		return
	}

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

	result, err := client.API().GetEWTBalance(1, 20, "")
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
		"",
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

// gocRewardExample GOC：预提交 → message=string(json.Marshal(preResult.Data)) → 签名 → 提交
// 对应接口: POST /api/open/v1/goc/pre_reward 、 POST /api/open/v1/goc/reward
// TODO: public_key、der_hex 请换为密盾对 message 的真实签名结果。
func gocRewardExample(client *junyousdk.Client) {
	fmt.Println("\n=== GOC 预提交 + 提交示例 ===")

	preReq := junyousdk.PreGOCRewardRequest{
		OpenId: "8d007704b1954336e0928c465745c1e87782f5390c1ec784722e63eadf6af6bf",
		Amount: "1.00",
	}

	preResult, err := client.API().PreRewardGOC(preReq)
	if err != nil {
		log.Printf("GOC 预提交失败: %v\n", err)
		return
	}
	if !preResult.Success {
		log.Printf("GOC 预提交失败: %s\n", preResult.Message)
		return
	}
	fmt.Printf("预提交成功，返回数据: %#v\n", preResult.Data)

	msgBytes, err := json.Marshal(preResult.Data)
	if err != nil {
		log.Printf("序列化 message 失败: %v\n", err)
		return
	}
	message := string(msgBytes)
	fmt.Println("message", message)

	bizNo, _ := preResult.Data["biz_no"].(string)
	if bizNo == "" {
		log.Printf("预提交 data 缺少 biz_no\n")
		return
	}

	commitReq := junyousdk.CommitGOCRewardRequest{
		BizNo:     "20260323134001989317",
		Message:   `{"amount":"1.00","biz_desc":"GOC奖励","biz_no":"20260323134001989317","biz_type":"1002","from":"RnSFU7fC7gNSBJtnUrKo7f1yBmhAm9hqh","to":"Wv8Tpva2QjtH3LzYweUvDDBorTV47VjSu"}`,
		PublicKey: `{"Curvname":"P-256","X":"N7H1oyS-4s-ZbsoS4orISoDDYP-QGXlUCDEU0jeicOM","Y":"-bxc-QAmzoWoDF8AJambw77IL9mB9iKHgdX2kT0sVL0"}`,
		DerHex:    "304502206256c1bb7fd9f4ec17a4dc6150e955eda450b1e1bbd3610ebc2e0b73cdc2348d022100f1633e11aa3bfe580fec901e312e827a60828e31015e289e8f1a58f562241d50",
	}

	commitResult, err := client.API().RewardGOC(commitReq)
	if err != nil {
		log.Printf("GOC 提交失败: %v\n", err)
		return
	}
	if !commitResult.Success {
		log.Printf("GOC 提交失败: %s\n", commitResult.Message)
		return
	}

	fmt.Printf("GOC 提交成功，返回数据: %#v\n", commitResult)
}
