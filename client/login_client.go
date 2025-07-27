package main

import (
	"context"
	"fmt"
	"my-IMSystem/user-service/user"
	"time"

	"google.golang.org/grpc"
)

func main() {
	// 1. 连接 gRPC 服务端
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		panic(fmt.Sprintf("连接失败: %v", err))
	}
	defer conn.Close()

	// 2. 创建 gRPC 客户端
	client := user.NewUserClient(conn)

	// 3. 构造请求
	req := &user.LoginRequest{
		Username: "imuser", // ← 你数据库里已有的用户名
		Password: "im123456",   // ← 密码
	}

	// 4. 调用 Login 方法
	resp, err := client.Login(context.Background(), req)
	if err != nil {
		fmt.Printf("调用失败: %v\n", err)
		return
	}

	// 5. 打印返回结果
	fmt.Printf("登录结果：\nToken: %s\nMessage: %s\n", resp.Token, resp.Message)
}
