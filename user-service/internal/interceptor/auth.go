package interceptor

import (
	"context"
	"my-IMSystem/pkg/jwt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// AuthInterceptor 用于验证 JWT 的 gRPC 拦截器
func AuthInterceptor(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	// 跳过无需认证的方法
	if info.FullMethod == "/user.User/Register" || info.FullMethod == "/user.User/Login" {
		return handler(ctx, req)
	}

	// 从 metadata 中获取 token
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "Missing metadata")
	}

	tokens := md["authorization"]
	if len(tokens) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "Authorization token not provided")
	}

	claims, err := jwt.ParseToken(tokens[0])
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid token: %v", err)
	}

	// 将 UID 写入上下文（方便业务层获取）
	ctx = context.WithValue(ctx, "uid", claims.Uid)

	// 调用后续处理器
	return handler(ctx, req)
}
