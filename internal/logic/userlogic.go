package logic

import (
	"accesscontrol/internal/model"
	"accesscontrol/internal/svc"
	"accesscontrol/internal/types"
	"context"
	"database/sql"
	"net/http"
	"github.com/zeromicro/go-zero/core/logx"
)

type UserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLogic {
	return &UserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserLogic) QueryUser(req *types.LoginRequest) (*types.Response, error) {
	// 验证参数
	if req.PhoneNumber == "" {
		return &types.Response{
			Code:    400,
			Message: "手机号不能为空",
		}, nil
	}

	// 从数据库查询用户
	var user model.User
	err := l.svcCtx.DB.Get(&user, `SELECT id, "phoneNumber", status, "validTime" FROM users WHERE "phoneNumber" = $1`, req.PhoneNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return &types.Response{
				Code:    401,
				Message: "用户不存在",
			}, nil
		}
		logx.Errorf("查询用户失败: %v", err)
		return &types.Response{
			Code:    500,
			Message: "系统错误",
		}, nil
	}

	return &types.Response{
		Code:    http.StatusOK,
		Message: "查询成功",
		Data:    user,
	}, nil
}

func (l *UserLogic) AddUser(req *model.User) (*types.Response, error) {
	// 添加用户
	if req.PhoneNumber == "" || req.ValidTime == "" {
		return &types.Response{
			Code:    400,
			Message: "手机号或有效时间不能为空",
		}, nil
	}
	// 检查用户是否已存在
	var user model.User
	err := l.svcCtx.DB.Get(&user, `SELECT id, "phoneNumber", status, "validTime" FROM users WHERE "phoneNumber" = $1`, req.PhoneNumber)
	if err == nil {
		return &types.Response{
			Code:    409,
			Message: "用户已存在",
		}, nil
	}

	// 添加新用户
	_, err = l.svcCtx.DB.Exec(`INSERT INTO users ("phoneNumber", status, "validTime") VALUES ($1, $2, $3)`,
		req.PhoneNumber, 1, req.ValidTime)
	if err != nil {
		logx.Errorf("添加用户失败: %v", err)
		return &types.Response{
			Code:    500,
			Message: "系统错误",
		}, nil
	}

	return &types.Response{
		Code:    http.StatusOK,
		Message: "用户创建成功",
	}, nil
}

// 编辑用户
func (l *UserLogic) EditUser(req *model.User) (*types.Response, error) {
	// 编辑用户
	if req.Id == 0 || req.PhoneNumber == "" || req.ValidTime == "" {
		return &types.Response{
			Code:    400,
			Message: "用户ID、手机号或有效时间不能为空",
		}, nil
	}
	// 检查用户是否存在
	var user model.User
	err := l.svcCtx.DB.Get(&user, `SELECT id, "phoneNumber", status, "validTime" FROM users WHERE id = $1`, req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return &types.Response{
				Code:    404,
				Message: "用户不存在",
			}, nil
		}
		logx.Errorf("查询用户失败: %v", err)
		return &types.Response{
			Code:    500,
			Message: "系统错误",
		}, nil
	}

	// 更新用户信息
	_, err = l.svcCtx.DB.Exec(`UPDATE users SET "phoneNumber" = $1, status = $2, "validTime" = $3 WHERE id = $4`,
		req.PhoneNumber, req.Status, req.ValidTime, req.Id)
	if err != nil {
		logx.Errorf("更新用户失败: %v", err)
		return &types.Response{
			Code:    500,
			Message: "系统错误",
		}, nil
	}

	return &types.Response{
		Code:    http.StatusOK,
		Message: "用户信息更新成功",
	}, nil
}

// 删除用户 根据手机号码
func (l *UserLogic) DeleteUser(phoneNumber string) (*types.Response, error) {
	// 删除用户
	if phoneNumber == "" {
		return &types.Response{
			Code:    400,
			Message: "手机号不能为空",
		}, nil
	}
	_, err := l.svcCtx.DB.Exec(`DELETE FROM users WHERE "phoneNumber" = $1`, phoneNumber)
	if err != nil {
		logx.Errorf("删除用户失败: %v", err)
		return &types.Response{
			Code:    500,
			Message: "系统错误",
		}, nil
	}
	return &types.Response{
		Code:    http.StatusOK,
		Message: "用户删除成功",
	}, nil
}

// 用户列表 分页加载
func (l *UserLogic) ListUsers(page, pageSize int) (*types.Response, error) {
	// 查询总数
	var total int
	err := l.svcCtx.DB.Get(&total, `SELECT COUNT(*) FROM users`)
	if err != nil {
		logx.Errorf("查询用户总数失败: %v", err)
		return &types.Response{
			Code:    500,
			Message: "系统错误",
		}, nil
	}

	// 分页查询用户列表
	var users []model.User
	err = l.svcCtx.DB.Select(&users, `SELECT id, "phoneNumber", status, "validTime" FROM users ORDER BY id LIMIT $1 OFFSET $2`, pageSize, (page-1)*pageSize)
	if err != nil {
		logx.Errorf("查询用户失败: %v", err)
		return &types.Response{
			Code:    500,
			Message: "系统错误",
		}, nil
	}

	return &types.Response{
		Code:    http.StatusOK,
		Message: "用户列表查询成功",
		Data: map[string]interface{}{
			"total": total,
			"list":  users,
		},
	}, nil
}
