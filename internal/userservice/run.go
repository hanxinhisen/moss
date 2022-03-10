package userservice

import "github.com/hanxinhisen/moss/internal/userservice/config"

func Run(cfg *config.Config) error {
	server, err := createAPIServer(cfg)
	if err != nil {
		return err
	}
	return server.PrepareRun().Run()
}
