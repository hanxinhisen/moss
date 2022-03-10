// Created by Hisen at 2022/3/2.
package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/hanxinhisen/moss/internal/pkg/middleware"
	"github.com/hanxinhisen/moss/pkg/log"
	"github.com/marmotedu/component-base/pkg/core"
	"github.com/marmotedu/component-base/pkg/version"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strings"
	"time"
)

type GenericAPIServer struct {
	middlewares         []string
	SecureServingInfo   *SecureServingInfo
	InsecureServingInfo *InsecureServingInfo
	ShutdownTimeout     time.Duration

	*gin.Engine
	healthz                      bool
	enableMetrics                bool
	enableProfiling              bool
	insecureServer, secureServer *http.Server
}

func initGenericAPIServer(s *GenericAPIServer) {
	s.Setup()
	s.InstallMiddlewares()
	s.InstallAPIS()
}

func (s *GenericAPIServer) Setup() {
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Infof("%-6s %-s --> %s (%d handlers)", httpMethod, absolutePath, handlerName, nuHandlers)

	}
}

func (s *GenericAPIServer) InstallMiddlewares() {
	s.Use(middleware.RequestID())
	s.Use(middleware.Context())
	for _, m := range s.middlewares {
		mw, ok := middleware.Middlewares[m]
		if !ok {
			log.Warnf("can not found middleware:%s", m)
			continue
		}
		log.Infof("install middleware %s", m)
		s.Use(mw)
	}
}
func (s *GenericAPIServer) InstallAPIS() {
	// 原生api server 是否开启健康监测接口
	if s.healthz {
		s.GET("/healthz", func(context *gin.Context) {
			core.WriteResponse(context, nil, map[string]string{"status": "ok"})
		})
	}
	if s.enableMetrics {
		prometheus := ginprometheus.NewPrometheus("gin")
		prometheus.Use(s.Engine)
	}
	if s.enableProfiling {
		pprof.Register(s.Engine)
	}

	s.GET("/version", func(context *gin.Context) {
		core.WriteResponse(context, nil, version.Get())
	})

}

func (s *GenericAPIServer) Close() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.secureServer.Shutdown(ctx); err != nil {
		log.Warnf("Shutdown secure server failed: %s", err.Error())
	}

	if err := s.insecureServer.Shutdown(ctx); err != nil {
		log.Warnf("Shutdown insecure server failed: %s", err.Error())
	}

}
func (s *GenericAPIServer) Run() error {
	// For scalability, use custom HTTP configuration mode here
	s.insecureServer = &http.Server{
		Addr:    s.InsecureServingInfo.Address,
		Handler: s,
		// ReadTimeout:    10 * time.Second,
		// WriteTimeout:   10 * time.Second,
		// MaxHeaderBytes: 1 << 20,

	}

	// For scalability, use custom HTTP configuration mode here
	s.secureServer = &http.Server{
		Addr:    s.SecureServingInfo.Address(),
		Handler: s,
		// ReadTimeout:    10 * time.Second,
		// WriteTimeout:   10 * time.Second,
		// MaxHeaderBytes: 1 << 20,
	}

	var eg errgroup.Group

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	eg.Go(func() error {
		log.Infof("Start to listening the incoming requests on http address: %s", s.InsecureServingInfo.Address)

		if err := s.insecureServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err.Error())

			return err
		}

		log.Infof("Server on %s stopped", s.InsecureServingInfo.Address)

		return nil
	})

	eg.Go(func() error {
		key, cert := s.SecureServingInfo.CertKey.KeyFile, s.SecureServingInfo.CertKey.CertFile
		if cert == "" || key == "" || s.SecureServingInfo.BindPort == 0 {
			return nil
		}

		log.Infof("Start to listening the incoming requests on https address: %s", s.SecureServingInfo.Address())

		if err := s.secureServer.ListenAndServeTLS(cert, key); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err.Error())

			return err
		}

		log.Infof("Server on %s stopped", s.SecureServingInfo.Address())

		return nil
	})

	// Ping the server to make sure the router is working.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if s.healthz {
		if err := s.ping(ctx); err != nil {
			return err
		}
	}

	if err := eg.Wait(); err != nil {
		log.Fatal(err.Error())
	}

	return nil
}

func (s *GenericAPIServer) ping(ctx context.Context) error {
	url := fmt.Sprintf("http://%s/healthz", s.InsecureServingInfo.Address)
	if strings.Contains(s.InsecureServingInfo.Address, "0.0.0.0") {
		url = fmt.Sprintf("http://127.0.0.1:%s/healthz", strings.Split(s.InsecureServingInfo.Address, ":")[1])
	}

	for {
		// Change NewRequest to NewRequestWithContext and pass context it
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return err
		}
		// Ping the server by sending a GET request to `/healthz`.

		resp, err := http.DefaultClient.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			log.Info("The router has been deployed successfully.")

			resp.Body.Close()

			return nil
		}

		// Sleep for a second to continue the next ping.
		log.Info("Waiting for the router, retry in 1 second.")
		time.Sleep(1 * time.Second)

		select {
		case <-ctx.Done():
			log.Fatal("can not ping http server within the specified time interval.")
		default:
		}
	}
	// return fmt.Errorf("the router has no response, or it might took too long to start up")
}
