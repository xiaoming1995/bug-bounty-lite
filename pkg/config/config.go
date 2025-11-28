package config

import (
	"fmt"
	"github.com/spf13/viper"
)

// Config 聚合所有配置结构体
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	DSN     string `mapstructure:"dsn"`
	MaxIdle int    `mapstructure:"max_idle"`
	MaxOpen int    `mapstructure:"max_open"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	Expire int    `mapstructure:"expire"`
}

// LoadConfig 读取配置文件的核心函数
func LoadConfig() *Config {
	// 1. 设置配置文件的名字和类型
	viper.SetConfigName("config") // 对应 config.yaml 的文件名 (不含后缀)
	viper.SetConfigType("yaml")   // 文件类型

	// 2. 设置查找路径 (按顺序查找)
	// 建议添加多个路径，防止在不同目录下运行 go run 找不到文件
	viper.AddConfigPath("./config") // 相对路径 config/
	viper.AddConfigPath(".")        // 当前目录

	// 3. 读取环境变量 (可选，用于 Docker 部署时覆盖配置)
	viper.AutomaticEnv()

	// 4. 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		// 如果读取失败，直接 panic，因为没有配置程序无法运行
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	// 5. 映射到结构体
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		panic(fmt.Errorf("Unable to decode into struct: %w \n", err))
	}

	// 简单打印一下，确认加载成功 (生产环境建议用 Logger)
	fmt.Printf("✅ Configuration loaded successfully. Mode: %s\n", cfg.Server.Mode)

	return &cfg
}