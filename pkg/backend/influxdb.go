package backend

import "BlankZhu/wakizashi/pkg/entity"

type InfluxClient struct {
}

func (ic *InfluxClient) Init() {
	// TODO
}

func (ic *InfluxClient) Close() {
	// TODO
}

func (ic *InfluxClient) Write(record *entity.TrafficRecord) error {
	// TODO
	return nil
}

func (ic *InfluxClient) WriteBatch(record []*entity.TrafficRecord) error {
	// TODO
	return nil
}
