package logic

import (
	"accesscontrol/internal/model"
	"accesscontrol/internal/svc"
	"accesscontrol/internal/types"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
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

func (l *UserLogic) Login(req *types.LoginRequest) (*types.Response, error) {
	// 验证参数
	if req.PhoneNumber == "" || req.Password == "" {
		return &types.Response{
			Code:    400,
			Message: "手机号或密码不能为空",
		}, nil
	}

	// 从数据库查询用户
	var user model.User
	err := l.svcCtx.DB.Get(&user, `SELECT id, password, "phoneNumber", created_at, updated_at FROM users WHERE "phoneNumber" = $1`, req.PhoneNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return &types.Response{
				Code:    401,
				Message: "手机号或密码错误",
			}, nil
		}
		logx.Errorf("查询用户失败: %v", err)
		return &types.Response{
			Code:    500,
			Message: "系统错误",
		}, nil
	}

	// 验证密码（SHA256 加密比对）
	hashedPassword := sha256.Sum256([]byte(req.Password))
	if user.Password != hex.EncodeToString(hashedPassword[:]) {
		return &types.Response{
			Code:    401,
			Message: "手机号或密码错误",
		}, nil
	}

	// 生成JWT token
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.Auth.AccessExpire
	accessToken, err := l.generateToken(now, accessExpire, user.Id)
	if err != nil {
		return nil, err
	}
	fmt.Printf("user.Id user.PhoneNumber: %d , %s , now: %d accessExpire: %d\n", user.Id, user.PhoneNumber, now, accessExpire)

	return &types.Response{
		Code:    0,
		Message: "登录成功",
		Data: types.LoginResponse{
			AccessToken:  accessToken,
			AccessExpire: now + accessExpire,
			UserInfo: model.User{
				Id:          user.Id,
				PhoneNumber: user.PhoneNumber,
				CreatedAt:   user.CreatedAt,
				UpdatedAt:   user.UpdatedAt,
			},
		},
	}, nil
}

func (l *UserLogic) generateToken(iat, exp int64, userId int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + exp
	claims["iat"] = iat
	claims["userId"] = userId
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(l.svcCtx.Config.Auth.AccessSecret))
}
