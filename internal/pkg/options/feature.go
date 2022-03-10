// Created by Hisen at 2022/3/1.
package options

import (
	"github.com/hanxinhisen/moss/internal/pkg/server"
	"github.com/spf13/pflag"
)

type FeatureOptions struct {
	EnableProfiling bool `json:"profiling"      mapstructure:"profiling"`
	EnableMetrics   bool `json:"enable-metrics" mapstructure:"enable-metrics"`
}

func NewFeatureOptions() *FeatureOptions {
	defaults := server.NewConfig()

	return &FeatureOptions{
		EnableMetrics:   defaults.EnableMetrics,
		EnableProfiling: defaults.EnableProfiling,
	}
}

func (o *FeatureOptions) ApplyTo(c *server.Config) error {
	c.EnableProfiling = o.EnableProfiling
	c.EnableMetrics = o.EnableMetrics
	return nil

}

func (o *FeatureOptions) AddFlags(fs *pflag.FlagSet) {
	if fs == nil {
		return
	}

	fs.BoolVar(&o.EnableProfiling, "feature.profiling", o.EnableProfiling,
		"Enable profiling via web interface host:port/debug/pprof/")

	fs.BoolVar(&o.EnableMetrics, "feature.enable-metrics", o.EnableMetrics,
		"Enables metrics on the apiserver at /metrics")
}

func (o *FeatureOptions) Validate() []error {
	return []error{}
}
