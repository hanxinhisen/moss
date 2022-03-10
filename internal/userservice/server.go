// Created by Hisen at 2022/3/2.
package userservice

import (
	"context"
	"fmt"
	genericoptions "github.com/hanxinhisen/moss/internal/pkg/options"
	genericapiserver "github.com/hanxinhisen/moss/internal/pkg/server"
	"github.com/hanxinhisen/moss/internal/userservice/config"
	cachev1 "github.com/hanxinhisen/moss/internal/userservice/controller/v1/cache"
	pb "github.com/hanxinhisen/moss/internal/userservice/proto/v1"
	"github.com/hanxinhisen/moss/internal/userservice/store"
	"github.com/hanxinhisen/moss/internal/userservice/store/mysql"
	"github.com/hanxinhisen/moss/pkg/log"
	"github.com/hanxinhisen/moss/pkg/shutdown"
	"github.com/hanxinhisen/moss/pkg/shutdown/shutdownmanagers/posixsignal"
	"github.com/hanxinhisen/moss/pkg/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

type apiServer struct {
	gs               *shutdown.GracefulShutdown
	redisOptions     *genericoptions.RedisOptions
	gRPCAPIServer    *grpcAPIServer
	genericAPIServer *genericapiserver.GenericAPIServer
}

type ExtraConfig struct {
	Addr         string
	MaxMsgSize   int
	ServerCert   genericoptions.GeneratableKeyCert
	mysqlOptions *genericoptions.MySQLOptions
	// etcdOptions      *genericoptions.EtcdOptions
}
type completedExtraConfig struct {
	*ExtraConfig
}

func (c *ExtraConfig) complete() *completedExtraConfig {
	if c.Addr == "" {
		c.Addr = "127.0.0.1:8081"
	}
	return &completedExtraConfig{c}
}

func (c *completedExtraConfig) New() (*grpcAPIServer, error) {
	creds, err := credentials.NewServerTLSFromFile(c.ServerCert.CertKey.CertFile, c.ServerCert.CertKey.KeyFile)
	if err != nil {
		log.Fatalf("Failed to generate credentials %s", err.Error())
	}
	opts := []grpc.ServerOption{grpc.MaxRecvMsgSize(c.MaxMsgSize), grpc.Creds(creds)}
	grpcServer := grpc.NewServer(opts...)

	storeIns, _ := mysql.GetMysqlFactoryOr(c.mysqlOptions)
	// storeIns, _ := etcd.GetEtcdFactoryOr(c.etcdOptions, nil)
	store.SetClient(storeIns)
	cacheIns, err := cachev1.GetCacheInsOr(storeIns)
	if err != nil {
		log.Fatalf("Failed to get cache instance: %s", err.Error())
	}

	pb.RegisterCacheServer(grpcServer, cacheIns)

	reflection.Register(grpcServer)

	return &grpcAPIServer{grpcServer, c.Addr}, nil
}

func createAPIServer(cfg *config.Config) (*apiServer, error) {
	// 创建优雅关闭
	gs := shutdown.New()
	gs.AddShutdownManager(posixsignal.NewPosixSignalManager())

	//构建默认配置文件
	genericConfig, err := buildGenericConfig(cfg)
	if err != nil {
		return nil, err
	}
	//构建附加配置文件，这里主要grpc配置文件
	extraConfig, err := buildExtraConfig(cfg)
	if err != nil {
		return nil, err
	}
	// 生成默认http服务
	genericServer, err := genericConfig.Complete().New()
	if err != nil {
		return nil, err
	}
	// 生成grpc服务
	extraServer, err := extraConfig.complete().New()
	if err != nil {
		return nil, err
	}
	server := &apiServer{
		gs:               gs,
		redisOptions:     cfg.RedisOptions,
		genericAPIServer: genericServer,
		gRPCAPIServer:    extraServer,
	}
	return server, nil

}

func buildExtraConfig(cfg *config.Config) (*ExtraConfig, error) {
	return &ExtraConfig{
		Addr:         fmt.Sprintf("%s:%d", cfg.GRPCOptions.BindAddress, cfg.GRPCOptions.BindPort),
		MaxMsgSize:   cfg.GRPCOptions.MaxMsgSize,
		ServerCert:   cfg.SecureServing.ServerCert,
		mysqlOptions: cfg.MySQLOptions,
	}, nil
}

func buildGenericConfig(cfg *config.Config) (genericConfig *genericapiserver.Config, lastErr error) {
	genericConfig = genericapiserver.NewConfig()
	if lastErr = cfg.GenericServerRunOptions.ApplyTo(genericConfig); lastErr != nil {
		return
	}
	if lastErr = cfg.FeatureOptions.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	if lastErr = cfg.SecureServing.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	if lastErr = cfg.InsecureServing.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	return

}

type preparedAPIServer struct {
	*apiServer
}

func (s *apiServer) PrepareRun() preparedAPIServer {
	// 初始化路由
	initRoute(s.genericAPIServer.Engine)
	// 初始化redis
	s.initRedisStore()
	s.gs.AddShutdownCallback(shutdown.ShutdownFunc(func(string string) error {
		mysqlStore, _ := mysql.GetMysqlFactoryOr(nil)
		if mysqlStore != nil {
			_ = mysqlStore.Close()
		}
		s.gRPCAPIServer.Close()
		s.genericAPIServer.Close()
		return nil
	}))
	return preparedAPIServer{s}
}

func (s preparedAPIServer) Run() error {
	// 启动grpc
	go s.gRPCAPIServer.Run()

	// start shutdown managers
	if err := s.gs.Start(); err != nil {
		log.Fatalf("start shutdown manager failed: %s", err.Error())
	}
	// 启动http服务
	return s.genericAPIServer.Run()
}

func (s *apiServer) initRedisStore() {
	ctx, cancel := context.WithCancel(context.Background())

	s.gs.AddShutdownCallback(shutdown.ShutdownFunc(func(s string) error {
		cancel()
		return nil
	}))

	config := &storage.Config{
		Host:                  s.redisOptions.Host,
		Port:                  s.redisOptions.Port,
		Addrs:                 s.redisOptions.Addrs,
		MasterName:            s.redisOptions.MasterName,
		Username:              s.redisOptions.Username,
		Password:              s.redisOptions.Password,
		Database:              s.redisOptions.Database,
		MaxIdle:               s.redisOptions.MaxIdle,
		MaxActive:             s.redisOptions.MaxActive,
		Timeout:               s.redisOptions.Timeout,
		EnableCluster:         s.redisOptions.EnableCluster,
		UseSSL:                s.redisOptions.UseSSL,
		SSLInsecureSkipVerify: s.redisOptions.SSLInsecureSkipVerify,
	}
	go storage.ConnectToRedis(ctx, config)

}
