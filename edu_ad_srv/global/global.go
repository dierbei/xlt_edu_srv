package global

import (
	"github.com/xlt/edu_srv/edu_ad_srv/config"
	"gorm.io/gorm"
)

var (
	// Nacos
	NacosConfig = &config.NacosConfig{}

	// 全局服务配置
	ServerConfig = &config.ServerConfig{}

	// MySQL
	MySQLConn *gorm.DB
)
