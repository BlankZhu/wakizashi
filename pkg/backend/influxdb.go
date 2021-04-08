package backend

import (
	"BlankZhu/wakizashi/pkg/config"
	"BlankZhu/wakizashi/pkg/entity"
	"time"

	iclient "github.com/influxdata/influxdb1-client/v2"
	"github.com/sirupsen/logrus"
)

type influxClient struct {
	client iclient.Client
	cfg    config.BackendConfig
}

func CreateInfluxClient(cfg config.BackendConfig) *influxClient {
	cli, err := iclient.NewHTTPClient(
		iclient.HTTPConfig{
			Addr:     cfg.InfluxCfg.Host,
			Username: cfg.InfluxCfg.User,
			Password: cfg.InfluxCfg.Password,
			Timeout:  time.Duration(cfg.Timeout) * time.Second,
		},
	)
	if err != nil {
		logrus.Fatalf("failed to create influxDB client, detail: %s", err)
	}
	ret := &influxClient{
		client: cli,
		cfg:    cfg,
	}
	return ret
}

func (ic *influxClient) Connect() error {
	_, _, err := ic.client.Ping(time.Duration(ic.cfg.Timeout) * time.Second)
	return err
}

func (ic *influxClient) Close() error {
	return ic.client.Close()
}

func (ic *influxClient) Write(record *entity.TrafficRecord) error {
	bps, _ := iclient.NewBatchPoints(iclient.BatchPointsConfig{
		Database:  ic.cfg.Database,
		Precision: "s",
	})
	// makeup point
	tags := map[string]string{
		"probeIP": record.ProbeIP,
		"srcIP":   record.SrcIP,
		"dstIP":   record.DstIP,
	}
	fields := map[string]interface{}{
		"size": record.Size,
	}
	pt, err := iclient.NewPoint(ic.cfg.Table, tags, fields, time.Unix(record.Timestamp, 0))

	if err != nil {
		return err
	}
	bps.AddPoint(pt)
	return ic.client.Write(bps)
}

func (ic *influxClient) WriteBatch(record []*entity.TrafficRecord) error {
	bps, _ := iclient.NewBatchPoints(iclient.BatchPointsConfig{
		Database:  ic.cfg.Database,
		Precision: "s",
	})

	for _, p := range record {
		// makeup point
		tags := map[string]string{
			"probeIP": p.ProbeIP,
			"srcIP":   p.SrcIP,
			"dstIP":   p.DstIP,
		}
		fields := map[string]interface{}{
			"size": p.Size,
		}
		pt, err := iclient.NewPoint(ic.cfg.Table, tags, fields, time.Unix(p.Timestamp, 0))

		if err != nil {
			return err
		}
		bps.AddPoint(pt)
	}

	return ic.client.Write(bps)
}
