package handler

import (
	"accesscontrol/internal/svc"
	"net/http"

	"github.com/zeromicro/go-zero/rest"
)

func SetupRoutes(server *rest.Server, svcCtx *svc.ServiceContext) {
	// 公开路由 (无需认证)
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/api/query",
				Handler: QueryUserHandler(svcCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/add",
				Handler: AddUserHandler(svcCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/edit",
				Handler: EditUserHandler(svcCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/delete",
				Handler: DeleteUserHandler(svcCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/list",
				Handler: ListUsersHandler(svcCtx),
			},
		},
	)
}
