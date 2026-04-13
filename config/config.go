package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"os"
)

var (
	Port          string
	JwtSignKey    string
	JwtExpireTime int64 // 过期时间,单位：分钟
	Username      string
	Password      string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env not found")
	}
	viper.SetDefault("PORT", ":8081")
	viper.AutomaticEnv()
	Port = viper.GetString("PORT")
}

func GetNotionToken() string {
	token := os.Getenv("NOTION_TOKEN")
	if token == "" {
		panic("NOTION_TOKEN 未配置")
	}
	return token
}
