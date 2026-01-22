package config

import (
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf

	// 数据库配置 (PostgreSQL)
	Database struct {
		Host     string
		Port     int
		User     string
		Password string
		DBName   string
		SSLMode  string
	}

	// Redis配置
	Redis struct {
		Host     string
		Password string
		DB       int
	}

	// JWT配置
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
}
