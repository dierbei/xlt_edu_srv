package registry

import (
	"go.uber.org/zap"

	"github.com/xlt/edu_srv/edu_ad_srv/global"
)

func RegisterServer(serviceID string, freePort int) {
	if err := NewRegistry().Register(
		global.ServerConfig.Host,
		freePort,
		global.ServerConfig.ServerName,
		global.ServerConfig.Tags,
		serviceID,
	); err != nil {
		zap.S().Errorw("registry.NewRegistry().Register failed", "msg", err.Error())
	}
}
