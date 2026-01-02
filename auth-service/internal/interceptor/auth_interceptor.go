package interceptor

import (
	"context"
	"strings"

	"my-IMSystem/pkg/jwt"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type contextKey string

const (
	uidKey      contextKey = "uid"
	deviceIdKey contextKey = "deviceId"
)

func AuthInterceptor(secretKey string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		//  放行的接口（无需校验 token）
		unauthMethods := []string{
			"/auth.Auth/Register",
			"/auth.Auth/Login",
			"/auth.Auth/VerifyToken",
			"/auth.Auth/ParseToken",
			"/auth.Auth/RefreshToken",
			"/auth.Auth/GenerateToken",
		}

		// 匹配是否在放行列表中
		for _, m := range unauthMethods {
			if info.FullMethod == m {
				return handler(ctx, req)
			}
		}

		// 从 metadata 中获取 authorization 字段
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "缺少 metadata")
		}

		authHeader := md["authorization"]
		if len(authHeader) == 0 {
			return nil, status.Error(codes.Unauthenticated, "未携带 token")
		}

		// Bearer token 或直接就是 token
		token := strings.TrimPrefix(authHeader[0], "Bearer ")

		// 解析 token
		claims, err := jwt.ParseToken(token, []byte(secretKey))
		if err != nil {
			logx.Errorf("token 解析失败: %v", err)
			return nil, status.Error(codes.Unauthenticated, "token 无效: "+err.Error())
		}

		// 注入 userId / deviceId 到 context 中，供下游 logic 使用
		ctx = context.WithValue(ctx, uidKey, claims.Uid)
		ctx = context.WithValue(ctx, deviceIdKey, claims.DeviceId)

		// 继续处理请求
		return handler(ctx, req)
	}
}
