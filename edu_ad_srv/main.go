package main

import (
	"fmt"
	"net"

	uuid "github.com/satori/go.uuid"
	"github.com/xlt/edu_srv/edu_ad_srv/pkg/registry"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/xlt/edu_srv/edu_ad_srv/global"
	"github.com/xlt/edu_srv/edu_ad_srv/internal/handler"
	"github.com/xlt/edu_srv/edu_ad_srv/internal/proto"
	"github.com/xlt/edu_srv/edu_ad_srv/pkg/initialize"
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
	// 健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	serviceID := fmt.Sprintf("%s", uuid.NewV4())
	// 注册到Consul
	registry.RegisterServer(serviceID, freePort)

	zap.S().Infow("server.Serve success", "port", freePort, "serviveID", serviceID)
	go func() {
		err = server.Serve(listen)
		if err != nil {
			zap.S().Errorw("server.Serve failed, err:", "msg", err.Error())
		}
	}()

	// 注销服务回收资源
	registry.Close(serviceID, freePort)
}
