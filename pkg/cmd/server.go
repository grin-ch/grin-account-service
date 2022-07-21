package cmd

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/grin-ch/grin-account-service/cfg"
	"github.com/grin-ch/grin-account-service/pkg/model"
	"github.com/grin-ch/grin-account-service/pkg/service"
	"github.com/grin-ch/grin-api/api/account"
	"github.com/grin-ch/grin-api/api/captcha"
	center "github.com/grin-ch/grin-etcd-center"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// RunServer 运行服务
func RunServer() error {
	var cfgName, cfgPath string
	flag.StringVar(&cfgName, "cfgName", "service", "log file")
	flag.StringVar(&cfgPath, "cfgPath", "./cfg", "log path")
	flag.Parse()

	cfg.SetServerConfig(cfgName, cfgPath)

	return grpcServer(cfg.GetConfig())
}

// 运行grpc服务
func grpcServer(c *cfg.ServerConfig) error {
	initLogger(c.LogPath, c.LogLevel, c.LogColor, c.LogCaller)

	grpcListener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", c.Port))
	if err != nil {
		log.Errorf("tcp listen err:%s", err.Error())
		return err
	}

	// grpc server
	svc := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandler(recoveryFunc)),
			grpc_auth.UnaryServerInterceptor(authFunc),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_recovery.StreamServerInterceptor(grpc_recovery.WithRecoveryHandler(recoveryFunc)),
			grpc_auth.StreamServerInterceptor(authFunc),
		),
	)
	gracefulShutdown(svc)

	// 获取注册中心连接
	etcdCenter, err := center.NewEtcdCenter(c.RegEndpoint, c.RegTimeout)
	if err != nil {
		log.Errorf("new etcd client err:%s", err.Error())
		return err
	}
	// 服务注册
	registrar := etcdCenter.Registrar(c.Name, c.Host, c.Port, c.RegTimeout)
	err = registrar.Registry()
	if err != nil {
		log.Errorf("server registry err: %s", err.Error())
		return err
	}
	defer registrar.Deregistry()

	// 注册本机grpc服务
	if err := registryGrpcServices(svc, c, etcdConncter(etcdCenter.Builder())); err != nil {
		log.Errorf("grpc server registry err:%s", err.Error())
		return err
	}
	log.Infof("grpc server is running: %s", fmt.Sprintf("%s:%d", c.Host, c.Port))
	return svc.Serve(grpcListener)
}

// 优雅退出
func gracefulShutdown(svc *grpc.Server) {
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		sign := <-quit

		// 关闭服务链接
		svc.GracefulStop()
		log.Infof("grpc server shutdown: %s", sign)
		os.Exit(1)
	}()
}

type connecter func(string) (*grpc.ClientConn, error)

func registryGrpcServices(svc *grpc.Server, c *cfg.ServerConfig, connect connecter) error {
	pv, err := model.RegistryDatabase(c.Dsn())
	if err != nil {
		log.Errorf("registry database error:%s", err.Error())
		return err
	}

	// 加载服务connect
	captchaConn, err := connect("grin-captcha-service")
	if err != nil {
		return err
	}
	captchaClient := captcha.NewCaptchaServiceClient(captchaConn)
	account.RegisterUserServiceServer(svc, service.NewUserService(pv, captchaClient))

	return nil
}
