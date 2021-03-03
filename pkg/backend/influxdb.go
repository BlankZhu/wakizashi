package backend

import (
	"BlankZhu/wakizashi/pkg/config"
	"BlankZhu/wakizashi/pkg/entity"
	"sync"

	iclient "github.com/influxdata/influxdb1-client/v2"
	"github.com/sirupsen/logrus"
)

var iClient *influxClient
var iOnce sync.Once

type influxClient struct {
	client *iclient.Client
}

func (ic *influxClient) Init(cfg *config.BackendConfig) {
	iOnce.Do(func() {
		cli, err := iclient.NewHTTPClient(
			iclient.HTTPConfig{
				Addr:     cfg.Host,
				Username: cfg.User,
				Password: cfg.Password,
			},
		)
		if err != nil {
			logrus.Fatalf("failed to create influxDB client, detail: %s", err)
		}
		iClient = &influxClient{
			client: &cli,
		}
	})
}

func (ic *influxClient) Close() {
	// TODO
}

func (ic *influxClient) Write(record *entity.TrafficRecord) error {
	// TODO
	return nil
}

func (ic *influxClient) WriteBatch(record []*entity.TrafficRecord) error {
	// TODO
	return nil
}
