package registry

import (
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func Close(serviceID string, freePort int) {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err := NewRegistry().DeRegister(serviceID); err != nil {
		zap.S().Errorw("registry.NewRegistry().DeRegister failed", "msg", err.Error())
	}
	zap.S().Infow("注销服务成功", "port", freePort, "serviveID", serviceID)
}
