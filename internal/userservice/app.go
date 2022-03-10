// Created by Hisen at 2022/3/1.
package userservice

import (
	"github.com/hanxinhisen/moss/internal/userservice/config"
	"github.com/hanxinhisen/moss/internal/userservice/options"
	"github.com/hanxinhisen/moss/pkg/app"
	"github.com/hanxinhisen/moss/pkg/log"
)

const commandDesc = `
moss cmdb system
`

// NewApp 创建APP
func NewApp(basename string) *app.App {
	// 初始化默认选项
	opts := options.NewOptions()
	app := app.NewApp("cmdb",
		basename,
		app.WithOptions(opts),
		app.WithDescription(commandDesc),
		app.WithDefaultValidArgs(),
		app.WithRunFunc(run(opts)),
	)
	return app
}

func run(opts *options.Options) app.RunFunc {
	return func(basename string) error {
		log.Init(opts.Log)
		defer log.Flush()
		cfg, err := config.CreateConfigFromOptions(opts)
		if err != nil {
			return err
		}
		return Run(cfg)
	}
}
