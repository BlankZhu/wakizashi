package backend

import (
	"BlankZhu/wakizashi/pkg/config"
	"BlankZhu/wakizashi/pkg/entity"
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	client *mongo.Client
	cfg    config.BackendConfig
}

func createMongoClient(cfg config.BackendConfig) *MongoClient {
	cli, err := mongo.NewClient(options.Client().ApplyURI(cfg.MongoCfg.MongoURI))
	if err != nil {
		logrus.Fatalf("failed to create mongoDB client, detail: %s", err)
	}
	ret := &MongoClient{
		client: cli,
		cfg:    cfg,
	}
	return ret
}

func (mc *MongoClient) Connect() error {
	ctx, cancel := mc.makeContextWithTimeout()
	defer cancel()
	return mc.client.Connect(ctx)
}

func (mc *MongoClient) Close() error {
	ctx, cancel := mc.makeContextWithTimeout()
	defer cancel()
	return mc.client.Disconnect(ctx)
}

func (mc *MongoClient) Write(record *entity.TrafficRecord) error {
	coll := mc.client.Database(mc.cfg.Database).Collection(mc.cfg.Table)
	ctx, cancel := mc.makeContextWithTimeout()
	defer cancel()
	_, err := coll.InsertOne(ctx, *record)
	return err
}

func (mc *MongoClient) WriteBatch(record []*entity.TrafficRecord) error {
	coll := mc.client.Database(mc.cfg.Database).Collection(mc.cfg.Table)
	ctx, cancel := mc.makeContextWithTimeout()
	defer cancel()

	iRecord := make([]interface{}, 0, len(record))
	for _, p := range record {
		iRecord = append(iRecord, *p)
	}
	_, err := coll.InsertMany(ctx, iRecord)
	return err
}

func (mc *MongoClient) makeContextWithTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(mc.cfg.Timeout)*time.Second)
}
