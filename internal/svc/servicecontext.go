package svc

import (
	"fmt"

	"accesscontrol/internal/config"
	"accesscontrol/internal/repository"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/zeromicro/go-zero/core/limit"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type ServiceContext struct {
	Config      config.Config
	DB          *sqlx.DB
	UserRepo    repository.UserRepository
	Redis       *redis.Redis
	RateLimiter *limit.TokenLimiter
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 构建 PostgreSQL 连接字符串
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)

	// 连接数据库
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		panic(fmt.Sprintf("数据库连接失败: %v", err))
	}

	// 设置连接池参数
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(10)

	// 初始化 Redis
	rds := redis.New(c.Redis.Host, func(r *redis.Redis) {
		r.Type = redis.NodeType
		r.Pass = c.Redis.Password
	})

	// 初始化 RateLimiter (100 req/s)
	limiter := limit.NewTokenLimiter(100, 100, rds, "api-rate-limit")

	return &ServiceContext{
		Config:      c,
		DB:          db,
		UserRepo:    repository.NewUserRepository(db),
		Redis:       rds,
		RateLimiter: limiter,
	}
}
