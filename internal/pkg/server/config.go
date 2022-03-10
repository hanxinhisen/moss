// Created by Hisen at 2022/3/1.
package server

import (
	"github.com/gin-gonic/gin"
	"net"
	"strconv"
	"time"
)

const (
	RecommendedHomeDir   = ".moss"
	RecommendedEnvPrefix = "MOSS"
)

type Config struct {
	SecureServing   *SecureServingInfo
	InsecureServing *InsecureServingInfo
	Jwt             *JwtInfo
	Mode            string
	Middlewares     []string
	Healthz         bool
	EnableProfiling bool
	EnableMetrics   bool
}

func NewConfig() *Config {
	return &Config{
		Healthz:         true,
		Mode:            gin.ReleaseMode,
		Middlewares:     []string{},
		EnableProfiling: true,
		EnableMetrics:   true,
		Jwt: &JwtInfo{
			Realm:      "moss jwt",
			Timeout:    1 * time.Hour,
			MaxRefresh: 1 * time.Hour,
		}}
}

type SecureServingInfo struct {
	BindAddress string
	BindPort    int
	CertKey     CertKey
}

func (s *SecureServingInfo) Address() string {
	return net.JoinHostPort(s.BindAddress, strconv.Itoa(s.BindPort))
}

type CertKey struct {
	CertFile string
	KeyFile  string
}

type InsecureServingInfo struct {
	Address string
}

type JwtInfo struct {
	Realm      string
	Key        string
	Timeout    time.Duration
	MaxRefresh time.Duration
}

type CompleteConfig struct {
	*Config
}

func (c *Config) Complete() CompleteConfig {
	return CompleteConfig{c}
}

func (c CompleteConfig) New() (*GenericAPIServer, error) {
	gin.SetMode(c.Mode)
	s := &GenericAPIServer{
		SecureServingInfo:   c.SecureServing,
		InsecureServingInfo: c.InsecureServing,
		healthz:             c.Healthz,
		enableMetrics:       c.EnableMetrics,
		enableProfiling:     c.EnableProfiling,
		middlewares:         c.Middlewares,
		Engine:              gin.New(),
	}
	initGenericAPIServer(s)
	return s, nil

}
