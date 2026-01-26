package logic

import (
	"accesscontrol/internal/model"
	"accesscontrol/internal/svc"
	"accesscontrol/internal/types"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
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
		return nil, errors.New("手机号不能为空")
	}

	// 1. 尝试从 Redis 获取缓存
	cacheKey := "user:phone:" + req.PhoneNumber
	var user model.User
	cacheVal, _ := l.svcCtx.Redis.Get(cacheKey)
	if cacheVal != "" {
		if err := json.Unmarshal([]byte(cacheVal), &user); err == nil {
			return &types.Response{
				Code:    http.StatusOK,
				Message: "查询成功（来自缓存）",
				Data:    user,
			}, nil
		}
	}

	// 2. 缓存未命中，从数据库查询用户
	u, err := l.svcCtx.UserRepo.FindOneByPhone(l.ctx, req.PhoneNumber)
	if err != nil {
		if err == sql.ErrNoRows || err.Error() == "sql: no rows in result set" {
			return nil, errors.New("用户不存在")
		}
		logx.Errorf("查询用户失败: %v", err)
		return nil, err
	}
	user = *u

	// 3. 写入 Redis 缓存 (设置 10 分钟过期时间)
	if data, err := json.Marshal(user); err == nil {
		_ = l.svcCtx.Redis.Setex(cacheKey, string(data), 600)
	}

	return &types.Response{
		Code:    http.StatusOK,
		Message: "查询成功",
		Data:    user,
	}, nil
}

func (l *UserLogic) AddUser(req *types.RegisterRequest) (*types.Response, error) {
	// 添加用户
	if req.PhoneNumber == "" || req.ValidTime == "" {
		return nil, errors.New("手机号或有效时间不能为空")
	}
	// 检查用户是否已存在
	_, err := l.svcCtx.UserRepo.FindOneByPhone(l.ctx, req.PhoneNumber)
	if err == nil {
		return nil, errors.New("用户已存在")
	}

	// 添加新用户
	user := &model.User{
		PhoneNumber: req.PhoneNumber,
		Status:      req.Status,
		ValidTime:   req.ValidTime,
	}
	err = l.svcCtx.UserRepo.Insert(l.ctx, user)
	if err != nil {
		logx.Errorf("添加用户失败: %v", err)
		return nil, err
	}

	return &types.Response{
		Code:    http.StatusOK,
		Message: "用户创建成功",
	}, nil
}

// 编辑用户
func (l *UserLogic) EditUser(req *types.UpdateUserRequest) (*types.Response, error) {
	// 编辑用户
	if req.Id == 0 || req.PhoneNumber == "" || req.ValidTime == "" {
		return nil, errors.New("用户ID、手机号或有效时间不能为空")
	}
	// 检查用户是否存在
	user, err := l.svcCtx.UserRepo.FindOne(l.ctx, req.Id)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 更新用户信息
	user.PhoneNumber = req.PhoneNumber
	user.Status = req.Status
	user.ValidTime = req.ValidTime

	err = l.svcCtx.UserRepo.Update(l.ctx, user)
	if err != nil {
		logx.Errorf("更新用户失败: %v", err)
		return nil, err
	}

	// 3. 清除 Redis 缓存 (保证下次查询能拿到新数据)
	cacheKey := "user:phone:" + user.PhoneNumber
	_, _ = l.svcCtx.Redis.Del(cacheKey)

	return &types.Response{
		Code:    http.StatusOK,
		Message: "用户信息更新成功",
	}, nil
}

// 删除用户 根据手机号码
func (l *UserLogic) DeleteUser(phoneNumber string) (*types.Response, error) {
	// 删除用户
	if phoneNumber == "" {
		return nil, errors.New("手机号不能为空")
	}

	err := l.svcCtx.UserRepo.Delete(l.ctx, phoneNumber)
	if err != nil {
		logx.Errorf("删除用户失败: %v", err)
		return nil, err
	}

	// 清除 Redis 缓存
	cacheKey := "user:phone:" + phoneNumber
	_, _ = l.svcCtx.Redis.Del(cacheKey)

	return &types.Response{
		Code:    http.StatusOK,
		Message: "用户删除成功",
	}, nil
}

// 用户列表 分页加载
func (l *UserLogic) ListUsers(page, pageSize int) (*types.Response, error) {
	// 使用 Repository 进行分页查询
	users, total, err := l.svcCtx.UserRepo.List(l.ctx, pageSize, (page-1)*pageSize)
	if err != nil {
		logx.Errorf("查询用户失败: %v", err)
		return nil, err
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
