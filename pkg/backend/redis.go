package backend

import "BlankZhu/wakizashi/pkg/entity"

type RedisClient struct {
}

func (rc *RedisClient) Init() {
	// TODO
}

func (rc *RedisClient) Close() {
	// TODO
}

func (rc *RedisClient) Write(record *entity.TrafficRecord) error {
	// TODO
	return nil
}

func (rc *RedisClient) WriteBatch(record []*entity.TrafficRecord) error {
	// TODO
	return nil
}
