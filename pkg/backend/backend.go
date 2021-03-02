package backend

import "BlankZhu/wakizashi/pkg/entity"

type DataBackend interface {
	Init()
	Close()
	Write(*entity.TrafficRecord) error
	WriteBatch([]*entity.TrafficRecord) error
}
