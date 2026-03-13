package middleware

import (
	"context"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/martketplace-vkr/pkg/server/http"
)

const (
	userIDKey     = "userID"
	userKey       = "user"
	permissionKey = "x-user-permissionkey"
	userID        = "x-user-id"
	userLogin     = "x-user-login"
	userRole      = "x-user-role"
)

func Auth() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		md, _ := metadata.FromIncomingContext(ctx)

		user := &http.User{}

		permissionKeyHeaderVal, ok := getHeaderValue(md, permissionKey)
		if !ok {
			return handler(ctx, req)
		}

		permission, err := strconv.Atoi(permissionKeyHeaderVal)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "missing permission key")
		}
		user.PermissionKey = int64(permission)

		userIDHeaderValue, ok := getHeaderValue(md, userID)
		if !ok {
			return handler(ctx, req)
		}

		userID, err := strconv.Atoi(userIDHeaderValue)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "missing user id")
		}

		user.ID = int64(userID)

		user.Login, ok = getHeaderValue(md, userLogin)
		if !ok {
			return handler(ctx, req)
		}

		user.Role, ok = getHeaderValue(md, userRole)
		if !ok {
			return handler(ctx, req)
		}

		ctx = context.WithValue(ctx, userKey, user)

		return handler(ctx, req)
	}
}

func getHeaderValue(md metadata.MD, key string) (value string, ok bool) {
	values := md.Get(key)

	if len(values) > 0 {
		return values[0], true
	}

	return "", false
}
