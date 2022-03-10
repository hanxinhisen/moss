// Created by Hisen at 2022/3/2.
package mysql

import (
	"fmt"
	"github.com/hanxinhisen/moss/internal/pkg/logger"
	genericoptions "github.com/hanxinhisen/moss/internal/pkg/options"
	"github.com/hanxinhisen/moss/internal/userservice/store"
	"github.com/hanxinhisen/moss/pkg/db"
	"github.com/marmotedu/errors"
	"gorm.io/gorm"

	"sync"
)

type datastore struct {
	db *gorm.DB
}

func (ds *datastore) Close() error {
	db, err := ds.db.DB()
	if err != nil {
		return errors.Wrap(err, "get gorm db instance failed")
	}

	return db.Close()
}

func (ds *datastore) Users() store.UserStore {
	return newUsers(ds)
}

var (
	mysqlFactory store.Factory
	once         sync.Once
)

func GetMysqlFactoryOr(opts *genericoptions.MySQLOptions) (store.Factory, error) {
	if opts == nil && mysqlFactory == nil {
		return nil, fmt.Errorf("failed to get mysql store fatory")
	}

	var err error
	var dbIns *gorm.DB

	once.Do(func() {
		options := &db.Options{
			Host:                  opts.Host,
			Username:              opts.Username,
			Password:              opts.Password,
			Database:              opts.Database,
			MaxIdleConnections:    opts.MaxIdleConnections,
			MaxOpenConnections:    opts.MaxOpenConnections,
			MaxConnectionLifeTime: opts.MaxConnectionLifeTime,
			LogLevel:              opts.LogLevel,
			Logger:                logger.New(opts.LogLevel),
		}

		dbIns, err = db.New(options)
		mysqlFactory = &(datastore{db: dbIns})
	})

	if mysqlFactory == nil || err != nil {
		return nil, fmt.Errorf("failed to get mysql store fatory, mysqlFactory: %+v, error: %w", mysqlFactory, err)
	}
	return mysqlFactory, nil
}
