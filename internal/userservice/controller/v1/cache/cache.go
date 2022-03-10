// Created by Hisen at 2022/3/3.
package cache

import (
	"context"
	"fmt"
	pb "github.com/hanxinhisen/moss/internal/userservice/proto/v1"
	"github.com/hanxinhisen/moss/internal/userservice/store"
	"sync"
)

type Cache struct {
	store store.Factory
}

var (
	cacheServer *Cache
	once        sync.Once
)

func GetCacheInsOr(store store.Factory) (*Cache, error) {
	if store != nil {
		once.Do(func() {
			cacheServer = &Cache{store: store}
		})

	}
	if cacheServer == nil {
		return nil, fmt.Errorf("got nil cache server")
	}
	return cacheServer, nil
}

func (c Cache) ListSecrets(ctx context.Context, request *pb.ListSecretsRequest) (*pb.ListSecretsResponse, error) {
	//TODO implement me
	return nil, nil
}

func (c Cache) ListPolicies(ctx context.Context, request *pb.ListPoliciesRequest) (*pb.ListPoliciesResponse, error) {
	//TODO implement me
	return nil, nil
}
