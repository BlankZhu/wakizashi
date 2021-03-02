package backend

import "BlankZhu/wakizashi/pkg/entity"

type MongoClient struct {
}

func (mc *MongoClient) Init() {
	// TODO
}

func (mc *MongoClient) Close() {
	// TODO
}

func (mc *MongoClient) Write(record *entity.TrafficRecord) error {
	// TODO
	return nil
}

func (mc *MongoClient) WriteBatch(record []*entity.TrafficRecord) error {
	// TODO
	return nil
}
