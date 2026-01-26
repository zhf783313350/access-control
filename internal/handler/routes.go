package handler

import (
	"accesscontrol/internal/middleware"
	"accesscontrol/internal/svc"
	"net/http"

	"github.com/zeromicro/go-zero/rest"
)

func SetupRoutes(server *rest.Server, serverCtx *svc.ServiceContext) {
	// 初始化中间件
	rateLimitMatch := middleware.NewRateLimitMiddleware(serverCtx.RateLimiter).Handle

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/api/query",
				Handler: rateLimitMatch(QueryUserHandler(serverCtx)),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/add",
				Handler: rateLimitMatch(AddUserHandler(serverCtx)),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/edit",
				Handler: rateLimitMatch(EditUserHandler(serverCtx)),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/delete",
				Handler: rateLimitMatch(DeleteUserHandler(serverCtx)),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/list",
				Handler: rateLimitMatch(ListUsersHandler(serverCtx)),
			},
		},
	)
}
