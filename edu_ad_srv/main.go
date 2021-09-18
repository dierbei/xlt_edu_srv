package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/xlt/edu_srv/edu_ad_srv/global"
	"github.com/xlt/edu_srv/edu_ad_srv/internal/handler"
	"github.com/xlt/edu_srv/edu_ad_srv/internal/proto"
	"github.com/xlt/edu_srv/edu_ad_srv/pkg/initialize"
	"github.com/xlt/edu_srv/edu_ad_srv/pkg/registry"
	"github.com/xlt/edu_srv/edu_ad_srv/pkg/utils"
)

func main() {
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitMySQL()

	freePort := utils.GetFreePort()

	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", global.ServerConfig.Host, freePort))
	if err != nil {
		zap.S().Errorw("net.Listen failed, err:", "msg", err.Error())
	}

	server := grpc.NewServer()
	proto.RegisterSpaceServer(server, &handler.SpaceServer{})
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	// uuid 保证服务唯一
	serviceID := fmt.Sprintf("%s", uuid.NewV4())

	if err := registry.NewRegistry().Register(
		global.ServerConfig.Host,
		freePort,
		global.ServerConfig.ServerName,
		global.ServerConfig.Tags,
		serviceID,
	); err != nil {
		zap.S().Errorw("registry.NewRegistry().Register failed", "msg", err.Error())
	}

	zap.S().Infow("server.Serve success", "port", freePort, "serviveID", serviceID)
	go func() {
		err = server.Serve(listen)
		if err != nil {
			zap.S().Errorw("server.Serve failed, err:", "msg", err.Error())
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err := registry.NewRegistry().DeRegister(serviceID); err != nil {
		zap.S().Errorw("registry.NewRegistry().DeRegister failed", "msg", err.Error())
	}
	zap.S().Infow("注销服务成功", "port", freePort, "serviveID", serviceID)
}
