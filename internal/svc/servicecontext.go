package svc

import (
	"fmt"

	"accesscontrol/internal/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type ServiceContext struct {
	Config config.Config
	DB     *sqlx.DB
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

	return &ServiceContext{
		Config: c,
		DB:     db,
	}
}
