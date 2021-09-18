package registry

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/xlt/edu_srv/edu_ad_srv/global"
	"go.uber.org/zap"
)

type Registry interface {
	Register(address string, port int, name string, tags []string, id string) error
	DeRegister(serviceId string) error
}

func NewRegistry() Registry {
	return &Consul{
		Host: global.ServerConfig.ConsulConfig.Host,
		Port: global.ServerConfig.ConsulConfig.Port,
	}
}

type Consul struct {
	Host string
	Port int
}

func (r *Consul) Register(address string, port int, name string, tags []string, id string) error {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", r.Host, r.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	//生成对应的检查对象
	healthCheck := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", global.ServerConfig.Host, port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s",
	}

	//生成注册对象
	serviceRegistration := new(api.AgentServiceRegistration)
	serviceRegistration.Name = name
	serviceRegistration.ID = id
	serviceRegistration.Port = port
	serviceRegistration.Tags = tags
	serviceRegistration.Address = address
	serviceRegistration.Check = healthCheck
	if err = client.Agent().ServiceRegister(serviceRegistration); err != nil {
		zap.S().Errorw("client.Agent().ServiceRegister failed", "msg", err.Error())
		return err
	}
	return nil
}

func (r *Consul) DeRegister(serviceId string) error {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", r.Host, r.Port)
	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}
	return client.Agent().ServiceDeregister(serviceId)
}
