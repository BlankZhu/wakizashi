package backend

import (
	"BlankZhu/wakizashi/pkg/config"
	"BlankZhu/wakizashi/pkg/entity"
)

type RedisClient struct {
}

func createRedisClient(cfg config.BackendConfig) *RedisClient {
	// TODO
	ret := &RedisClient{}
	return ret
}

func (rc *RedisClient) Connect() error {
	// TODO
	return nil
}

func (rc *RedisClient) Close() error {
	// TODO
	return nil
}

func (rc *RedisClient) Write(record *entity.TrafficRecord) error {
	// TODO
	return nil
}

func (rc *RedisClient) WriteBatch(record []*entity.TrafficRecord) error {
	// TODO
	return nil
}
