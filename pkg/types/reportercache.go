package types

import (
	"BlankZhu/wakizashi/pkg/entity"
	"sync"
)

// ReporterCache cache type used by wakizashi's reporter
type ReporterCache struct {
	sync.RWMutex
	Data map[string]*entity.TrafficRecord
}

// Init initializes the ReporterCache
func (rc *ReporterCache) Init() {
	rc.Data = make(map[string]*entity.TrafficRecord)
}
