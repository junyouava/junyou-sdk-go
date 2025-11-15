package main

import (
	"fmt"
	"log"

	junyousdk "github.com/junyouava/junyou-sdk-go"
)

func main() {
	// 初始化客户端
	config := junyousdk.DefaultConfig().
		WithAccessId("your-access-id").
		WithAccessKey("your-access-key")

	client, err := junyousdk.NewClient(config)
	if err != nil {
		log.Fatalf("创建客户端失败: %v\n", err)
	}

	// 示例1: 注册
	registerExample(client)

	// 示例2: 获取登录令牌
	loginExample(client)

	// 示例3: 释放权证
	ewtExample(client)
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
