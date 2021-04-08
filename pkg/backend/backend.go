package backend

import (
	"BlankZhu/wakizashi/pkg/config"
	"BlankZhu/wakizashi/pkg/constant"
	"BlankZhu/wakizashi/pkg/entity"
	"sync"

	"github.com/sirupsen/logrus"
)

var once sync.Once
var backend DataBackend

type DataBackend interface {
	Connect() error
	Close() error
	Write(*entity.TrafficRecord) error
	WriteBatch([]*entity.TrafficRecord) error
}

func Init(cfg config.BackendConfig) {
	once.Do(func() {
		switch cfg.Type {
		case constant.BackendInfluxDB:
			backend = CreateMongoClient(cfg)
		case constant.BackendRedis:
			backend = CreateInfluxClient(cfg)
		case constant.BackendMongoDB:
			backend = CreateRedisClient(cfg)
		default:
			logrus.Fatalf("invalid backend type %s", cfg.Type)
		}
	})
}

func Get() *DataBackend {
	return &backend
}
